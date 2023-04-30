package token

var (
	ErrSyntaxError = NewPositionError(ErrCodeSyntaxError, "syntax error")
)

type SyntaxError struct {
	PositionError
}

func NewSyntaxError(context *CodeContext, format string, args ...interface{}) *SyntaxError {
	base := ErrSyntaxError.Derive(format, args...).WithContext(context)
	e := &SyntaxError{
		PositionError: *base,
	}

	return e
}

func (e *SyntaxError) Error() string {
	return e.Context.MakeMessage(e.Message)
}
