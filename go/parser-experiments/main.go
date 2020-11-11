package main

import (
	"fmt"
	"os"

	"experimental/lexer"
	"experimental/parser"
)

func parseOne(input string) {
	fmt.Println("----------------------------------------------------------------")
	fmt.Printf("Parsing \"%s\"\n", input)
	theLexer := lexer.NewLexer([]byte(input))
	theParser := parser.NewParser()
	_, err := theParser.Parse(theLexer)
	if err != nil {
		//fmt.Println(err)
		fmt.Println("Fail")
	} else {
		fmt.Println("OK")
	}
}

func main() {
	if len(os.Args) == 1 {
		inputs := []string{
			"",
			";",
			"x",
			"x;x",
			"x;",
			";x",
			"x ; {}",
			"{} ; x",
			"{} x",
		}
		for i, input := range inputs {
			if i > 0 {
				fmt.Println()
			}
			parseOne(input)
		}
	} else {
		for i, arg := range os.Args[1:] {
			if i > 0 {
				fmt.Println()
			}
			parseOne(arg)
		}
	}
	fmt.Println("----------------------------------------------------------------")
}
