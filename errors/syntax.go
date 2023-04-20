package errors

var (
	ErrSyntaxError = NewRawError(ErrCodeSyntaxError, "syntax error")
)

type SyntaxError struct {
	BaseError
}

func NewSyntaxError(context *CodeContext, format string, args ...interface{}) *SyntaxError {
	base := ErrSyntaxError.Derive(format, args...).WithContext(context)
	e := &SyntaxError{
		BaseError: *base,
	}

	return e
}

func (e *SyntaxError) Error() string {
	return e.Context.MakeMessage(e.Message)
}
