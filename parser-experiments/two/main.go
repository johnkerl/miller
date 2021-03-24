package main

import (
	"fmt"
	"os"

	"miller/src/dsl"

	"miller/src/parsing/lexer"
	"miller/src/parsing/parser"
)

const GREEN = "\033[32;01m"
const RED = "\033[31;01m"
const TEXTDEFAULT = "\033[0m"

func parseOne(input string, printError bool) bool {
	theLexer := lexer.NewLexer([]byte(input))
	theParser := parser.NewParser()
	iast, err := theParser.Parse(theLexer)
	if err == nil {
		fmt.Printf("%sOK%s   %s\n", GREEN, TEXTDEFAULT, input)
		iast.(*dsl.AST).Print()
		fmt.Println()
		return true
	} else {
		if printError {
			fmt.Println(err)
		}
		fmt.Printf("%sFail%s %s\n", RED, TEXTDEFAULT, input)
		fmt.Println()
		return false
	}
}

func main() {
	printError := false
	args := os.Args[1:] // os.Args[0] is program name in Go
	if len(args) >= 1 && args[0] == "-v" {
		printError = true
		args = args[1:]
	}

	if len(args) == 0 {
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
			if parseOne(input, printError) == false {
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
			if parseOne(input, printError) == true {
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
		for _, arg := range args {
			parseOne(arg, printError)
		}
	}
}
