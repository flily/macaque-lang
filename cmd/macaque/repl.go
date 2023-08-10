package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/flily/macaque-lang/compiler"
	"github.com/flily/macaque-lang/lex"
	"github.com/flily/macaque-lang/parser"
	"github.com/flily/macaque-lang/vm"
)

func Repl(args *Arguments) {
	m := vm.NewNaiveVMInterpreter()
	compiler := compiler.NewCompiler()

	if len(args.Files) > 0 {
		filename := args.Files[0]
		scanner := lex.NewRecursiveScanner(filename)
		parser := parser.NewLLParser(scanner)
		if err := parser.ReadTokens(); err != nil {
			fmt.Printf("parse file %s error.\n%s\n", filename, err)
			return
		}

		program, err := parser.Parse()
		if err != nil {
			fmt.Printf("parse file %s error.\n%s\n", filename, err)
			return
		}

		page, err := compiler.Compile(program)
		if err != nil {
			fmt.Printf("compile file %s error.\n%s\n", filename, err)
			return
		}

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

		scanner := lex.NewRecursiveScanner("stdin")
		scanner.SetContent([]byte(input))

		parser := parser.NewLLParser(scanner)
		err = parser.ReadTokens()
		if err != nil {
			fmt.Printf("parse input error.\n%s\n", err)
			return
		}

		line, err := parser.Parse()
		if err != nil {
			fmt.Printf("parse input error.\n%s\n", err)
			return
		}

		code, err := compiler.CompileCode(line)
		if err != nil {
			fmt.Printf("compile input error.\n%s\n", err)
			return
		}

		if m.CodePage != nil {
			m.MergeCodeBlock(code, compiler.Context)

		} else {
			page := compiler.Link(code)

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
