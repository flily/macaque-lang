package compiler

import (
	"github.com/flily/macaque-lang/opcode"
)

type CodeBuffer struct {
	Code []opcode.Opcode
}

func (b *CodeBuffer) Length() int {
	return len(b.Code)
}

func (b *CodeBuffer) Write(code int, operands ...int) {
	c := opcode.Code(code, operands...)
	b.Code = append(b.Code, c)
}

func (b *CodeBuffer) Append(buf *CodeBuffer) {
	b.Code = append(b.Code, buf.Code...)
}

func (b *CodeBuffer) AppendCode(buf *CodeBuffer, err error) error {
	if err != nil {
		return err
	}

	b.Append(buf)
	return nil
}

func NewCodeBuffer() *CodeBuffer {
	b := &CodeBuffer{
		Code: make([]opcode.Opcode, 0),
	}

	return b
}

type CompileResult struct {
	Code   *CodeBuffer
	Values int
}

func NewCompileResult() *CompileResult {
	r := &CompileResult{
		Code:   NewCodeBuffer(),
		Values: 0,
	}

	return r
}

func (r *CompileResult) Append(result *CompileResult, err error) error {
	if err != nil {
		return err
	}

	r.Code.Append(result.Code)
	r.Values += result.Values
	return nil
}

func (r *CompileResult) AppendCode(result *CompileResult) {
	r.Code.Append(result.Code)
}

func (r *CompileResult) Write(code int, operands ...int) {
	r.Code.Write(code, operands...)
}
