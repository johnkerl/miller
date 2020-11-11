package main

import (
	"fmt"
	"os"

	"experimental/lexer"
	"experimental/parser"
)

func parseOne(input string) {
	theLexer := lexer.NewLexer([]byte(input))
	theParser := parser.NewParser()
	_, err := theParser.Parse(theLexer)

	green := "\033[32;01m"
	red := "\033[31;01m"
	textdefault := "\033[0m"

	if err != nil {
		//fmt.Println(err)
		fmt.Printf("%sFail%s %s\n", red, textdefault, input)
	} else {
		fmt.Printf("%sOK%s   %s\n", green, textdefault, input)
	}
}

func main() {
	if len(os.Args) == 1 {
		inputs := []string{
			// Expect pass
			"",
			";",
			";;",
			"x",
			"x;x",
			"x;",
			"x;;",
			";x",
			";;x",
			"x ; {}",
			"{} ; x",
			"{} x",
			"{ x }",
			"{ x; x }",
			"x; { x; x }",
			"{ x; x } x",
			"{ x; x } ; x",

			// Expect fail
			"x x",
			"x {}",
		}
		for _, input := range inputs {
			parseOne(input)
		}
	} else {
		for _, arg := range os.Args[1:] {
			parseOne(arg)
		}
	}
}
