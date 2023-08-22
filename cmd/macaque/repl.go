package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/flily/macaque-lang/compiler"
	"github.com/flily/macaque-lang/vm"
)

func Repl(args *Arguments) {
	m := vm.NewNaiveVMInterpreter()
	var cc *compiler.Compiler

	if len(args.Files) > 0 {
		filename := args.Files[0]
		c, page, err := compiler.CompileFile(filename)
		if err != nil {
			fmt.Printf("compile file %s error.\n%s\n", filename, err)
			return
		}

		cc = c
		m.LoadCodePage(page)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">>> ")
		input, err := reader.ReadSlice('\n')
		if err != nil {
			fmt.Printf("read input error.\n%s\n", err)
			return
		}

		code, err := cc.CompileCode("stdin", input)
		if err != nil {
			fmt.Printf("compile input error.\n%s\n", err)
			return
		}

		if m.CodePage != nil {
			m.MergeCodeBlock(code, cc.Context)

		} else {
			page := cc.Link(code)

			m.LoadCodePage(page)
			main := page.Main().Func(nil)
			m.SetEntry(main)
		}

		page := m.CodePage
		main := page.Main()
		main.Relink()

		result, err := m.Resume(main.Func(nil))
		if err != nil {
			fmt.Printf("runtime error.\n%s\n", err)
			return
		}

		if len(result) > 0 {
			fmt.Printf("##> ")
			for _, r := range result {
				fmt.Printf("%s ", r.Inspect())
			}
			fmt.Println()
		}

		top := m.Top()
		if top != nil {
			fmt.Printf("STACK:\n")
			sp := int(m.GetSP())
			for i := sp - 1; i >= 0; i-- {
				fmt.Printf("  - %s\n", m.GetStackObject(i).Inspect())
			}

		} else {
			fmt.Println("TOP: nil")
		}
	}
}
