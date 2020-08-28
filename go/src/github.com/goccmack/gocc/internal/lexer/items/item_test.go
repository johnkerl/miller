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
	"testing"

	"github.com/goccmack/gocc/internal/lexer/symbols"
)

/*
id : _a | _b _a =>
	id : • _a | _b _a
	id : _a | • _b _a

id : • a | _b =>
	id : (_a | _b _a) •

id : a | • _b _a =>
	id : _a | _b • _a =>
		id : (_a | _b _A) •
*/
func Test1(t *testing.T) {
	src := `_a : 'a'; _b : 'b' ; id : _a | _b _a ;`
	g := parse(src, t)
	symbols := symbols.NewSymbols(g.LexPart)

	items1 := NewItem("id", g.LexPart, symbols).Emoves()
	exp1 := "id : • _a | _b _a"
	exp2 := "id : _a | • _b _a"
	matchItems(t, items1, exp1, exp2)

	exp3 := "id : (_a | _b _a) •"
	matchItems(t, findItem(exp1, items1).MoveRegDefId("_a"), exp3)

	exp4 := "id : _a | _b • _a"
	matchItems(t, findItem(exp2, items1).MoveRegDefId("_b"), exp4)
}

/*

/*
id : _a {_b} =>
	id : • _a {_b} =>
		id : _a {• _b}
		id : _a {_b} •
*/
func Test2(t *testing.T) {
	src := `_a : 'a'; _b : 'b' ; id : _a {_b};`
	g := parse(src, t)
	symbols := symbols.NewSymbols(g.LexPart)

	items1 := NewItem("id", g.LexPart, symbols).Emoves()

	exp1 := "id : • _a {_b}"
	matchItems(t, items1, exp1)

	items2 := findItem(exp1, items1).MoveRegDefId("_a")
	exp2 := "id : _a {• _b}"
	exp3 := "id : _a {_b} •"
	matchItems(t, items2, exp2, exp3)
}

/*
id : [.] . =>
	id : [• .] .
	id : [.] • .

id : [.] • . =>
	id : [.] . •
*/
func Test3(t *testing.T) {
	src := `id : [.] .;`
	g := parse(src, t)
	symbols := symbols.NewSymbols(g.LexPart)

	items := NewItem("id", g.LexPart, symbols).Emoves()

	exp1 := "id : [• .] ."
	exp2 := "id : [.] • ."
	matchItems(t, items, exp1, exp2)
	exp3 := "id : [.] . •"
	matchItems(t, findItem(exp2, items).MoveDot(), exp3)
}

/*
id : _a (_a | _b) _c =>
	id : •a (_a | _b) _c =>
		id : _a (• _a | _b) _c
		id : _a (_a | • _b) _c

id : _a (• _a | _b) _c =>
	id : _a (_a | _b) •_c =>
		id : _a (_a | _b) _c •

id : _a (_a | • _b) _c
	id : _a (_a | _b) •_c =>
		id : _a (_a | _b) _c •
*/
func Test4(t *testing.T) {
	src := "id : _a (_a | _b) _c ;"
	g := parse(src, t)
	symbols := symbols.NewSymbols(g.LexPart)

	items := NewItem("id", g.LexPart, symbols).Emoves()

	exp1 := "id : • _a (_a | _b) _c"
	matchItems(t, items, exp1)

	items1 := findItem(exp1, items).MoveRegDefId("_a")
	exp2 := "id : _a (• _a | _b) _c"
	exp3 := "id : _a (_a | • _b) _c"
	matchItems(t, items1, exp2, exp3)

	items2 := findItem(exp2, items1).MoveRegDefId("_a")
	exp4 := "id : _a (_a | _b) • _c"
	matchItems(t, items2, exp4)
	items3 := findItem(exp4, items2).MoveRegDefId("_c")
	exp5 := "id : _a (_a | _b) _c •"
	matchItems(t, items3, exp5)

	items4 := findItem(exp3, items1).MoveRegDefId("_b")
	exp6 := "id : _a (_a | _b) • _c"
	matchItems(t, items4, exp6)
}

/*
id : {_a | _b} =>
	id : {• _a | _b}
	id : {_a | • _b}
	id : {_a | _b} •

id : {• _a | _b} =>
	id : {• _a | _b}
	id : {_a | • _b}
	id : {_a | _b} •
*/
func Test5(t *testing.T) {
	src := `id : {_a | _b};`
	g := parse(src, t)
	symbols := symbols.NewSymbols(g.LexPart)

	items := NewItem("id", g.LexPart, symbols).Emoves()

	exp1, exp2, exp3 := "id : {• _a | _b}", "id : {_a | • _b}", "id : {_a | _b} •"
	matchItems(t, items, exp1, exp2, exp3)

	items1 := findItem(exp1, items).MoveRegDefId("_a")
	matchItems(t, items1, exp1, exp2, exp3)

	items2 := findItem(exp2, items).MoveRegDefId("_b")
	matchItems(t, items2, exp1, exp2, exp3)
}

func findItem(item string, items []*Item) *Item {
	for _, item1 := range items {
		if item1.String() == item {
			return item1
		}
	}
	return nil
}

func findItemWithPosLevel(items []*Item, level int) *Item {
	for _, item := range items {
		if item.pos.level() == level {
			return item
		}
	}
	return nil
}

func matchItems(t *testing.T, generatedItems []*Item, expectedItems ...string) {
	dumpGeneratedItems := false
	if len(generatedItems) != len(expectedItems) {
		dumpGeneratedItems = true
		t.Errorf("Expected %d items, got %d\n", len(expectedItems), len(generatedItems))
	}
	itemMap := make(map[string]bool)
	for _, item := range generatedItems {
		itemMap[item.String()] = true
	}

	for _, item := range expectedItems {
		if _, ok := itemMap[item]; !ok {
			dumpGeneratedItems = true
			t.Errorf("Expected item >%s< not found\n", item)
		}
	}

	if dumpGeneratedItems {
		fmt.Printf("item_test.matchItems: Generated items:\n")
		for i, item := range generatedItems {
			fmt.Printf("\t%d: %s\n", i, item)
			// fmt.Printf("%s", item.pos)
		}
	}
}

// func parse(src string, t *testing.T) *ast.Grammar {
// 	scanner := new(scanner.Scanner)
// 	scanner.Init([]byte(src), token.FRONTENDTokens)
// 	parser := parser.NewParser(parser.ActionTable, parser.GotoTable, parser.ProductionsTable, token.FRONTENDTokens)
// 	g, err := parser.Parse(scanner)
// 	if err != nil {
// 		t.Fatal(fmt.Sprintf("Parse error: %s\n", err))
// 		return nil
// 	}

// 	return g.(*ast.Grammar)
// }
