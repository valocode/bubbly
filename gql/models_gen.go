// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gql

import (
	"fmt"
	"io"
	"strconv"
)

type ResultsType string

const (
	ResultsTypeCodeScan ResultsType = "code_scan"
	ResultsTypeTestRun  ResultsType = "test_run"
)

var AllResultsType = []ResultsType{
	ResultsTypeCodeScan,
	ResultsTypeTestRun,
}

func (e ResultsType) IsValid() bool {
	switch e {
	case ResultsTypeCodeScan, ResultsTypeTestRun:
		return true
	}
	return false
}

func (e ResultsType) String() string {
	return string(e)
}

func (e *ResultsType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ResultsType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ResultsType", str)
	}
	return nil
}

func (e ResultsType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Type string

const (
	TypeJSON Type = "json"
	TypeCsv  Type = "csv"
	TypeXML  Type = "xml"
	TypeYaml Type = "yaml"
	TypeHTTP Type = "http"
)

var AllType = []Type{
	TypeJSON,
	TypeCsv,
	TypeXML,
	TypeYaml,
	TypeHTTP,
}

func (e Type) IsValid() bool {
	switch e {
	case TypeJSON, TypeCsv, TypeXML, TypeYaml, TypeHTTP:
		return true
	}
	return false
}

func (e Type) String() string {
	return string(e)
}

func (e *Type) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Type(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Type", str)
	}
	return nil
}

func (e Type) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}