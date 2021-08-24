package store

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-multierror"
	"github.com/valocode/bubbly/ent"
)

func HandleMultiVError(vErr multierror.Error) error {
	var ret = make([]string, 0, len(vErr.Errors))
	for _, err := range vErr.Errors {
		ret = append(ret, err.Error())
	}
	return NewValidationError(nil, strings.Join(ret, ", "))
}

func HandleValidatorError(err error, msg string) error {
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return NewServerError(err, "validation error for "+msg)
	}
	if vErrs, ok := err.(validator.ValidationErrors); ok {
		var ret []string
		for _, vErr := range vErrs {
			ns := vErr.Namespace()
			if idx := strings.Index(ns, "."); idx > 0 {
				// Just double check we won't cause an error accessing idx+1
				if idx+1 <= len(ns) {
					ns = ns[idx+1:]
				}
			}
			ret = append(ret, ns+" "+vErr.Tag())
		}
		return NewValidationError(nil, msg+": "+strings.Join(ret, ", "))
	}
	return NewValidationError(nil, err.Error())
}

func HandleEntError(entErr error, msg string) error {

	switch {
	case ent.IsConstraintError(entErr):
		return NewConflictError(entErr, msg)
	case ent.IsNotFound(entErr):
		return NewNotFoundError(entErr, msg)
	case ent.IsNotSingular(entErr):
		// Handle not singluar error the same as not found.
		// The error is triggered by looking for one thing, and receiving more
		// which can be treated as a not found error.
		// Rethink sometime...
		return NewNotFoundError(entErr, msg)
	case ent.IsValidationError(entErr):
		return NewValidationError(entErr, msg)
	default:
		return NewServerError(entErr, msg)
	}
}

func NewConflictError(err error, msg string) *ConflictError {
	return &ConflictError{
		Msg: msg,
		Err: err,
	}
}

// ConflictError returns when validating a request fails
type ConflictError struct {
	Msg string
	Err error
}

// Error implements the error interface.
func (e *ConflictError) Error() string {
	return "conflict on resource " + e.Msg + ": " + e.Err.Error()
}

// IsConflictError returns a boolean indicating whether the error is a conflict error.
func IsConflictError(err error) bool {
	if err == nil {
		return false
	}
	var e *ConflictError
	return errors.As(err, &e)
}

func NewValidationError(err error, msg string) *ValidationError {
	return &ValidationError{
		Msg: msg,
		Err: err,
	}
}

// ValidationError returns when validating a request fails
type ValidationError struct {
	Msg string
	Err error
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	if e.Err == nil {
		return "validation error: " + e.Msg
	}
	return fmt.Sprintf("validation error on %s: %s", e.Msg, e.Err.Error())
}

// IsValidationError returns a boolean indicating whether the error is a validation error.
func IsValidationError(err error) bool {
	if err == nil {
		return false
	}
	var e *ValidationError
	return errors.As(err, &e)
}

func NewNotFoundError(err error, msg string, a ...interface{}) *NotFoundError {
	return &NotFoundError{
		msg: fmt.Sprintf(msg, a...),
		Err: err,
	}
}

// NotFoundError returns when trying to fetch a specific entity and it was not found in the database.
type NotFoundError struct {
	msg string
	Err error
}

// Error implements the error interface.
func (e *NotFoundError) Error() string {
	return e.msg
}

// IsNotFound returns a boolean indicating whether the error is a not found error.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	var e *NotFoundError
	return errors.As(err, &e)
}

func NewServerError(err error, msg string) *ServerError {
	return &ServerError{
		Msg: msg,
		Err: err,
	}
}

// ServerError returns when an unexpected error happened on the server
type ServerError struct {
	Msg string
	Err error
}

// Error implements the error interface.
func (e *ServerError) Error() string {
	return e.Err.Error()
}

// IsServerError returns a boolean indicating whether the error is a server error.
func IsServerError(err error) bool {
	if err == nil {
		return false
	}
	var e *ServerError
	return errors.As(err, &e)
}
