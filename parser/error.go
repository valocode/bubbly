package parser

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
)

func NewParserError(val interface{}, diags hcl.Diagnostics) *ParserError {
	return &ParserError{
		Diags: diags,
		Value: val,
	}
}

type ParserError struct {
	Diags hcl.Diagnostics
	Value interface{}
}

func (e *ParserError) Error() string {
	var msgs []string
	for _, diag := range e.Diags.Errs() {
		msgs = append(msgs, diag.Error())
	}
	return "\n" + strings.Join(msgs, "\n")
}
