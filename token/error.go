package token

const (
	UnknownError            = iota
	ErrCodeSyntaxError      = 1
	ErrCodeRuntimeError     = 2
	ErrCodeCompilationError = 3
	ErrScannerError         = 100
)
