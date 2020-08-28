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
	"io"
	"strings"

	"github.com/goccmack/gocc/internal/ast"
	"github.com/goccmack/gocc/internal/lexer/symbols"
	"github.com/goccmack/gocc/internal/util"
)

type Item struct {
	Id        string
	Prod      ast.LexProduction
	pos       *itemPos
	ProdIndex ast.LexProdIndex
	hashKey   string
	str       string
	*symbols.Symbols
}

/*
For a lex production,
  T : s
func NewItem returns
  T : •s
without executing the ℇ-moves for s.

Func *Item.Emoves must be called to return the set of basic items for T : •s.
*/
func NewItem(prodId string, lexPart *ast.LexPart, symbols *symbols.Symbols) *Item {
	prod := lexPart.Production(prodId)
	item := &Item{
		Id:        prodId,
		Prod:      prod,
		pos:       newItemPos(prod.LexPattern()),
		ProdIndex: lexPart.ProdIndex(prodId),
		Symbols:   symbols,
	}
	item.getHashKey()
	return item
}

func (this *Item) Clone() *Item {
	return &Item{
		Id:        this.Id,
		Prod:      this.Prod,
		pos:       this.pos.clone(),
		ProdIndex: this.ProdIndex,
		Symbols:   this.Symbols,
	}
}

/*
"this" may be a non-basic item. The function returns the set of basic items after ℇ-moves have be executed on this.

For a general description of dotted items (items) and ℇ-moves of items, see:

  Modern Compiler Design. Dick Grune, et al. Second Edition. Springer 2012.

ℇ-moves for lex items:

  Lex T be any production.
  Let w,x,y,z be any strings of lexical symbols.
  Let s = x|...|y have one or more alternatives.

  Then

  T : •s              =>  T : •x|...|z
                          ...
                          T :  x|...|•z

  T : x•[y]z          =>  T : x[•y]z
                          T : x[y]•z

  T : x[y•]z          =>  T : x[y]•z

  T : x•{y}z          =>  T : x{•y}z
                          T : x{y}•z

  T : x{y•}z          =>  T : x{•y}z
                          T : x{y}•z

  T : w•(x|...|y)z    =>  T : w(•x|...|y)z
                          ...
                          T : w( x|...|•y)z

  T : x(...|y•|...)z  =>  T : x(...|y|...)•z
*/
func (this *Item) Emoves() (items []*Item) {
	newItems := util.NewStack(8).Push(this)
	for newItems.Len() > 0 {
		item := newItems.Pop().(*Item)

		if item.Reduce() || item.nextIsTerminal() {
			items = append(items, item)
			continue
		}

		nt, pos := item.pos.top()

		switch node := nt.(type) {
		case *ast.LexPattern:
			newItems.Push(item.eMovesLexPattern(node, pos)...)
		case *ast.LexGroupPattern:
			newItems.Push(item.eMovesGroupPattern(node, pos)...)
		case *ast.LexOptPattern:
			newItems.Push(item.eMovesOptPattern(node, pos)...)
		case *ast.LexRepPattern:
			newItems.Push(item.eMovesRepPattern(node, pos)...)
		case *ast.LexAlt:
			if newItem := item.eMovesLexAlt(node, pos); newItem != item {
				newItems.Push(newItem)
			} else {
				items = append(items, item)
			}
		default:
			panic(fmt.Sprintf("Unexpected type in items.Emoves(): %T", nt))
		}
	}

	return
}

func (this *Item) eMovesLexPattern(nt *ast.LexPattern, pos int) []interface{} {
	if pos == 0 {
		return this.newLexPatternBasicItems(nt, pos)
	}

	postItem := this.Clone()
	if postItem.pos.level() == 0 {
		postItem.pos.stack[0].pos = nt.Len()
	} else {
		postItem.pos.pop()
		postItem.pos.inc()
	}
	postItem.getHashKey()

	return []interface{}{postItem}
}

