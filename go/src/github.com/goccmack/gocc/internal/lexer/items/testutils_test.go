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
	"fmt"
	"testing"

	"github.com/goccmack/gocc/internal/ast"
	"github.com/goccmack/gocc/internal/frontend/parser"
	"github.com/goccmack/gocc/internal/frontend/scanner"
	"github.com/goccmack/gocc/internal/frontend/token"
)

func findSet(sets *ItemSets, items []string) *ItemSet {
	for _, set := range sets.List() {
		if setEqualItems(set, items) {
			return set
		}
	}
	return nil
}

func setEqualItems(set *ItemSet, items []string) bool {
	if set.Size() != len(items) {
		return false
	}
	for _, item := range items {
		if !setContainItem(set, item) {
			return false
		}
	}
	return true
}

func setContainItem(set *ItemSet, item string) bool {
	for _, setItem := range set.Items {
		if setItem.str == item {
			return true
		}
	}
	return false
}

func parse(src string, t *testing.T) *ast.Grammar {
	scanner := new(scanner.Scanner)
	scanner.Init([]byte(src), token.FRONTENDTokens)
	parser := parser.NewParser(parser.ActionTable, parser.GotoTable, parser.ProductionsTable, token.FRONTENDTokens)
	g, err := parser.Parse(scanner)
	if err != nil {
		t.Fatal(fmt.Sprintf("Parse error: %s\n", err))
		return nil
	}

	return g.(*ast.Grammar)
}
