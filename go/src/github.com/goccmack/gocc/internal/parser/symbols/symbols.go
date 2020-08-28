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
Support for the symbols of the language defined by the input grammar, G. This package supports code generation.
*/
package symbols

import (
	"fmt"
	"strings"

	"github.com/goccmack/gocc/internal/ast"
)

type Symbols struct {
	//key: symbol id
	//val: symbol type
	idMap map[string]int

	//key: symbol ntTypeMap index
	//val: symbol type
	ntIdMap   map[string]int
	ntTypeMap []string

	//key: symbol id
	//val: symbol type
	stringLitIdMap map[string]int
	stringLitList  []string

	//key: symbol type
	//val: symbol id
	typeMap []string
}

func NewSymbols(grammar *ast.Grammar) *Symbols {
	symbols := &Symbols{
		idMap:          make(map[string]int),
		typeMap:        make([]string, 0, 16),
		ntIdMap:        make(map[string]int),
		ntTypeMap:      make([]string, 0, 16),
		stringLitIdMap: make(map[string]int),
		stringLitList:  make([]string, 0, 16),
	}

	symbols.Add("INVALID")
	symbols.Add("$")

	if grammar.SyntaxPart == nil {
		return symbols
	}

	for _, p := range grammar.SyntaxPart.ProdList {
		if _, exist := symbols.ntIdMap[p.Id]; !exist {
			symbols.ntTypeMap = append(symbols.ntTypeMap, p.Id)
			symbols.ntIdMap[p.Id] = len(symbols.ntTypeMap) - 1
		}
		symbols.Add(p.Id)
		for _, sym := range p.Body.Symbols {
			symStr := sym.SymbolString()
			symbols.Add(symStr)
			if _, ok := sym.(ast.SyntaxStringLit); ok {
				if _, exist := symbols.ntIdMap[symStr]; exist {
					panic(fmt.Sprintf("string_lit \"%s\" conflicts with production name %s", symStr, symStr))
				}
				if _, exist := symbols.stringLitIdMap[symStr]; !exist {
					symbols.stringLitIdMap[symStr] = symbols.Type(symStr)
					symbols.stringLitList = append(symbols.stringLitList, symStr)
				}
			}
		}
	}
	return symbols
}

func (this *Symbols) Add(symbols ...string) {
	for _, sym := range symbols {
		if _, exist := this.idMap[sym]; !exist {
			this.typeMap = append(this.typeMap, sym)
			this.idMap[sym] = len(this.typeMap) - 1
		}
	}
}

func (this *Symbols) Id(typ int) string {
	return this.typeMap[typ]
}

func (this *Symbols) IsTerminal(sym string) bool {
	_, nt := this.ntIdMap[sym]
	return !nt
}

func (this *Symbols) List() []string {
	return this.typeMap
}

/*
Return a slice containing the ids of all symbols declared as string literals in the grammar.
*/
func (this *Symbols) ListStringLitSymbols() []string {
	return this.stringLitList
}

func (this *Symbols) ListTerminals() []string {
	terminals := make([]string, 0, 16)
	for _, sym := range this.typeMap {
		if this.IsTerminal(sym) {
			terminals = append(terminals, sym)
		}
	}
	return terminals
}

func (this *Symbols) StringLitType(id string) int {
	if typ, exist := this.stringLitIdMap[id]; exist {
		return typ
	}
	return -1
}

/*
Return the id of the NT with index idx, or "" if there is no NT symbol with index, idx.
*/
func (this *Symbols) NTId(idx int) string {
	if idx < 0 || idx >= len(this.ntTypeMap) {
		return ""
	}
	return this.ntTypeMap[idx]
}

/*
Return the number of NT symbols in the grammar
*/
func (this *Symbols) NumNTSymbols() int {
	return len(this.ntTypeMap)
}

/*
Returns a slice containing all the non-terminal symbols of the grammar.
*/
func (this *Symbols) NTList() []string {
	return this.ntTypeMap
}

/*
Returns the NT index of a symbol (index in 0..|NT|-1) or -1 if the symbol is not in NT.
*/
func (this *Symbols) NTType(symbol string) int {
	if idx, exist := this.ntIdMap[symbol]; exist {
		return idx
	}
	return -1
}

/*
Returns the total number of symbols in grammar: the sum of the terminals and non-terminals.
*/
func (this *Symbols) NumSymbols() int {
	return len(this.typeMap)
}

func (this *Symbols) String() string {
	w := new(strings.Builder)
	for i, sym := range this.typeMap {
		fmt.Fprintf(w, "%3d: %s\n", i, sym)
	}
	return w.String()
}

func (this *Symbols) Type(id string) int {
	if typ, ok := this.idMap[id]; ok {
		return typ
	}
	return -1
}
