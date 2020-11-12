package main

import (
	"fmt"
	"os"

	"dsl"
	"experimental/lexer"
	"experimental/parser"
)

const GREEN = "\033[32;01m"
const RED = "\033[31;01m"
const TEXTDEFAULT = "\033[0m"

func parseOne(input string) bool {
	theLexer := lexer.NewLexer([]byte(input))
	theParser := parser.NewParser()
	iast, err := theParser.Parse(theLexer)
	if err == nil {
		fmt.Printf("%sOK%s   %s\n", GREEN, TEXTDEFAULT, input)
		iast.(*dsl.AST).Print()
		fmt.Println()
		return true
	} else {
		//fmt.Println(err)
		fmt.Printf("%sFail%s %s\n", RED, TEXTDEFAULT, input)
		fmt.Println()
		return false
	}
}

func main() {
	if len(os.Args) == 1 {
		ok := true

		fmt.Println("----------------------------------------------------------------")
		fmt.Println("EXPECT OK")
		goods := []string{
			"",
			";",
			";;",
			";;;",
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
			"{} {};",
			"{};{};",
			"{} {} {}",
			"{};{}; {}",
			"{} ; ; {}",
			"x; x;",
			"{x; x;}",
			"{x; x;} x",
			"{x; x;} x;",
			"{x; x;} {};",
			"{} x; {}",
			"{} x; {};",
			"{} x; x; {}",
			"{} x; x; {};",
			"{} x; x; x; {}",
			"{} x; x; x; {};",
			"{} {} ;;; {;;} x; x; x; x; x; {}",
		}
		for _, input := range goods {
			if parseOne(input) == false {
				ok = false
			}
		}

		fmt.Println()
		fmt.Println("----------------------------------------------------------------")
		fmt.Println("EXPECT FAIL")
		bads := []string{
			"x x",
			"x {}",
		}
		for _, input := range bads {
			if parseOne(input) == true {
				ok = false
			}
		}

		fmt.Println()
		fmt.Println("----------------------------------------------------------------")
		if ok {
			fmt.Printf("%sALL AS EXPECTED%s\n", GREEN, TEXTDEFAULT)
		} else {
			fmt.Printf("%sNOT ALL AS EXPECTED%s\n", RED, TEXTDEFAULT)
			os.Exit(1)
		}

	} else {
		for _, arg := range os.Args[1:] {
			parseOne(arg)
		}
	}
}
