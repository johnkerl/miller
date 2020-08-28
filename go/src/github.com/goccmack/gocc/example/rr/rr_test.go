//Copyright 2012 Vastech SA (PTY) LTD
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package rr

import (
	"fmt"
	"testing"

	"github.com/goccmack/gocc/example/rr/lexer"
	"github.com/goccmack/gocc/example/rr/parser"
)

func parse(src string) (ast string, err error) {
	l := lexer.NewLexer([]byte(src))
	p := parser.NewParser()
	res, err := p.Parse(l)
	if err == nil {
		ast = res.(string)
	}
	return
}

func test(t *testing.T, src, exp string) {
	ast, err := parse(src)
	if err != nil {
		t.Fatalf("\tError: %s\n", err.Error())
	}
	if ast != exp {
		t.Fatalf("\tError: ast= `%s`\n", ast)
	}
}

type TD struct {
	src string
	exp string
}

var testData = []TD{
	{"a", "B "},
	{"a a", "A1 "},
	{"a a a", "A1 "},
	{"c a", "A1 "},
	{"c a a a a", "A1 "},
}

func Test(t *testing.T) {
	for _, td := range testData {
		fmt.Printf("\tsrc: `%s`; exp: `%s`\n", td.src, td.exp)
		test(t, td.src, td.exp)
	}
}
