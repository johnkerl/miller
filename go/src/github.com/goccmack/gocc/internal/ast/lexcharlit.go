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

import (
	"github.com/goccmack/gocc/internal/frontend/token"
	"github.com/goccmack/gocc/internal/util"
)

type LexCharLit struct {
	Val rune
	Lit []byte
	s   string
}

func NewLexCharLit(tok interface{}) (*LexCharLit, error) {
	return newLexCharLit(tok), nil
}

func newLexCharLit(tok interface{}) *LexCharLit {
	c := new(LexCharLit)
	t := tok.(*token.Token)

	c.Val = util.LitToRune(t.Lit)
	c.Lit = t.Lit
	c.s = util.RuneToString(c.Val)

	return c
}

func newLexCharLitFromRune(c rune) *LexCharLit {
	cl := &LexCharLit{
		Val: c,
		s:   util.RuneToString(c),
	}
	cl.Lit = []byte(cl.s)
	return cl
}

func (this *LexCharLit) IsTerminal() bool {
	return true
}
func (this *LexCharLit) String() string {
	return this.s
}
