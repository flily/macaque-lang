package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/flily/macaque-lang/compiler"
	"github.com/flily/macaque-lang/lex"
	"github.com/flily/macaque-lang/parser"
	"github.com/flily/macaque-lang/vm"
)

type Arguments struct {
	CompileMode     bool
	InteractiveMode bool
	Files           []string
}

func readFile(filename string) []byte {
	fd, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	content, err := io.ReadAll(fd)
	if err != nil {
		panic(err)
	}

	return content
}

func execFile(filename string) {
	content := readFile(filename)

	scanner := lex.NewRecursiveScanner(filename)
	scanner.SetContent(content)
	parser := parser.NewLLParser(scanner)
	err := parser.ReadTokens()
	if err != nil {
		fmt.Printf("parse file %s error.\n%s\n", filename, err)
		return
	}

	program, err := parser.Parse()
	if err != nil {
		fmt.Printf("parse file %s error.\n%s\n", filename, err)
		return
	}

	compiler := compiler.NewCompiler()
	page, err := compiler.Compile(program)
	if err != nil {
		fmt.Printf("compile file %s error.\n%s\n", filename, err)
		return
	}

	machine := vm.NewNaiveVM()
	machine.LoadCodePage(page)
	main := page.Main().Func(nil)
	result, err := machine.Run(main)
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

	top := machine.Top()
	if top != nil {
		fmt.Printf("TOP: %s\n", machine.Top().Inspect())
	} else {
		fmt.Println("TOP: nil")
	}
}

func main() {
	args := &Arguments{}

	flag.BoolVar(&args.CompileMode, "c", false, "Compile mode")
	flag.BoolVar(&args.InteractiveMode, "i", false, "Interactive mode")
	flag.Parse()

	if flag.NArg() < 0 {
		fmt.Println("Usage: macaque [-c] [-i] <file>")
		return
	}

	args.Files = flag.Args()
	if len(args.Files) <= 0 {
		args.InteractiveMode = true
	}

	if args.InteractiveMode {
		Repl(args)
	} else {
		execFile(args.Files[0])
	}
}
