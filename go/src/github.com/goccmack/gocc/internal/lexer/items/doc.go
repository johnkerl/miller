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

/*
Package items implements dotted items for FSA generation during the lexer generation process.

GENERALISED SUBSET ALGORITHM

The problem:
Lexers for some applications, e.g.: email header parsing have a number of char range input symbols, which overlap to a large extent.

The lexer symbol space can be very large, e.g.: UTF-8, so it is impractical to create items for every possible value of the char ranges.

Possible solutions:

1. Automatic conflict handling
Impractical because of overlapping char ranges for each state.

2. Separate symbol classification functions for each state
Better, but overlaps can still occur

3. Create a set of intersections of the valid char ranges of each set. Use the intersections as input symbols to the lexer.
This allows the full state space to be generated.

DEFINITION: LexPart = (Tokens, RegDefs, TermSym)

ALGORITHM: Generate lexer sets
	S0 := CreateSet0()
	itemSets = {S0}
	transitions = {}
	repeat {
		for set in itemSets {
			for symRange in set.SymRanges() {
				if nextSet := set.Next(symRange); nextSet not in itemSets {
					add nextSet to itemSets
					add (set, symRange, nextSet) to transitions
				}
			}
		}
	} until no more sets added to itemSets
	return itemSets, transitions


ALGORITHM: CreateSet0()
	S0 = {}
	for all T : x in Tokens {
		add T : •x to S0
	}
	return S0.Closure()


DEFINITION: Let R be the id of a regular definition.

ALGORITHM: set.Next()
	Output: the set containing all items in set that can move over symRange.

	repeat
		nextSet = {}
		for I : x •c y in set and c is in symRange{
			add I : x c •y to nextSet
		}
		for R : •z in nextSet and I : x •R y in set {
			add I : x •R y to nextSet
		}
		for R : z• in nextSet and I : x •R y in set {
			add I : x R• y to nextSet
		}
	until no items were added to nextSet

	return nextSet.Closure()


ALGORITHM: set.Closure()

	closure = {X : x •y | X : x •y in set and y not null}
	repeat
		for item I : x •R y in set and R : •z in RegDefs {
			add R : •z to closure
		}
	until no more items were added to closure


Disjunct character ranges

DEFINITIONS
 	The alphabet, S = [s[1],s[n]], of a regular grammar is an ordered set of n symbols. Given any s[i] != s[n] of S the next symbol in S
 	can be found  	by s[i] + 1. s[n] + 1 is undefined. For any s[i] != s[1] the previous symbol can be found by s[i] - 1. |S| is the
 	size of S.

 	Closed intervals of the alphabet, S, may be declared as symbol ranges, r = [r[1],r[k]] where the size of r, |r| = k. r is an ordered
 	subset of S with n elements without gaps, i.e.: r[i] + 1 = r[i+1].

 	X : x •a y is an item  of a regular grammar, with a the expected symbol and x, y possibly empty strings of symbols from S.

	In the algorithm symbols may be elements or ranges of S. If s is an element of S, |s| = 1. If |s| = 1, s = [s,s]



ALGORITHM: set.SymRanges() - Create the lexer symbol ranges of an item set, S:

 	Input: An item set of the grammar.

 	Output: A set of non-overlapping symbol ranges.

	Let SR = {} be a initially empty set of symbol ranges.

	for all items i in S {
		let s be the expected symbol of i
		if {
		// s outside all r in SR
		case for all r in I, s[|s|] < r[1] or s[1] > r[|r|]:
			add [s[1],s[|s|] to I

		// s equal to some r in SR
		case r in SR and s[1] = r[1] and s[|s|] = r[|r|]:
			do nothing

		//s inside some r in SR
		case r in SR and r[1] = s[1] and s[|s|] < r[|r|]:
			replace r in SR by {[r[1],s[|s|]], [s[|s|]+1,r[|r|]]}
		case r in SR and r[1] < s[1] and s[|s|] < r[|r|]:
			replace r in SR by {[r[1],s[1]-1], [s[s[1],s[|s|]], s[|s]+1,r[|r|]]}
		case r in SR and r[1] < s[1] and s[|s|] = r[|r|]:
			repace r in SR by {[r[1],s[1]-1], [s[1],r[|r|]]}

		// overlap to left
		case r in SR and s[1] < r[1] <= s[|s|] < r[|r|]:
			replace r in SR by {[s[1],r[1]-1], [r[1],s[|s|]], [s[|s|]+1,r[|r|]]}
		case r in SR and s[1] < r[1] and s[|s|] = r[|r|]:
			add [s[1],r[1]-1] to SR
		}

		// overlap to right
		case r in SR and s[1] = r[1] and s[|s|] > r[|r|]:
			add [r[|r|]+1,s[|s|]] to SR
		case r in SR and [s1] > r[1] and s[|s|] > r[|r|]:
			replace r in SR by {[r[1],s[1]-1], [s[1],r[|r|]], [r[|r|]+1,s[|s|]]}
	}


*/
package items
