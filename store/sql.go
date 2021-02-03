package store

import (
	"fmt"
	"strings"

	"github.com/verifa/bubbly/api/core"
	"github.com/verifa/bubbly/parser"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"

	sq "github.com/Masterminds/squirrel"
)

const tableIDField = "_id"
const tableJoinSuffix = "_id"

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func sqlTableCreate(table core.Table) (string, error) {
	var (
		index        = 0
		fieldLen     = len(table.Fields) + len(table.Joins)
		tableFields  = make([]string, fieldLen)
		uniqueFields = make([]string, 0)
	)

	tableFields = append(tableFields, tableIDField+" SERIAL PRIMARY KEY")
	// Add the fields to the SQL table
	for _, field := range table.Fields {
		sqlType, err := sqlType(field.Type)
		if err != nil {
			return "", fmt.Errorf("failed to create SQL statement for table: %s: %w", table.Name, err)
		}
		tableFields[index] = field.Name + " " + sqlType

		if field.Unique {
			uniqueFields = append(uniqueFields, field.Name)
		}
		index++
	}
	// Add the joins as fields to the SQL table
	for _, join := range table.Joins {
		tableFields[index] = join.Table + "_id SERIAL"
		index++
	}

	if len(uniqueFields) > 0 {
		tableFields = append(tableFields, fmt.Sprintf("UNIQUE (%s)", strings.Join(uniqueFields, ",")))
	}

	// TODO:nate wait until merged with master
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ( %s );", table.Name, strings.Join(tableFields, ",")), nil
}

func psqlDataNodeUpsert(node *dataNode, table core.Table) (sq.InsertBuilder, error) {
	var (
		data           = node.Data
		fieldNames     = node.orderedFields()
		conflictValues = make([]string, 0, len(data.Fields))
		uniqueFields   = make([]string, 0)
		sqlOnConflict  = ""
		sqlReturning   = ""
	)

	for _, field := range table.Fields {
		if field.Unique {
			uniqueFields = append(uniqueFields, field.Name)
		}
	}

	for _, name := range fieldNames {
		conflictValues = append(conflictValues, name+"=EXCLUDED."+name)
	}

	// Create the UPSERT / ON CONFLICT part of the SQL statement, if any.
	if len(uniqueFields) > 0 {
		sqlOnConflict = fmt.Sprintf(
			"ON CONFLICT (%s) DO UPDATE SET %s",
			strings.Join(uniqueFields, ","),
			strings.Join(conflictValues, ","),
		)
	}

	// Create the RETURNING part of the SQL statement, if any.
	sqlReturning = "RETURNING " + strings.Join(node.orderedRefFields(), ",")
	values, err := sqlArgValues(node)
	if err != nil {
		return sq.InsertBuilder{}, fmt.Errorf("failed to get SQL arguments: %w", err)
	}
	return psql.Insert(data.TableName).
		Columns(fieldNames...).
		Values(values...).
		Suffix(sqlOnConflict).
		Suffix(sqlReturning), nil
}

func applyGraphQLArgs(sql sq.SelectBuilder, args map[string]interface{}) sq.SelectBuilder {
	for k, v := range args {

		switch {
		case k == firstID:
			sql = sql.Limit(uint64(v.(int)))
			continue
		case k != filterID:
			sql = sql.Where(sq.Eq{k: v})
			continue
		default:
			// Do nothing, proceed with the function and check the suffixes
		}

		hasSuffix := func(n string, filter string) (bool, string) {
			if !strings.HasSuffix(n, filter) {
				return false, ""
			}
			return true, strings.ReplaceAll(n, filter, "")
		}

		filter := v.(map[string]interface{})
		for n, val := range filter {
			if ok, f := hasSuffix(n, filterGreaterThan); ok {
				sql = sql.Where(f+" > ?", val)
				continue
			}
			if ok, f := hasSuffix(n, filterLessThan); ok {
				sql = sql.Where(f+" < ?", val)
				continue
			}
			if ok, f := hasSuffix(n, filterGreaterThanOrEqualTo); ok {
				sql = sql.Where(f+" >= ?", val)
				continue
			}
			if ok, f := hasSuffix(n, filterLessThanOrEqualTo); ok {
				sql = sql.Where(f+" <= ?", val)
				continue
			}
			// TODO: this is not so easy because with raw SQL statements we need
			// to use "WHERE some_field IN ($3, $4)" and provide args at
			// position 3 and 4 (in this case) to the SQL exec...
			// How do we know these indexes?
			// Important: test "_not_in" first so we don't
			// accidentally match "_in".
			if ok, f := hasSuffix(n, filterNotIn); ok {
				sql = sql.Where(sq.Eq{f: val.([]interface{})})
				continue
			}
			if ok, f := hasSuffix(n, filterIn); ok {
				sql = sql.Where(sq.NotEq{f: val.([]interface{})})
			}
		}
	}
	return sql
}

func sqlArgValues(node *dataNode) ([]interface{}, error) {
	var (
		data   = node.Data
		values = make([]interface{}, 0, len(data.Fields))
	)

	for _, f := range node.orderedFields() {
		val, err := sqlValue(node, data.Fields[f])
		if err != nil {
			return nil, fmt.Errorf("failed to get SQL value from cty.Value for field: %s: %w", f, err)
		}
		values = append(values, val)
	}
	return values, nil
}

func sqlValue(node *dataNode, val cty.Value) (interface{}, error) {
	// Check if the value is a capsule value, in which case it needs special
	// treatment
	if val.Type().IsCapsuleType() {
		// Get the underlying DataRef type
		ref := val.EncapsulatedValue().(*parser.DataRef)
		if parent, ok := node.Parents[ref.TableName]; ok {
			if ret, ok := parent.Return[ref.Field]; ok {
				return ret, nil
			}
		}
		return nil, fmt.Errorf("could not find data ref: %s.%s", ref.TableName, ref.Field)
	}

	// If not a capsule type, it is a regular cty.Value
	return valueFromCty(val)
}

func valueFromCty(val cty.Value) (interface{}, error) {
	switch ty := val.Type(); {
	case ty == cty.Bool:
		return val.True(), nil
	case ty == cty.Number:
		var number int
		if err := gocty.FromCtyValue(val, &number); err != nil {
			return nil, fmt.Errorf("failed to convert cty.Value to int: %s: %w", val.GoString(), err)
		}
		return number, nil
	case ty == cty.String:
		return val.AsString(), nil
	case ty.IsObjectType():
		var (
			m   = val.AsValueMap()
			ret = make(map[string]interface{}, len(m))
		)
		for k, v := range m {
			var err error
			ret[k], err = valueFromCty(v)
			if err != nil {
				return nil, err
			}
		}
		return ret, nil
	default:
		return nil, fmt.Errorf("unsupported cty type: %s", ty.GoString())
	}
}

// sqlType takes a cty.Type and returns a string representation of the
// corresponding SQL type
func sqlType(ty cty.Type) (string, error) {
	switch {
	case ty == cty.Bool:
		return "BOOL", nil
	case ty == cty.Number:
		return "INT8", nil
	case ty == cty.String:
		return "TEXT", nil
	case ty.IsObjectType():
		return "JSONB", nil
	default:
		return "", fmt.Errorf("unsupported SQL type: %s", ty.GoString())
	}
}
