package errors

import (
	"fmt"
)

const (
	UnknownError = iota
	LexicalError = 1
)

type BaseError struct {
	Kind    int
	Message string
}

func NewError(kind int, format string, args ...interface{}) error {
	return &BaseError{
		Kind:    kind,
		Message: fmt.Sprintf(format, args...),
	}
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
