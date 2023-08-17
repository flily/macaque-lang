package compiler

import "fmt"

type CompilerFlag uint64

const (
	FlagNone      = 0x0000
	FlagPackValue = 0x0001

	// Non-passable flags
	FlagNonPassable = 0x00ff
	FlagCleanStack  = 0x0100
	FlagWithReturn  = 0x0200
)

func NewFlag(flags ...uint64) CompilerFlag {
	f := uint64(FlagNone)
	for _, flag := range flags {
		f |= flag
	}

	return CompilerFlag(f)
}

func (f CompilerFlag) Has(flag uint64) bool {
	return (uint64(f) & flag) != 0
}

func (f CompilerFlag) With(flag uint64) CompilerFlag {
	return f | CompilerFlag(flag)
}

func (f CompilerFlag) Without(flag uint64) CompilerFlag {
	return f & ^CompilerFlag(flag)
}

func (f *CompilerFlag) Set(flag uint64) CompilerFlag {
	*f |= CompilerFlag(flag)
	return *f
}

func (f *CompilerFlag) Clear(flag uint64) CompilerFlag {
	*f &= ^CompilerFlag(flag)
	return *f
}

func (f *CompilerFlag) ClearNonPassable() CompilerFlag {
	*f &= CompilerFlag(FlagNonPassable)
	return *f
}

func (f CompilerFlag) String() string {
	n := uint64(f)
	w0 := (n >> 48) & 0xffff
	w1 := (n >> 32) & 0xffff
	w2 := (n >> 16) & 0xffff
	w3 := (n >> 0) & 0xffff
	return fmt.Sprintf("CompilerFlag(%04x-%04x-%04x-%04x)", w0, w1, w2, w3)
}
