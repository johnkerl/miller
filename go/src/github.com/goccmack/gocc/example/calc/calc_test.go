package calc

import (
	"testing"

	"github.com/goccmack/gocc/example/calc/lexer"
	"github.com/goccmack/gocc/example/calc/parser"
)

type TI struct {
	src    string
	expect int64
}

var testData = []*TI{
	{"1 + 1", 2},
	{"1 * 1", 1},
	{"1 + 2 * 3", 7},
}

func Test1(t *testing.T) {
	p := parser.NewParser()
	for _, ts := range testData {
		s := lexer.NewLexer([]byte(ts.src))
		sum, err := p.Parse(s)
		if err != nil {
			t.Error(err)
		}
		if sum != ts.expect {
			t.Errorf("Error: %s = %d. Expected %d\n", ts.src, sum, ts.expect)
		}
	}
}
