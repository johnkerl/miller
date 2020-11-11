package main

import (
	"fmt"
	"os"

	"dsl"
	"experimental/lexer"
	"experimental/parser"
)

func parseOne(input string) {
	theLexer := lexer.NewLexer([]byte(input))
	theParser := parser.NewParser()
	iast, err := theParser.Parse(theLexer)

	green := "\033[32;01m"
	red := "\033[31;01m"
	textdefault := "\033[0m"

	if err == nil {
		fmt.Printf("%sOK%s   %s\n", green, textdefault, input)
		iast.(*dsl.AST).Print()
	} else {
		//fmt.Println(err)
		fmt.Printf("%sFail%s %s\n", red, textdefault, input)
	}
	fmt.Println()
}

func main() {
	if len(os.Args) == 1 {

		fmt.Println("EXPECT OK")
		goods := []string{
			"",
			";",
			";;",
			"x",
			"x;x",
			"x;x;x",
			"x;x;x;x",
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
			"{};{}",
			"{} {}",
		}
		for _, input := range goods {
			parseOne(input)
		}

		fmt.Println()
		fmt.Println("EXPECT FAIL")
		bads := []string{
			"x x",
			"x {}",
		}
		for _, input := range bads {
			parseOne(input)
		}

	} else {
		for _, arg := range os.Args[1:] {
			parseOne(arg)
		}
	}
}