func (this *Item) eMovesGroupPattern(nt *ast.LexGroupPattern, pos int) []interface{} {
	if pos == 0 {
		return this.newLexPatternBasicItems(nt.LexPattern, pos)
	}

	postItem := this.Clone()
	postItem.pos.pop()
	postItem.pos.inc()
	postItem.getHashKey()

	return []interface{}{postItem}
}

func (this *Item) eMovesLexAlt(nt *ast.LexAlt, pos int) *Item {
	if pos >= len(nt.Terms) {
		postItem := this.Clone()
		postItem.pos.pop()
		postItem.pos.setToEnd()
		postItem.getHashKey()
		return postItem
	}

	nextTerm := nt.Terms[pos]

	if nextTerm.LexTerminal() {
		return this
	}

	newItem := this.Clone()
	newItem.pos.push(nextTerm.(ast.LexNTNode), 0)
	newItem.getHashKey()

	return newItem
}

func (this *Item) eMovesOptPattern(nt *ast.LexOptPattern, pos int) (items []interface{}) {
	if pos == 0 {
		items = append(items, this.newLexPatternBasicItems(nt.LexPattern, pos)...)
	}

	postItem := this.Clone()
	postItem.pos.pop()
	postItem.pos.inc()
	postItem.getHashKey()
	items = append(items, postItem)

	return
}

func (this *Item) eMovesRepPattern(nt *ast.LexRepPattern, pos int) (items []interface{}) {
	items = append(items, this.newLexPatternBasicItems(nt.LexPattern, pos)...)

	postItem := this.Clone()
	postItem.pos.pop()
	postItem.pos.inc()
	postItem.getHashKey()
	items = append(items, postItem)

	return
}

func (this *Item) getHashKey() {
	w := new(strings.Builder)
	fmt.Fprintf(w, "%d:", this.ProdIndex)
	for i, stackElem := range this.pos.stack {
		if i < len(this.pos.stack)-1 {
			fmt.Fprintf(w, "%d:", stackElem.pos)
		} else {
			fmt.Fprintf(w, "%d", stackElem.pos)
		}
	}
	this.hashKey = w.String()
}

func (this *Item) HashKey() string {
	if this.hashKey == "" {
		panic("Uninitialised hash key")
	}
	return this.hashKey
}

func (this *Item) newLexPatternBasicItems(nt *ast.LexPattern, pos int) []interface{} {
	numAlts := len(nt.Alternatives)
	items := make([]interface{}, numAlts)
	for i, alt := range nt.Alternatives {
		altItem := this.Clone()
		altItem.pos.setPos(i)
		altItem.pos.push(alt, 0)
		altItem.getHashKey()
		items[i] = altItem
	}
	return items
}

func (this *Item) nextIsTerminal() bool {
	nt, pos := this.pos.top()
	if pos >= nt.Len() {
		return false
	}
	return nt.Element(pos).LexTerminal()
}

func (this *Item) Equal(that *Item) bool {
	if this.hashKey == "" || that.hashKey == "" {
		panic("nil hashkey")
	}
	return this.hashKey == that.hashKey
}

/*
This function returns the expected symbol for shift items and
nil for reduce items.
This is the position of a basic item -- no ℇ-moves are possible.
*/
func (this *Item) ExpectedSymbol() (node ast.LexTNode) {
	nt, pos := this.pos.top()
	if pos >= nt.Len() {
		return nil
	}
	return nt.Element(pos).(ast.LexTNode)
}

func (this *Item) match(rng CharRange) bool {
	lexTerm := this.ExpectedSymbol()
	if lexTerm == nil {
		return false
	}
	switch t := lexTerm.(type) {
	case *ast.LexDot:
		return false
	case *ast.LexCharLit:
		return rng.From == t.Val && rng.To == t.Val
	case *ast.LexCharRange:
		return rng.From >= t.From.Val && rng.From <= t.To.Val &&
			rng.To <= t.To.Val
	case *ast.LexRegDefId:
		return false
	}
	panic(fmt.Sprintf("Unexpected lexTerm type: %T", lexTerm))
}

