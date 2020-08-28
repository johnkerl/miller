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

package ast

import (
	"fmt"
)

type LexProdMap struct {
	//key: production id
	idMap map[string]LexProdIndex

	//value: lex production id
	idxMap map[LexProdIndex]string
}

type LexProdIndex int

func NewLexProdMap(prodList *LexProductions) *LexProdMap {
	lpm := &LexProdMap{
		idMap:  make(map[string]LexProdIndex),
		idxMap: make(map[LexProdIndex]string),
	}
	lpm.Add(prodList.Productions...)

	return lpm
}

func newLexProdMap() *LexProdMap {
	return &LexProdMap{
		idMap:  make(map[string]LexProdIndex),
		idxMap: make(map[LexProdIndex]string),
	}
}

func (this *LexProdMap) Index(id string) LexProdIndex {
	idx, exist := this.idMap[id]
	if exist {
		return idx
	}
	return -1
}

func (this *LexProdMap) Id(index LexProdIndex) string {
	id, exist := this.idxMap[index]
	if exist {
		return id
	}
	return ""
}

func (this *LexProdMap) Add(prods ...LexProduction) {
	for _, prod := range prods {
		if _, exist := this.idMap[prod.Id()]; exist {
			panic(fmt.Sprintf("Production %s already exists", prod.Id()))
		}
		idx := len(this.idxMap)
		this.idMap[prod.Id()] = LexProdIndex(idx)
		this.idxMap[LexProdIndex(idx)] = prod.Id()
	}
}
