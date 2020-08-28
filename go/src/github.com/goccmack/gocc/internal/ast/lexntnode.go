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

package ast

type LexNTNode interface {
	LexNode
	Element(int) LexNode
	Len() int
	Walk(LexNodeVisitor) LexNodeVisitor
}

// Element

func (this *LexAlt) Element(i int) LexNode {
	return this.Terms[i]
}

func (this *LexGroupPattern) Element(i int) LexNode {
	return this.LexPattern.Alternatives[i]
}

func (this *LexOptPattern) Element(i int) LexNode {
	return this.LexPattern.Alternatives[i]
}

func (this *LexPattern) Element(i int) LexNode {
	return this.Alternatives[i]
}

func (this *LexRepPattern) Element(i int) LexNode {
	return this.LexPattern.Alternatives[i]
}

// Len

func (this *LexAlt) Len() int {
	return len(this.Terms)
}

func (this *LexGroupPattern) Len() int {
	return len(this.LexPattern.Alternatives)
}

func (this *LexOptPattern) Len() int {
	return len(this.LexPattern.Alternatives)
}

func (this *LexPattern) Len() int {
	return len(this.Alternatives)
}

func (this *LexRepPattern) Len() int {
	return len(this.LexPattern.Alternatives)
}
