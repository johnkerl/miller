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
	"github.com/goccmack/gocc/internal/lexer/symbols"
)

type ItemSets struct {
	sets    []*ItemSet
	lexPart *ast.LexPart
	symbols *symbols.Symbols
}

func GetItemSets(lexPart *ast.LexPart) *ItemSets {
	itemSets := &ItemSets{
		sets:    make([]*ItemSet, 0, 256),
		lexPart: lexPart,
		symbols: symbols.NewSymbols(lexPart),
	}

	itemSets.Add(ItemsSet0(lexPart, itemSets.symbols))

	return itemSets.Closure()
}

func (this *ItemSets) Add(items ItemList) (setNo int) {
	if yes, setNo := this.Contain(items); yes {
		return setNo
	}
	setNo = this.Size()
	this.sets = append(this.sets, NewItemSet(setNo, this.lexPart, this.symbols, items))
	return setNo
}

func (this *ItemSets) Closure() *ItemSets {
	for i := 0; i < len(this.sets); i++ {
		for symI, rng := range this.sets[i].SymbolClasses.List() {
			if items := this.sets[i].Next(rng); len(items) != 0 {
				setNo, nextState := this.Add(items), this.sets[i].Transitions[symI]
				if nextState != -1 && nextState != setNo {
					panic(fmt.Sprintf("Next set conflict in (S%d, %s) -> %d. Existing setNo: %d", i, rng, setNo, nextState))
				}
				this.sets[i].Transitions[symI] = setNo
			}
		}
		for impI, imprt := range this.symbols.ImportIdList {
			if items := this.sets[i].NextImport(imprt); len(items) != 0 {
				setNo, nextState := this.Add(items), this.sets[i].ImportTransitions[this.symbols.ImportType(imprt)]
				if nextState != -1 && nextState != setNo {
					panic(fmt.Sprintf("Next set conflict in (S%d, %s) -> %d. Existing setNo: %d", i, imprt, setNo, nextState))
				}
				this.sets[i].ImportTransitions[impI] = setNo
			}
		}
		if items := this.sets[i].NextDot(); len(items) != 0 {
			setNo := this.Add(items)
			this.sets[i].DotTransition = setNo
		}
	}
	return this
}

func (this *ItemSets) Contain(items ItemList) (yes bool, index int) {
	for i, thisSet := range this.sets {
		if thisSet.Equal(items) {
			return true, i
		}
	}
	return false, -1
}

func (this *ItemSets) Size() int {
	return len(this.sets)
}

func (this *ItemSets) List() []*ItemSet {
	return this.sets
}

func (this *ItemSets) Symbols() *symbols.Symbols {
	return this.symbols
}

func (this *ItemSets) String() string {
	buf := new(strings.Builder)
	fmt.Fprintf(buf, "Item sets:\n")
	for i, s := range this.sets {
		fmt.Fprintf(buf, "S%d{\n", i)
		for _, item := range s.Items {
			fmt.Fprintf(buf, "\t%s\n", item)
		}
		fmt.Fprintf(buf, "}\n")
		fmt.Fprintf(buf, "Transitions:\n")
		symClasses := s.SymbolClasses.List()
		for idx, symClass := range symClasses {
			fmt.Fprintf(buf, "\t%s -> S%d\n", symClass, s.Transitions[idx])
		}
		for idx, imprt := range this.symbols.ImportIdList {
			fmt.Fprintf(buf, "\t%s -> %d\n", imprt, s.ImportTransitions[idx])
		}
		if s.SymbolClasses.MatchAny {
			fmt.Fprintf(buf, ". -> S%d\n", s.DotTransition)
		}
		if s.Action() == nil {
			fmt.Fprintf(buf, "Action: nil\n")
		} else {
			fmt.Fprintf(buf, "Action: %s\n", s.Action())
		}
		fmt.Fprintf(buf, "Symbols classes: %s\n\n", s.SymbolClasses)
	}
	return buf.String()
}
