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
)

type itemPos struct {
	stack []stackElement
}

type stackElement struct {
	node ast.LexNTNode
	pos  int
}

func newItemPos(lexPattern *ast.LexPattern) (pos *itemPos) {
	pos = &itemPos{
		stack: make([]stackElement, 0, 8),
	}
	pos.push(lexPattern, 0)
	return
}

func (this *itemPos) clone() *itemPos {
	clone := &itemPos{
		stack: make([]stackElement, len(this.stack)),
	}
	copy(clone.stack, this.stack[0:len(this.stack)])
	return clone
}

func (this *itemPos) inc() {
	this.stack[this.level()].pos++
}

func (this *itemPos) pop() (ntNode ast.LexNTNode, pos int) {
	ntNode, pos = this.top()
	this.stack = this.stack[:len(this.stack)-1]
	return
}

func (this *itemPos) push(node ast.LexNTNode, pos int) {
	this.stack = append(this.stack, stackElement{node, pos})
}

func (this *itemPos) top() (nt ast.LexNTNode, pos int) {
	top := len(this.stack) - 1
	nt, pos = this.stack[top].node, this.stack[top].pos
	return
}

// returns level of stacking in the item. Bottom is 0
func (this *itemPos) level() int {
	return len(this.stack) - 1
}

/*
This function returns the ast.Node at the top of the stack
*/
func (this *itemPos) ntNode() ast.LexNTNode {
	n, _ := this.top()
	return n
}

// returns the position within the top level of the stack
func (this *itemPos) pos() int {
	_, pos := this.top()
	return pos
}

func (this *itemPos) setPos(i int) {
	this.stack[this.level()].pos = i
}

func (this *itemPos) setToEnd() {
	node, _ := this.top()
	this.stack[this.level()].pos = node.Len()
}

func (this *itemPos) equal(that *itemPos) bool {
	if len(this.stack) != len(that.stack) {
		return false
	}

	for i := 0; i < len(this.stack); i++ {
		if this.stack[i].node != that.stack[i].node ||
			this.stack[i].pos != that.stack[i].pos {

			return false
		}
	}
	return true
}

func (this *itemPos) String() string {
	buf := new(strings.Builder)
	for i := 0; i < len(this.stack); i++ {
		fmt.Fprintf(buf, "\t%T:%v; pos %d\n", (*this).stack[i].node, (*this).stack[i].node, this.stack[i].pos)
	}
	return buf.String()
}
