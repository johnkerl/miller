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
	// "fmt"
	// "github.com/goccmack/gocc/internal/ast"
	"testing"
)

func TestItemPos1(t *testing.T) {
	src := `id : . {.};`

	g := parse(src, t)
	prod := g.LexPart.Production("id")
	itempos := newItemPos(prod.LexPattern())
	if itempos.level() != 0 {
		t.Fatalf("itempos.level() == %d", itempos.level())
	}
	if itempos.pos() != 0 {
		t.Fatalf("itempos.pos() == %d", itempos.pos())
	}
	if itempos.ntNode() != prod.LexPattern() {
		t.Fatalf("itempos.ntNode() == %T", itempos.ntNode())
	}
}

func TestItemPosClone1(t *testing.T) {
	src := `id : . {.};`

	g := parse(src, t)
	prod := g.LexPart.Production("id")
	itempos := newItemPos(prod.LexPattern())
	clone := itempos.clone()
	if clone.level() != itempos.level() {
		t.Fatalf("clone.level() == %d, itempos.level() == %d", clone.level(), itempos.level())
	}
}

func TestItemPos2(t *testing.T) {
	src := `id : . {.};`

	g := parse(src, t)
	prod := g.LexPart.Production("id")
	itempos := newItemPos(prod.LexPattern())
	node := prod.LexPattern().Alternatives[0]
	itempos.push(node, 0)
	if itempos.level() != 1 {
		t.Fatalf("itempos.level() == %d", itempos.level())
	}
	if itempos.pos() != 0 {
		t.Fatalf("itempos.pos() == %d", itempos.pos())
	}
	if itempos.ntNode() != node {
		t.Fatalf("itempos.ntNode() == %T", itempos.ntNode())
	}
	if itempos.level() != 1 {
		t.Fatalf("itempos.level() == %d", itempos.level())
	}
}
