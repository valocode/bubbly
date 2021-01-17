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

	return fmt.Sprintf("CREATE TABLE %s ( %s );", table.Name, strings.Join(tableFields, ",")), nil
}

func sqlDataBlockUpsert(sc *saveContext, data core.Data, table core.Table, refFields []string) (sq.InsertBuilder, error) {
	var (
		lenFields      = len(data.Fields) + len(data.Joins)
		index          = 0
		fieldNames     = make([]string, lenFields)
		conflictValues = make([]string, lenFields)
		uniqueFields   = make([]string, 0)
		sqlOnConflict  = ""
		sqlReturning   = ""
	)

	for _, field := range table.Fields {
		if field.Unique {
			uniqueFields = append(uniqueFields, field.Name)
		}
	}

	for _, field := range data.Fields {
		fieldNames[index] = field.Name
		conflictValues[index] = field.Name + "=EXCLUDED." + field.Name
		index++
	}

	for _, join := range data.Joins {
		fieldNames[index] = join.Table + tableJoinSuffix
		conflictValues[index] = join.Table + tableJoinSuffix + "=EXCLUDED." + join.Table + tableJoinSuffix
		index++
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
	if len(refFields) > 0 {
		sqlReturning = "RETURNING " + strings.Join(refFields, ",")
	}
	values, err := sqlArgValues(sc, data)
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

func sqlArgValues(sc *saveContext, data core.Data) ([]interface{}, error) {
	var (
		lenValues = len(data.Fields) + len(data.Joins)
		index     = 0
		values    = make([]interface{}, lenValues)
	)

	for _, field := range data.Fields {
		val, err := sqlValue(sc, field.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to get SQL value from cty.Value for field: %s: %w", field.Name, err)
		}
		values[index] = val
		index++
	}
	for _, join := range data.Joins {
		val, err := sqlValue(sc, join.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to get SQL value from cty.Value for join: %s: %w", join.Table, err)
		}
		values[index] = val
		index++
	}
	return values, nil
}

func sqlValue(sc *saveContext, val cty.Value) (interface{}, error) {
	// Check if the value is a capsule value, in which case it needs special
	// treatment
	if val.Type().IsCapsuleType() {
		// Get the underlying DataRef type
		ref := val.EncapsulatedValue().(*parser.DataRef)
		if refVals, ok := sc.DataRefs[ref.TableName]; ok {
			for field, value := range refVals {
				if ref.Field == field {
					if value == nil {
						return nil, fmt.Errorf("data ref value has not been set: %s.%s", ref.TableName, ref.Field)
					}
					return value, nil
				}
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
		return "INT", nil
	case ty == cty.String:
		return "TEXT", nil
	case ty.IsObjectType():
		return "JSONB", nil
	default:
		return "", fmt.Errorf("unsupported SQL type: %s", ty.GoString())
	}
}