func (this *Item) Move(rng CharRange) []*Item {
	if !this.match(rng) {
		return nil
	}
	movedItem := this.Clone()
	movedItem.pos.inc()
	movedItem.getHashKey()

	items := movedItem.Emoves()

	return items
}

func (this *Item) MoveDot() []*Item {
	if this.ExpectedSymbol() != ast.LexDOT {
		return nil
	}
	movedItem := this.Clone()
	movedItem.pos.inc()
	movedItem.getHashKey()
	items := movedItem.Emoves()

	return items
}

func (this *Item) MoveRegDefId(id string) []*Item {
	if rid, ok := this.ExpectedSymbol().(*ast.LexRegDefId); ok && id == rid.Id {
		movedItem := this.Clone()
		movedItem.pos.inc()
		movedItem.getHashKey()
		items := movedItem.Emoves()
		return items
	}
	return nil
}

/*
returns true for reduce items, like:
	T : xyz •

and false otherwise
*/
func (this *Item) Reduce() bool {
	if this.pos.level() != 0 {
		return false
	}
	ntNode, pos := this.pos.top()
	return pos >= ntNode.Len()
}

func (this *Item) getString() {
	w := new(strings.Builder)
	fmt.Fprintf(w, "%s : ", this.Prod.Id())
	if this.Reduce() {
		le, _ := this.pos.top()
		if le.Len() > 1 {
			fmt.Fprintf(w, "(")
		}
		WriteStringNode(w, this.Prod.LexPattern(), this.pos)
		if le.Len() > 1 {
			fmt.Fprintf(w, ") •")
		} else {
			fmt.Fprintf(w, " •")
		}
	} else {
		WriteStringNode(w, this.Prod.LexPattern(), this.pos)
	}
	this.str = w.String()
}

func (this *Item) String() string {
	if this.str == "" {
		this.getString()
	}
	return this.str
}

func WriteStringNode(w io.Writer, node ast.LexNode, pos *itemPos) {
	topNode, idx := pos.top()

	switch n := node.(type) {
	case *ast.LexPattern:
		writeStringPattern(w, n, pos)
	case *ast.LexDot:
		fmt.Fprintf(w, ".")
	case *ast.LexCharLit:
		fmt.Fprintf(w, "%s", n.String())
	case *ast.LexCharRange:
		fmt.Fprintf(w, " %s-%s", n.From.String(), n.To.String())
	case *ast.LexGroupPattern:
		fmt.Fprintf(w, "(")
		writeStringPattern(w, n.LexPattern, pos)
		fmt.Fprintf(w, ")")
	case *ast.LexAlt:
		for i, term := range n.Terms {
			if n == topNode && i == idx {
				fmt.Fprintf(w, "• ")
			}
			WriteStringNode(w, term, pos)
			if i < len(n.Terms)-1 {
				fmt.Fprintf(w, " ")
			}
		}
	case *ast.LexOptPattern:
		fmt.Fprintf(w, "[")
		writeStringPattern(w, n.LexPattern, pos)
		fmt.Fprintf(w, "]")
	case *ast.LexRegDefId:
		fmt.Fprintf(w, "%s", n.String())
	case *ast.LexRepPattern:
		fmt.Fprintf(w, "{")
		writeStringPattern(w, n.LexPattern, pos)
		fmt.Fprintf(w, "}")
	default:
		panic(errors.New(fmt.Sprintf("Unexpected type of node, %T", node)))
	}
}

func writeStringPattern(w io.Writer, pattern *ast.LexPattern, pos *itemPos) {
	for i, alt := range pattern.Alternatives {
		WriteStringNode(w, alt, pos)
		if i < len(pattern.Alternatives)-1 {
			fmt.Fprintf(w, " | ")
		}
	}
}
