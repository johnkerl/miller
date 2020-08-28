package dsl

import (
	"fmt"
	"testing"

	"miller/dsl"
	"miller/parsing/lexer"
	"miller/parsing/parser"
)

func testOne(sourceString []byte) (*dsl.AST, error) {
	fmt.Printf("Input: %s\n", sourceString)
	theLexer := lexer.NewLexer(sourceString)
	theParser := parser.NewParser()
	interfaceAST, err := theParser.Parse(theLexer)
	if err == nil {
		return interfaceAST.(*dsl.AST), nil
	} else {
		return nil, err
	}
}

func TestPass(t *testing.T) {
	ast, err := testOne([]byte("$x = 3"))
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("AST:")
	ast.Print()
}

func TestFail(t *testing.T) {
	_, err := testOne([]byte("a b ; d e f"))
	if err == nil {
		t.Fatal("Expected parse error")
	} else {
		fmt.Printf("Parsing failed as expected: %v\n", err)
	}
}
