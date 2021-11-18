package schema

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/99designs/gqlgen/graphql"
)

type Metadata map[string]interface{}

type Labels map[string]string

func LabelsFromMap(m map[string]string) *Labels {
	var labels = make(Labels, len(m))
	for k, v := range m {
		labels[k] = v
	}
	return &labels
}

func (l *Labels) String() string {
	if l == nil {
		return ""
	}
	var mapStr = make([]string, 0, len(*l))
	for k, v := range *l {
		mapStr = append(mapStr, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(mapStr, ", ")
}

func (l *Labels) ToMap() map[string]string {
	if l == nil {
		return nil
	}
	var m = make(map[string]string, len(*l))
	for k, v := range *l {
		m[k] = v
	}
	return m
}

func MarshalLabels(val map[string]string) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		err := json.NewEncoder(w).Encode(val)
		if err != nil {
			panic(err)
		}
	})
}

func UnmarshalLabels(v interface{}) (map[string]string, error) {
	if m, ok := v.(map[string]string); ok {
		return m, nil
	}

	return nil, fmt.Errorf("%T is not a label", v)
}
