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

// key: string of range, e.g.: 'a'-'z'
type CharRangeSymbols struct {
	idmap   map[string]int
	typeMap []*ast.LexCharRange
}

func NewCharRangeSymbols() *CharRangeSymbols {
	return &CharRangeSymbols{
		idmap:   make(map[string]int),
		typeMap: make([]*ast.LexCharRange, 0, 16),
	}
}

func (this *CharRangeSymbols) Add(cr *ast.LexCharRange) {
	this.typeMap = append(this.typeMap, cr)
	this.idmap[cr.String()] = len(this.typeMap) - 1
}

func (this *CharRangeSymbols) Len() int {
	return len(this.typeMap)
}

func (this *CharRangeSymbols) List() []*ast.LexCharRange {
	return this.typeMap
}

func (this *CharRangeSymbols) StringList() []string {
	symbols := make([]string, len(this.typeMap))
	for i, sym := range this.typeMap {
		symbols[i] = sym.String()
	}
	return symbols
}
