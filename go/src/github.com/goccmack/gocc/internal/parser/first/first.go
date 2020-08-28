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

package first

import (
	"fmt"
	"strings"

	"github.com/goccmack/gocc/internal/ast"
	"github.com/goccmack/gocc/internal/parser/symbols"
)

/*
key: Id of production
*/
type FirstSets struct {
	firstSets map[string]SymbolSet
	symbols   *symbols.Symbols
}

const EMPTY = "empty"

//Returns the FirstSets of the Grammar.
func GetFirstSets(g *ast.Grammar, symbols *symbols.Symbols) *FirstSets {
	firstSets := &FirstSets{
		firstSets: make(map[string]SymbolSet),
		symbols:   symbols,
	}

	if g.SyntaxPart == nil {
		return firstSets
	}

	for again := true; again; {
		again = false
		for _, prod := range g.SyntaxPart.ProdList {
			switch {
			case len(prod.Body.Symbols) == 0:
				if firstSets.AddToken(prod.Id, EMPTY) {
					again = true
				}
			case symbols.IsTerminal(prod.Body.Symbols[0].SymbolString()):
				if firstSets.AddToken(prod.Id, prod.Body.Symbols[0].SymbolString()) {
					again = true
				}
			default:
				first := FirstS(firstSets, stringList(prod.Body.Symbols))
				if !first.Equal(firstSets.GetSet(prod.Id)) {
					if firstSets.AddSet(prod.Id, first) {
						again = true
					}
				}
			}
		}
	}

	return firstSets
}

func stringList(symbols []ast.SyntaxSymbol) []string {
	sl := make([]string, len(symbols))
	for i, sym := range symbols {
		sl[i] = sym.SymbolString()
	}
	return sl
}

func (this *FirstSets) AddSet(prodName string, terminals SymbolSet) (symbolsAdded bool) {
	for symbol := range terminals {
		if this.AddToken(prodName, symbol) {
			symbolsAdded = true
		}
	}
	return
}

func (this *FirstSets) AddToken(prodName string, terminal string) (symbolAdded bool) {
	set, ok := this.firstSets[prodName]
	if !ok {
		set = make(SymbolSet)
		this.firstSets[prodName] = set
	}
	if _, contain := set[terminal]; !contain {
		set[terminal] = true
		symbolAdded = true
	}
	return
}

//Returns a set.
func (this *FirstSets) GetSet(prodName string) SymbolSet {
	if set, ok := this.firstSets[prodName]; ok {
		return set
	}
	return nil
}

//Returns a string representing the FirstSets.
func (this *FirstSets) String() string {
	buf := new(strings.Builder)
	for _, nt := range this.symbols.NTList() {
		set := this.firstSets[nt]
		fmt.Fprintf(buf, "%s: %s\n", nt, set)
	}
	return buf.String()
}

//Returns the First SymbolSet given the ast.SyntaxSymbol.
func First(fs *FirstSets, sym string) SymbolSet {
	if fs.symbols.IsTerminal(sym) {
		return SymbolSet{sym: true}
	}
	return fs.GetSet(sym)
}

/*
Returns First of the string, xyz, e.g.: for the item,

  X  : w • xyz

  Let x, y, z be strings from the union of T and NT.
  First(xy...z) =
	First(x) if First(x) does not contain ϵ
 	First(x) + First(y) if First(x) contains ϵ but First(y) does not contain ϵ
 	...
 	First(x) + First(y) + ... + First(z)
*/
func FirstS(firstSets *FirstSets, symbols []string) (first SymbolSet) {
	first = make(SymbolSet)
	if len(symbols) == 0 {
		return
	}
	fst := First(firstSets, symbols[0])
	first.AddSet(fst)
	_, containEmpty := fst["empty"]
	for i := 1; i < len(symbols) && containEmpty; i++ {
		fst = First(firstSets, symbols[i])
		first.AddSet(fst)
		_, containEmpty = fst["empty"]
	}
	if !containEmpty {
		delete(first, "empty")
	}
	return
}
