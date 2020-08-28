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

package golang

import (
	"github.com/goccmack/gocc/internal/config"
	"github.com/goccmack/gocc/internal/lexer/items"
	"github.com/goccmack/gocc/internal/token"
)

func Gen(pkg, outDir string, header string, itemsets *items.ItemSets, tokMap *token.TokenMap, cfg config.Config) {
	genLexer(pkg, outDir, itemsets, cfg)
	genTransitionTable(pkg, outDir, header, itemsets)
	genActionTable(pkg, outDir, itemsets, tokMap)
}
