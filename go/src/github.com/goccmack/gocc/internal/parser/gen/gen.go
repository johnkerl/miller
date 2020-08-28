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

/*
This package controls the generation of all parser-related code.
*/
package gen

import (
	"github.com/goccmack/gocc/internal/ast"
	"github.com/goccmack/gocc/internal/config"
	"github.com/goccmack/gocc/internal/parser/gen/golang"
	"github.com/goccmack/gocc/internal/parser/lr1/items"
	"github.com/goccmack/gocc/internal/parser/symbols"
	"github.com/goccmack/gocc/internal/token"
)

func Gen(pkg, outDir, header string, prods ast.SyntaxProdList, symbols *symbols.Symbols,
	itemsets *items.ItemSets, tokMap *token.TokenMap, cfg config.Config) (conflicts map[int]items.RowConflicts) {

	golang.GenAction(outDir)
	conflicts = golang.GenActionTable(outDir, prods, itemsets, tokMap, cfg.Zip())
	golang.GenErrors(pkg, outDir)
	golang.GenGotoTable(outDir, itemsets, symbols, cfg.Zip())
	golang.GenParser(pkg, outDir, prods, itemsets, symbols, cfg)
	golang.GenProductionsTable(pkg, outDir, header, prods, symbols, itemsets, tokMap)

	return
}
