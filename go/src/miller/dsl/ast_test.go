package dsl

import (
	"fmt"
	"testing"

	"miller/dsl/ast"
	"miller/parsing/lexer"
	"miller/parsing/parser"
)

func TestPass(t *testing.T) {
	sml, err := test([]byte("$x = 3"))
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Printf("Output: %s\n", sml)
}

func TestFail(t *testing.T) {
	_, err := test([]byte("a b ; d e f"))
	if err == nil {
		t.Fatal("Expected parse error")
	} else {
		fmt.Printf("Parsing failed as expected: %v\n", err)
	}
}

func test(src []byte) (astree ast.StatementList, err error) {
	fmt.Printf("input: %s\n", src)
	s := lexer.NewLexer(src)
	p := parser.NewParser()
	a, err := p.Parse(s)
	if err == nil {
		astree = a.(ast.StatementList)
	}
	return
}
