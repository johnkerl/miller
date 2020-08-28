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
	"strings"

	"github.com/goccmack/gocc/internal/ast"
	"github.com/goccmack/gocc/internal/parser/first"
	"github.com/goccmack/gocc/internal/parser/symbols"
)

//A list of a list of Items.
type ItemSets struct {
	sets []*ItemSet
}

// TODO: optimise loop
// g is a BNF grammar. Items returns the sets of Items of the grammar g.
func GetItemSets(g *ast.Grammar, s *symbols.Symbols, firstSets *first.FirstSets) *ItemSets {
	S := &ItemSets{
		sets: []*ItemSet{InitialItemSet(g, s, firstSets).Closure()},
	}
	symbols := s.List()
	included := -1
	for again := true; again; {
		again = false
		for i, I := range S.sets {
			if i > included {
				for _, X := range symbols {
					gto := I.Goto(X)
					if gto.Size() > 0 {
						idx := S.GetIndex(gto)
						if idx == -1 {
							S.sets, again = append(S.sets, gto), true
							idx = len(S.sets) - 1
							gto.SetNo = idx
						}
						I.AddTransition(X, idx)
					}
				}
				included = i
			}
		}
	}
	return S
}

//Returns whether the list of a list of items contains the list of items.
func (this *ItemSets) Contains(I *ItemSet) bool {
	for _, i := range this.sets {
		if i.Equal(I) {
			return true
		}
	}
	return false
}

//Returns the index of the list of items.
func (this *ItemSets) GetIndex(I *ItemSet) int {
	if I == nil || I.Size() == 0 {
		return -1
	}

	for i, items := range this.sets {
		if items.Equal(I) {
			return i
		}
	}
	return -1
}

/*
Return a slice containing all the item sets in increasing order of index
*/
func (this *ItemSets) List() []*ItemSet {
	return this.sets
}

/*
return set[SetNo]
*/
func (this *ItemSets) Set(SetNo int) *ItemSet {
	return this.sets[SetNo]
}

/*
Returns the number of item sets
*/
func (this *ItemSets) Size() int {
	return len(this.sets)
}

//Returns a string representing the list of the list of items.
func (this *ItemSets) String() string {
	buf := new(strings.Builder)
	for i, is := range this.sets {
		fmt.Fprintf(buf, "S%d%s\n", i, is.String())
	}
	return buf.String()
}

//Returns the inital Item of a Grammar.
func InitialItemSet(g *ast.Grammar, symbols *symbols.Symbols, fs *first.FirstSets) *ItemSet {
	set := NewItemSet(symbols, g.SyntaxPart.ProdList, fs)
	set.SetNo = 0
	prod := g.SyntaxPart.ProdList[0]
	set.AddItem(NewItem(0, prod, 0, "$"))
	return set
}
