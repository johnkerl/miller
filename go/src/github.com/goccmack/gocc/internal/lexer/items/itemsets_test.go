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
)

var testItemsets1src = `
a : '0'- '6' '0';
b : '4'-'9' '9';
c : '5' 'a' ;
`

var testItemSets1Data = []*itemSetTestData{
	{
		CharRange: CharRange{'0', '3'},
		itemSet:   []string{"a :  '0'-'6' • '0'"},
	},
	{
		CharRange: CharRange{'4', '4'},
		itemSet:   []string{"a :  '0'-'6' • '0'", "b :  '4'-'9' • '9'"},
	},
	{
		CharRange: CharRange{'5', '5'},
		itemSet:   []string{"a :  '0'-'6' • '0'", "b :  '4'-'9' • '9'", "c : '5' • 'a'"},
	},
	{
		CharRange: CharRange{'6', '6'},
		itemSet:   []string{"a :  '0'-'6' • '0'", "b :  '4'-'9' • '9'"},
	},
	{
		CharRange: CharRange{'7', '9'},
		itemSet:   []string{"b :  '4'-'9' • '9'"},
	},
}

// func TestItemSets1(t *testing.T) {
// 	g := parse(testItemsets1src, t)
// 	symbols := symbols.NewSymbols(g.LexPart)
// 	s0 := NewItemSet(0, g.LexPart, ItemsSet0(g.LexPart, symbols))
// 	for _, data := range testItemSets1Data {
// 		s1 := NewItemSet(1, g.LexPart, s0.Next(data.CharRange))
// 		checkItemsetResult(t, s1, data)
// 	}
// }

var testItemsets2src = `
<< import "unicode" >>

import(
	_upcase "unicode.IsUpper"
	_lowcase "unicode.IsLower"
	_digit "unicode.IsDigit"
)
_id_char : _upcase | _lowcase | '_' | _digit ;

_tokId : _lowcase {_id_char} ;

ignoredTokId : '!' _tokId ;
`

func _TestItemSets2(t *testing.T) {
	g := parse(testItemsets2src, t)
	itemSets := GetItemSets(g.LexPart)
	fmt.Printf("%s\n", itemSets)
}

var testItemsets3src = `
_lineComment : '/' '/' {.} '\n' ;

_blockComment : '/' '*' {.} '*' '/' ;

comment : _lineComment | _blockComment ;
`

func _TestItemSets3(t *testing.T) {
	g := parse(testItemsets3src, t)
	itemSets := GetItemSets(g.LexPart)
	fmt.Printf("%s\n", itemSets)
}

type itemSetTestData struct {
	CharRange
	itemSet
}

type itemSet []string

func (this itemSet) Contain(item string) bool {
	for _, it := range this {
		if it == item {
			return true
		}
	}
	return false
}

func checkItemsetResult(t *testing.T, set *ItemSet, data *itemSetTestData) {
	if len(set.Items) != len(data.itemSet) {
		testItemSetsErr(t, set, data, fmt.Sprintf("len(set.Items) == %v", len(set.Items)))
	}
}

func testItemSetsErr(t *testing.T, set *ItemSet, data *itemSetTestData, msg string) {
	t.Fatalf(msg)
}
