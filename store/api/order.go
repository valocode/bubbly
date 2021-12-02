package api

import (
	"fmt"
	"strings"

	"github.com/valocode/bubbly/ent"
)

func OrderFromSortBy(sortBy string) (*Order, error) {
	var (
		field     string
		direction OrderDirection
		vals      = strings.Split(sortBy, ":")
	)
	switch {
	case len(vals) == 1:
		direction = "asc"
	case len(vals) == 2:
		direction = OrderDirection(vals[1])
		if err := direction.Validate(); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid sort_by parameter: %s", sortBy)
	}
	field = vals[0]

	return &Order{
		Field:     field,
		Direction: direction,
	}, nil
}

type Order struct {
	Field     string         `query:"sort_by"`
	Direction OrderDirection `query:"order"`
}

func (o Order) Func(field string) ent.OrderFunc {
	switch o.Direction {
	case OrderDirectionAsc:
		return ent.Asc(field)
	case OrderDirectionDesc:
		return ent.Asc(field)
	}
	return nil
}

type OrderDirection string

const (
	OrderDirectionAsc  OrderDirection = "asc"
	OrderDirectionDesc OrderDirection = "desc"
)

func (o OrderDirection) String() string {
	return string(o)
}

func (o OrderDirection) Validate() error {
	switch o {
	case OrderDirectionAsc, OrderDirectionDesc:
		return nil
	default:
		return fmt.Errorf("order direction invalid: %s", o)
	}
}

type OrderField string

func (o OrderField) String() string {
	return string(o)
}
