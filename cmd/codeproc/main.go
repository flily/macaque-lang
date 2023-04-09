package main

import (
	"flag"
	"fmt"
)

func doHelp() {
	fmt.Printf("Usage: codeproc [COMMAND]\n")
}

func doLex(args []string) {
	for _, filename := range args {
		LexicalAnalysis(filename)
	}
}

func main() {
	flag.Parse()

	if flag.NArg() <= 0 {
		doHelp()
		return
	}

	cmd := flag.Args()[0]
	args := flag.Args()[1:]
	switch cmd {
	case "help":
		doHelp()

	case "lex":
		doLex(args)
	}
}
