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
	"github.com/goccmack/gocc/internal/ast"
)

// key: string of symbols - string(ast.CharLit.Lit). E.g.: "'a'"
type CharLitSymbols struct {
	idMap   map[string]int
	typeMap []*ast.LexCharLit
}

func NewCharLitSymbols() *CharLitSymbols {
	return &CharLitSymbols{
		idMap:   make(map[string]int),
		typeMap: make([]*ast.LexCharLit, 0, 16),
	}
}

func (this *CharLitSymbols) Add(cl *ast.LexCharLit) {
	this.typeMap = append(this.typeMap, cl)
	this.idMap[cl.String()] = len(this.typeMap) - 1
}

func (this *CharLitSymbols) GetSymbolId(id string) (sym *ast.LexCharLit, exist bool) {
	if idx, ok := this.idMap[id]; !ok {
		return nil, false
	} else {
		sym = this.typeMap[idx]
	}
	return
}

func (this *CharLitSymbols) Len() int {
	return len(this.typeMap)
}

func (this *CharLitSymbols) List() []*ast.LexCharLit {
	return this.typeMap
}

func (this *CharLitSymbols) StringList() []string {
	symbols := make([]string, len(this.typeMap))
	for i, sym := range this.typeMap {
		symbols[i] = sym.String()
	}
	return symbols
}
