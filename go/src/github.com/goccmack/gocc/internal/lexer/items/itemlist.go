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
	"errors"
	"fmt"
	"strings"

	"github.com/goccmack/gocc/internal/ast"
	"github.com/goccmack/gocc/internal/lexer/symbols"
)

// Each Itemset element is a ItemList
type ItemList []*Item

func NewItemList(len int) ItemList {
	if len < 8 {
		len = 8
	}
	return make(ItemList, 0, len)
}

func (this ItemList) AddExclusive(item *Item) (ItemList, error) {
	if this.Contain(item) {
		return nil, errors.New(fmt.Sprintf("Duplicate item: %s in set:\n%s", item.String(), this.PrefixString("\t")))
	}
	return append(this, item), nil
}

/*
If this does not already contain a copy of item, item is added to this.
*/
func (this ItemList) AddNoDuplicate(items ...*Item) ItemList {
	newList := this
	for _, item := range items {
		if !newList.Contain(item) {
			newList = append(newList, item)
		}
	}
	return newList
}

/*
See Algorithm: set.Closure() in package doc
*/
func (this ItemList) Closure(lexPart *ast.LexPart, symbols *symbols.Symbols) ItemList {
	closure := this
	for i := 0; i < len(closure); i++ {
		expSym := closure[i].ExpectedSymbol()
		if regDefId, isRegDefId := expSym.(*ast.LexRegDefId); isRegDefId {
			if !this.ContainShift(expSym.String()) && !symbols.IsImport(regDefId.Id) {
				closure = closure.AddNoDuplicate(NewItem(regDefId.Id, lexPart, symbols).Emoves()...)
			}
		}
	}
	return closure
}

func (this ItemList) Contain(that *Item) bool {
	for _, item := range this {
		if item.Equal(that) {
			return true
		}
	}
	return false
}

func (this ItemList) ContainShift(id string) bool {
	for _, item := range this {
		if item.Id == id && !item.Reduce() {
			return true
		}
	}
	return false
}

func (this ItemList) Equal(that ItemList) bool {
	if len(this) != len(that) {
		return false
	}
	for _, item := range this {
		if !that.Contain(item) {
			return false
		}
	}
	return true
}

func (this ItemList) indexOf(that *Item) int {
	for i, item := range this {
		if item.Equal(that) {
			return i
		}
	}
	return -1
}

func (this ItemList) Remove(item *Item) ItemList {
	idx := this.indexOf(item)
	if idx == -1 {
		panic(fmt.Sprintf("Cannot find item: %s", item.String()))
	}
	if idx < len(this)-1 {
		copy(this[idx:], this[idx+1:])
	}
	return this[:len(this)-1]
}

func (this ItemList) PrefixString(prefix string) string {
	buf := new(strings.Builder)
	for _, item := range this {
		fmt.Fprintf(buf, "%s%s\n", prefix, item.String())
	}
	return buf.String()
}
