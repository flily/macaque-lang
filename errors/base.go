package errors

import (
	"errors"
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

type HasKind interface {
	ErrorKind() int
}

type CanDerive interface {
	Derive(format string, args ...interface{}) error
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

func KindOf(err error) int {
	if err, ok := err.(HasKind); ok {
		return err.ErrorKind()
	}

	return UnknownError
}

func (e *BaseError) ErrorKind() int {
	return e.Kind
}

// Is wraps errors.Is
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func (e *BaseError) Is(target error) bool {
	if err, ok := target.(HasKind); ok {
		return err.ErrorKind() == e.ErrorKind()
	}

	return false
}

func Derive(base error, format string, args ...interface{}) error {
	if err, ok := base.(CanDerive); ok {
		return err.Derive(format, args...)
	}

	return NewError(UnknownError, format, args...)
}

func (e *BaseError) Derive(format string, args ...interface{}) error {
	err := NewRawError(e.Kind, format, args...)
	err.Base = e
	return err
}
