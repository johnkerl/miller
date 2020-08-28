package dsl

import (
	"fmt"
	"testing"

	"miller/dsl"
	"miller/parsing/lexer"
	"miller/parsing/parser"
)

func testOne(src []byte) (astree dsl.StatementList, err error) {
	fmt.Printf("Input: %s\n", src)
	s := lexer.NewLexer(src)
	p := parser.NewParser()
	a, err := p.Parse(s)
	if err == nil {
		astree = a.(dsl.StatementList)
	}
	return
}

func TestPass(t *testing.T) {
	sml, err := testOne([]byte("$x = 3"))
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Printf("Output: %s\n", sml)
}

func TestFail(t *testing.T) {
	_, err := testOne([]byte("a b ; d e f"))
	if err == nil {
		t.Fatal("Expected parse error")
	} else {
		fmt.Printf("Parsing failed as expected: %v\n", err)
	}
}
