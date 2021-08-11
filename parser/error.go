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
	var (
		msgs     []string
		prevDiag *hcl.Diagnostic
	)
	for _, err := range e.Diags.Errs() {
		// This is actually a *hcl.Diagnostic, we just need to cast it, e.g.
		if diag, ok := err.(*hcl.Diagnostic); ok {
			// The use of HCL dynamic blocks can create a lot of duplicate messages.
			// We only need to show one of those and they come sequentially, so
			// compare this diagnostic with the previous one
			if prevDiag != nil && prevDiag.Subject.String() == diag.Subject.String() &&
				prevDiag.Detail == diag.Detail {
				prevDiag = diag
				continue
			}
			var errMsg string

			errMsg = err.Error()
			// errMsg += "\nEvalContext:\n"
			// errMsg += appendEvalContext(diag.EvalContext)
			// If it's not a duplicate, add it
			msgs = append(msgs, errMsg)
			prevDiag = diag
			continue
		}
		// If err is not an hcl.Diagnostic just add the error message normally
		msgs = append(msgs, err.Error())
	}
	return "\n" + strings.Join(msgs, "\n")
}

func appendEvalContext(eCtx *hcl.EvalContext) string {
	var errMsg string
	if eCtx == nil {
		return ""
	}
	errMsg += "\nNext EvalContext Scope:\n"
	for name, val := range eCtx.Variables {
		errMsg += name + ": " + val.GoString() + "\n"
	}
	return errMsg + appendEvalContext(eCtx.Parent())
}
