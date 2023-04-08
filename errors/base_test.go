package errors

import (
	"io"
	"testing"
)

func TestBaseErrorMessage(t *testing.T) {
	err := NewError(ErrCodeSyntaxError, "test error")
	if err.Error() != "test error" {
		t.Errorf("error message is not correct: %s", err.Error())
	}

	e1 := NewError(ErrCodeSyntaxError, "test error2")
	if !Is(err, e1) {
		t.Errorf("error should be equal: %s", err.Error())
	}

	e2 := NewError(ErrScannerError, "test error3")
	if Is(err, e2) {
		t.Errorf("error should not be equal: %s", err.Error())
	}

	if Is(err, io.EOF) {
		t.Errorf("error should not be equal: %s", err.Error())
	}
}

func TestKindOf(t *testing.T) {
	e1 := NewError(ErrCodeSyntaxError, "test error")
	if KindOf(e1) != ErrCodeSyntaxError {
		t.Errorf("error kind is not correct: %d", KindOf(e1))
	}

	e2 := NewError(ErrScannerError, "test error")
	if KindOf(e2) != ErrScannerError {
		t.Errorf("error kind is not correct: %d", KindOf(e2))
	}

	e3 := io.EOF
	if KindOf(e3) != UnknownError {
		t.Errorf("error kind is not correct: %d", KindOf(e3))
	}
}

func TestBaseErrorDerive(t *testing.T) {
	err := NewRawError(ErrCodeSyntaxError, "test error")

	e1 := err.Derive("test error1")
	if KindOf(e1) != ErrCodeSyntaxError {
		t.Errorf("error kind is not correct: %d", KindOf(e1))
	}

	if e1.Error() != "test error1" {
		t.Errorf("error message is not correct: %s", e1.Error())
	}

	e2 := Derive(err, "test error2")
	if KindOf(e2) != ErrCodeSyntaxError {
		t.Errorf("error kind is not correct: %d", KindOf(e2))
	}

	if e2.Error() != "test error2" {
		t.Errorf("error message is not correct: %s", e2.Error())
	}

	e3 := Derive(io.EOF, "test error3")
	if KindOf(e3) != UnknownError {
		t.Errorf("error kind is not correct: %d", KindOf(e3))
	}

	if e3.Error() != "test error3" {
		t.Errorf("error message is not correct: %s", e3.Error())
	}
}
