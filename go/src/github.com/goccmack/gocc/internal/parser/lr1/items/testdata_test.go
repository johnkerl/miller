//Copyright 2013 Vastech SA (PTY) LTD
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

package items

import (
	"testing"

	"github.com/goccmack/gocc/internal/ast"
	"github.com/goccmack/gocc/internal/frontend/parser"
	"github.com/goccmack/gocc/internal/frontend/scanner"
	"github.com/goccmack/gocc/internal/frontend/token"
)

const G1 = `
A : B a ;

B : b ;
`

var G1S0 = []string{
	"A : •B a «$»",
	"B : •b «$»",
}

func parse(src string, t *testing.T) *ast.Grammar {
	lexer := new(scanner.Scanner)
	lexer.Init([]byte(src), token.FRONTENDTokens)
	p := parser.NewParser(parser.ActionTable, parser.GotoTable, parser.ProductionsTable, token.FRONTENDTokens)
	res, err := p.Parse(lexer)
	if err != nil {
		t.Fatal(err)
		return nil
	}
	return res.(*ast.Grammar)
}
