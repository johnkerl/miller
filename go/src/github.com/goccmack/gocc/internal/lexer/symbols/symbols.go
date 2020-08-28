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

package symbols

import (
	"fmt"
	"strings"

	"github.com/goccmack/gocc/internal/ast"
)

type Symbols struct {
	*CharLitSymbols
	*CharRangeSymbols
	ImportIdList []string
	importIdMap  map[string]int
	//key: import Id, value: external function
	importFuncMap map[string]string
	typeMap       []string
	idMap         map[string]int // key: symbol id, val: symbol type
}

func NewSymbols(lexpart *ast.LexPart) (sym *Symbols) {
	sym = &Symbols{
		CharLitSymbols:   NewCharLitSymbols(),
		CharRangeSymbols: NewCharRangeSymbols(),
		ImportIdList:     make([]string, 0, len(lexpart.Imports)+1),
		importIdMap:      make(map[string]int),
		importFuncMap:    make(map[string]string),
	}
	for _, tokDef := range lexpart.TokDefsList {
		tokDef.LexPattern().Walk(sym)
	}
	for _, regDef := range lexpart.RegDefsList {
		regDef.LexPattern().Walk(sym)
	}
	for _, igDef := range lexpart.IgnoredTokDefsList {
		igDef.LexPattern().Walk(sym)
	}
	for id, extFunc := range lexpart.Imports {
		sym.ImportIdList = append(sym.ImportIdList, id)
		sym.importIdMap[id] = len(sym.ImportIdList) - 1
		sym.importFuncMap[id] = extFunc.ExtFunc
	}
	sym.makeMaps()
	return
}

func (this *Symbols) makeMaps() {
	this.typeMap = make([]string, 0, this.CharLitSymbols.Len()+this.CharRangeSymbols.Len()+len(this.ImportIdList))
	this.typeMap = append(this.typeMap, this.CharLitSymbols.StringList()...)
	this.typeMap = append(this.typeMap, this.CharRangeSymbols.StringList()...)
	this.typeMap = append(this.typeMap, this.ImportIdList...)
	this.typeMap = append(this.typeMap, ".")
	this.idMap = make(map[string]int)
	for i, sym := range this.typeMap {
		this.idMap[sym] = i
	}
}

/*
This function returns the external function associated with the import id.
If there is no registered import id the function returns "".
*/
func (this *Symbols) ExternalFunction(id string) string {
	if extFunc, exist := this.importFuncMap[id]; exist {
		return extFunc
	}
	return ""
}

func (this *Symbols) ImportType(id string) int {
	if typ, exist := this.importIdMap[id]; exist {
		return typ
	}
	return -1
}

func (this *Symbols) List() []string {
	return this.typeMap
}

func (this *Symbols) NumSymbols() int {
	return len(this.typeMap)
}

func (this *Symbols) Type(id string) int {
	return this.idMap[id]
}

func (this *Symbols) IsImport(id string) bool {
	if _, isImport := this.importFuncMap[id]; isImport {
		return true
	}
	return false
}

func (this *Symbols) Visit(n ast.LexNode) ast.LexNodeVisitor {
	switch node := n.(type) {
	case *ast.LexCharLit:
		this.CharLitSymbols.Add(node)
	case *ast.LexCharRange:
		this.CharRangeSymbols.Add(node)
	}
	return this
}

func (this *Symbols) String() string {
	buf := new(strings.Builder)
	fmt.Fprintf(buf, "Lexical Symbols = { ")
	for i, sym := range this.List() {
		if i > 0 {
			fmt.Fprintf(buf, ", %s", sym)
		} else {
			fmt.Fprintf(buf, "%s", sym)
		}
	}
	fmt.Fprintf(buf, " }")
	return buf.String()
}
