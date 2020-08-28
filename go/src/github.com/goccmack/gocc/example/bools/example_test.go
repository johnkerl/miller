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

package example

import (
	"testing"

	"github.com/goccmack/gocc/example/bools/ast"
	"github.com/goccmack/gocc/example/bools/lexer"
	"github.com/goccmack/gocc/example/bools/parser"
)

func testEval(t *testing.T, exampleStr string, output bool) {
	lex := lexer.NewLexer([]byte(exampleStr))
	p := parser.NewParser()
	st, err := p.Parse(lex)
	if err != nil {
		panic(err)
	}
	if st.(ast.Val).Eval() != output {
		t.Fatalf("Should be %v for %v", output, exampleStr)
	}
}

func TestOr(t *testing.T) {
	testEval(t, "true | false", true)
}

func TestAnd(t *testing.T) {
	testEval(t, "true & false", false)
}

func TestSubString(t *testing.T) {
	testEval(t, `"true" in "false"`, false)
	testEval(t, `"true" in "trues"`, true)
}

func TestLess(t *testing.T) {
	testEval(t, "0 < 5", true)
	testEval(t, "0 > 5", false)
}

func TestMixed(t *testing.T) {
	testEval(t, "0 < 5 | false", true)
}

func TestGroup(t *testing.T) {
	testEval(t, "( true | false ) & ( true & true )", true)
}

func TestGroupMixed(t *testing.T) {
	testEval(t, `( true | false ) & 0 > 100000 | "t" in "taddle"`, true)
}
