package nolexer

import (
	"testing"

	"github.com/goccmack/gocc/example/nolexer/parser"
	"github.com/goccmack/gocc/example/nolexer/scanner"
)

func Test1(t *testing.T) {
	S := scanner.NewString("hiya world")
	P := parser.NewParser()
	if _, e := P.Parse(S); e != nil {
		t.Error(e.Error())
	}
}

func Test2(t *testing.T) {
	S := scanner.NewString("hello world")
	P := parser.NewParser()
	if _, e := P.Parse(S); e != nil {
		t.Error(e.Error())
	}
}
