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
	for _, err := range e.Diags.Errs() {
		// This is actually a *hcl.Diagnostic, we just need to cast it, e.g.
		// if diag, ok := err.(*hcl.Diagnostic); ok { ... }
		msgs = append(msgs, err.Error())
	}
	return "\n" + strings.Join(msgs, "\n")
}
