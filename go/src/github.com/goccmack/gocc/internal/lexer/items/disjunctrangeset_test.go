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
	"testing"

	"github.com/goccmack/gocc/internal/ast"
	// "unicode"
)

/*
ABCDEF GHIJKL MNOPQR STUVWX YZ
*/
func TestDisjunctSets0(t *testing.T) {
	set := NewDisjunctRangeSet()
	data := []CharRange{
		{'A', 'F'},
		{'G', 'L'},
		{'M', 'R'},
		{'S', 'X'},
		{'Y', 'Z'},
	}

	for _, i := range []int{2, 1, 0, 4, 3} {
		set.AddRange(data[i].From, data[i].To)
	}

	if list := set.List(); len(list) != 5 {
		t.Fatalf("len(list) == %d", len(list))
	}
}

/*
ABCDEF GHIJKL MNOPQR STUVWX YZ
			|===========|
|-------|	|===========| 1
|-----------+---|-------+ 2
|-----------+-----------| 3
|-----------+-----------+---| 4
			|---|-------+ 5
			|===========| 6
			|-----------+---| 7
			+===========+	|----| 8
			+---|----|--+ 9
			+---|-------| 10
			+---|-------+---| 11
*/

/*
ABCDEF 		 MNOPQR
			|===========|
|-------|	|===========| 1
*/
func TestDisjunctSets1(t *testing.T) {
	set := NewDisjunctRangeSet()
	data := []CharRange{
		{'A', 'F'},
		{'M', 'R'},
	}
	set.addRange(data[1], data[0])
	rng := set.List()
	if len(rng) != 2 {
		t.Fatalf("len(rng) == %d", len(rng))
	}
	for i, rng := range set.List() {
		if !data[i].Equal(rng) {
			t.Fatalf("set[%d]==%s, expected %s", i, rng.String(), data[i].String())
		}
	}
}

/*
ABCDEF 		 MNOPQR
			|===========|
|-----------+---|-------+ 2
*/
func TestDisjunctSets2(t *testing.T) {
	set := NewDisjunctRangeSet()
	input := []CharRange{
		{'M', 'R'},
		{'A', 'O'},
	}
	check := []CharRange{
		{'A', 'L'},
		{'M', 'O'},
		{'P', 'R'},
	}
	set.addRange(input...)
	ranges := set.List()
	if len(ranges) != len(check) {
		t.Fatalf("len(rng) == %d", len(ranges))
	}
	for i, rng := range ranges {
		if !check[i].Equal(rng) {
			t.Fatalf("set[%d]==%s, expected %s", i, rng.String(), check[i].String())
		}
	}
}

/*
ABCDEF 		GHIJKL
			|===========|
|-----------+-----------| 3
*/
func TestDisjunctSets3(t *testing.T) {
	set := NewDisjunctRangeSet()
	input := []CharRange{
		{'G', 'L'},
		{'A', 'F'},
	}
	check := []CharRange{
		{'A', 'F'},
		{'G', 'L'},
	}
	set.addRange(input...)
	ranges := set.List()
	if len(ranges) != len(check) {
		t.Fatalf("len(rng) == %d", len(ranges))
	}
	for i, rng := range ranges {
		if !check[i].Equal(rng) {
			t.Fatalf("set[%d]==%s, expected %s", i, rng.String(), check[i].String())
		}
	}
}

/*
ABCDEF GHIJKL MNOPQR STUVWX YZ
			|===========|
|-----------+-----------+---| 4
*/
func TestDisjunctSets4(t *testing.T) {
	set := NewDisjunctRangeSet()
	input := []CharRange{
		{'G', 'L'},
		{'A', 'R'},
	}
	check := []CharRange{
		{'A', 'F'},
		{'G', 'L'},
		{'M', 'R'},
	}
	set.addRange(input...)
	ranges := set.List()
	if len(ranges) != len(check) {
		t.Fatalf("len(rng) == %d", len(ranges))
	}
	for i, rng := range ranges {
		if !check[i].Equal(rng) {
			t.Fatalf("set[%d]==%s, expected %s", i, rng.String(), check[i].String())
		}
	}
}

/*
ABCDEF GHIJKL MNOPQR STUVWX YZ
			|===========|
			|---|-------+ 5
*/
func TestDisjunctSets5(t *testing.T) {
	set := NewDisjunctRangeSet()
	input := []CharRange{
		{'G', 'L'},
		{'G', 'I'},
	}
	check := []CharRange{
		{'G', 'I'},
		{'J', 'L'},
	}
	set.addRange(input...)
	ranges := set.List()
	if len(ranges) != len(check) {
		t.Fatalf("len(rng) == %d", len(ranges))
	}
	for i, rng := range ranges {
		if !check[i].Equal(rng) {
			t.Fatalf("set[%d]==%s, expected %s", i, rng.String(), check[i].String())
		}
	}
}

/*
ABCDEF GHIJKL MNOPQR STUVWX YZ
			|===========|
			|===========| 6
*/
func TestDisjunctSets6(t *testing.T) {
	set := NewDisjunctRangeSet()
	input := []CharRange{
		{'G', 'L'},
		{'G', 'L'},
	}
	check := []CharRange{
		{'G', 'L'},
	}
	set.addRange(input...)
	ranges := set.List()
	if len(ranges) != len(check) {
		t.Fatalf("len(rng) == %d", len(ranges))
	}
	for i, rng := range ranges {
		if !check[i].Equal(rng) {
			t.Fatalf("set[%d]==%s, expected %s", i, rng.String(), check[i].String())
		}
	}
}

