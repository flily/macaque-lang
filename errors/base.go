package errors

import (
	"fmt"
)

const (
	UnknownError       = iota
	ErrCodeSyntaxError = 1
	ErrScannerError    = 100
)

type BaseError struct {
	Kind    int
	Message string
	Base    *BaseError
}

func NewRawError(kind int, format string, args ...interface{}) *BaseError {
	e := &BaseError{
		Kind:    kind,
		Message: fmt.Sprintf(format, args...),
	}

	return e
}

func NewError(kind int, format string, args ...interface{}) error {
	return NewRawError(kind, format, args...)
}

func (e *BaseError) Error() string {
	return e.Message
}

func (e *BaseError) Is(target error) bool {
	if err, ok := target.(*BaseError); ok {
		return err.Kind == e.Kind
	}

	return false
}

func Derive(base error, format string, args ...interface{}) error {
	if err, ok := base.(*BaseError); ok {
		return err.Derive(format, args...)
	}

	return NewError(UnknownError, format, args...)
}

func (e *BaseError) Derive(format string, args ...interface{}) error {
	return Derive(e, format, args...)
}
