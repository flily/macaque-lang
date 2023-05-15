package vm

import (
	"github.com/flily/macaque-lang/errors"
)

type RuntimeError struct {
	errors.BaseError
}

func NewRuntimeError(format string, args ...interface{}) *RuntimeError {
	base := errors.NewRawError(errors.ErrCodeRuntimeError, format, args...)
	e := &RuntimeError{
		BaseError: *base,
	}

	return e
}
