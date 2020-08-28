package t1

import (
	"testing"

	"github.com/goccmack/gocc/internal/test/t1/errors"
	"github.com/goccmack/gocc/internal/test/t1/lexer"
	"github.com/goccmack/gocc/internal/test/t1/parser"
)

func Test1(t *testing.T) {
	ast, err := parser.NewParser().Parse(lexer.NewLexer([]byte(`c`)))
	if err != nil {
		t.Fail()
	}
	if slice, ok := ast.([]interface{}); !ok {
		t.Fail()
	} else {
		if len(slice) != 2 {
			t.Fatalf("len(slice)==%d", len(slice))
		}
		if slice[0] != nil {
			t.Fatal("slice[0] != nil")
		}
		if str, ok := slice[1].(string); !ok {
			t.Fatalf("%T", slice[1])
		} else {
			if str != "c" {
				t.Fatal(`str != "c"`)
			}
		}
	}
}

func Test2(t *testing.T) {
	ast, err := parser.NewParser().Parse(lexer.NewLexer([]byte(`b c`)))
	if err != nil {
		t.Fatal(err)
	}
	if slice, ok := ast.([]interface{}); !ok {
		t.Fail()
	} else {
		if len(slice) != 2 {
			t.Fatalf("len(slice)==%d", len(slice))
		}
		if str, ok := slice[0].(string); !ok {
			t.Fatal(`str, ok := slice[0].(string); !ok`)
		} else {
			if str != "b" {
				t.Fatal(`str != "b"`)
			}
		}
		if str, ok := slice[1].(string); !ok {
			t.Fatalf("%T", slice[1])
		} else {
			if str != "c" {
				t.Fatal(`str != "c"`)
			}
		}
	}
}

func Test3(t *testing.T) {
	toks := lexer.NewLexer([]byte(`c b`))
	_, err := parser.NewParser().Parse(toks)
	if err == nil {
		t.Fatal("No error for erronous input.")
	}
	if errs, ok := err.(*errors.Error); !ok {
		t.Fatal("Incompatible error type for erronous input.")
	} else {
		if errs.ErrorToken.Column != 3 {
			t.Fatal("errs.ErrorToken.Column = ", errs.ErrorToken.Column, " != 3")
		}
	}
}