/*
ABCDEF GHIJKL MNOPQR STUVWX YZ
			|===========|
			|-----------+---| 7
*/
func TestDisjunctSets7(t *testing.T) {
	set := NewDisjunctRangeSet()
	input := []CharRange{
		{'G', 'L'},
		{'G', 'O'},
	}
	check := []CharRange{
		{'G', 'L'},
		{'M', 'O'},
	}
	set.addRange(input...)
	ranges := set.List()
	if len(ranges) != len(check) {
		t.Fatalf("len(rng) == %d", len(ranges))
	}
	for i, rng := range ranges {
		if !check[i].Equal(rng) {
			t.Fatalf("set[%d]==%s, expected %s", i, rng.String(), check[i].String())
		}
	}
}

/*
ABCDEF GHIJKL MNOPQR STUVWX YZ
			|===========|
			+===========+	|----| 8
*/
func TestDisjunctSets8(t *testing.T) {
	set := NewDisjunctRangeSet()
	input := []CharRange{
		{'G', 'L'},
		{'S', 'X'},
	}
	check := []CharRange{
		{'G', 'L'},
		{'S', 'X'},
	}
	set.addRange(input...)
	ranges := set.List()
	if len(ranges) != len(check) {
		t.Fatalf("len(rng) == %d", len(ranges))
	}
	for i, rng := range ranges {
		if !check[i].Equal(rng) {
			t.Fatalf("set[%d]==%s, expected %s", i, rng.String(), check[i].String())
		}
	}
}

/*
ABCDEF GHIJKL MNOPQR STUVWX YZ
			|===========|
			+---|----|--+ 9
*/
func TestDisjunctSets9(t *testing.T) {
	set := NewDisjunctRangeSet()
	input := []CharRange{
		{'G', 'R'},
		{'J', 'O'},
	}
	check := []CharRange{
		{'G', 'I'},
		{'J', 'O'},
		{'P', 'R'},
	}
	set.addRange(input...)
	ranges := set.List()
	if len(ranges) != len(check) {
		t.Fatalf("len(rng) == %d", len(ranges))
	}
	for i, rng := range ranges {
		if !check[i].Equal(rng) {
			t.Fatalf("set[%d]==%s, expected %s", i, rng.String(), check[i].String())
		}
	}
}

/*
ABCDEF GHIJKL MNOPQR STUVWX YZ
			|===========|
			+---|-------| 10
*/
func TestDisjunctSets10(t *testing.T) {
	set := NewDisjunctRangeSet()
	input := []CharRange{
		{'G', 'R'},
		{'M', 'R'},
	}
	check := []CharRange{
		{'G', 'L'},
		{'M', 'R'},
	}
	set.addRange(input...)
	ranges := set.List()
	if len(ranges) != len(check) {
		t.Fatalf("len(rng) == %d", len(ranges))
	}
	for i, rng := range ranges {
		if !check[i].Equal(rng) {
			t.Fatalf("set[%d]==%s, expected %s", i, rng.String(), check[i].String())
		}
	}
}

/*
ABCDEF GHIJKL MNOPQR STUVWX YZ
			|===========|
			+---|-------+---| 11
*/
func TestDisjunctSets11(t *testing.T) {
	set := NewDisjunctRangeSet()
	input := []CharRange{
		{'G', 'L'},
		{'J', 'R'},
	}
	check := []CharRange{
		{'G', 'I'},
		{'J', 'L'},
		{'M', 'R'},
	}
	set.addRange(input...)
	ranges := set.List()
	if len(ranges) != len(check) {
		t.Fatalf("len(rng) == %d", len(ranges))
	}
	for i, rng := range ranges {
		if !check[i].Equal(rng) {
			t.Fatalf("set[%d]==%s, expected %s", i, rng.String(), check[i].String())
		}
	}
}

/*
ABC DEF GHI JKL MNO PQR STUVWX YZ
	|=======|	|=======|
|---+-------+---+-----------+
*/
func TestDisjunctSets12(t *testing.T) {
	set := NewDisjunctRangeSet()
	input := []CharRange{
		{'D', 'F'},
		{'J', 'L'},
		{'A', 'O'},
	}
	check := []CharRange{
		{'A', 'C'},
		{'D', 'F'},
		{'G', 'I'},
		{'J', 'L'},
		{'M', 'O'},
	}
	set.addRange(input...)
	ranges := set.List()
	if len(ranges) != len(check) {
		t.Fatalf("len(rng) == %d", len(ranges))
	}
	for i, rng := range ranges {
		if !check[i].Equal(rng) {
			t.Fatalf("set[%d]==%s, expected %s", i, rng.String(), check[i].String())
		}
	}
}

func TestDisjunctSets13(t *testing.T) {
	set := NewDisjunctRangeSet()
	sym := &ast.LexCharLit{
		Val: 'A',
	}
	set.AddLexTNode(sym)
	check := CharRange{'A', 'A'}
	ranges := set.List()
	if len(ranges) != 1 {
		t.Fatalf("len(rng) == %d", len(ranges))
	}
	if !ranges[0].Equal(check) {
		t.Fatalf("range==%s, expected %s", ranges[0], check)
	}
}

func TestDisjunctSets14(t *testing.T) {
	set := NewDisjunctRangeSet()
	set.AddRange('D', 'F')
	set.AddRange('J', 'L')
	set.AddLexTNode(ast.LexDOT)
	if len(set.set) != 2 {
		t.Fatalf("len(set.set) == %d", len(set.set))
	}
}
