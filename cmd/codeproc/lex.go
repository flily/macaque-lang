package main

import (
	"fmt"
	"io"
	"os"

	"github.com/flily/macaque-lang/lex"
	"github.com/flily/macaque-lang/token"
)

func LexicalAnalysis(filename string) {
	fd, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	code, err := io.ReadAll(fd)
	if err != nil {
		panic(err)
	}

	scanner := lex.NewRecursiveScanner(filename)
	scanner.SetContent(code)

	elementList := make([]*token.TokenContext, 0)
	for {
		elem, err := scanner.Scan()
		if err != nil {
			if err != io.EOF {
				panic(err)
			}

			break
		}

		elementList = append(elementList, elem)
	}

	for _, elem := range elementList {
		highlight := elem.ToContext().HighLight()
		fmt.Printf("%s\n%s  %s\n",
			elem.Position.Line.Content,
			highlight,
			elem.Token,
		)
	}
}
