package opcode

import (
	"testing"
)

func assertPanic(v interface{}, f func()) (has bool, ok bool) {
	defer func() {
		err := recover()
		has = err != nil
		if has {
			ok = err == v
		}
	}()

	f()
	return
}

func assertNoPanic(f func()) bool {
	has, _ := assertPanic(nil, f)
	return !has
}

func TestAssertPanic(t *testing.T) {
	f1 := func() { panic("lorem ipsum") }

	if has, ok := assertPanic("lorem ipsum", f1); !has || !ok {
		t.Errorf("f1 SHOULD panic, but got %v, panic result %v", has, ok)
	}

	if has, ok := assertPanic("loremipsum", f1); !has || ok {
		t.Errorf("f1 SHOULD panic, but got %v, panic result %v", has, ok)
	}

	f2 := func() {}

	if has, ok := assertPanic(nil, f2); has {
		t.Errorf("f2 SHOULD not panic, but got %v, panic result %v", has, ok)
	}

	if has, ok := assertPanic("lorem ipsum", f2); has {
		t.Errorf("f2 SHOULD not panic, but got %v, panic result %v", has, ok)
	}
}

func TestAssertNoPanic(t *testing.T) {
	f1 := func() { panic("lorem ipsum") }

	if assertNoPanic(f1) {
		t.Errorf("f1 SHOULD panic, but not.")
	}

	f2 := func() {}

	if !assertNoPanic(f2) {
		t.Errorf("f2 SHOULD NOT panic.")
	}
}
