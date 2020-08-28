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
	"bytes"
	"fmt"
	"strings"

	"github.com/goccmack/gocc/internal/frontend/token"
)

type LexTokDef struct {
	id      string
	pattern *LexPattern
}

func NewLexTokDef(tokId, lexPattern interface{}) (*LexTokDef, error) {
	tokDef := &LexTokDef{
		id:      string(tokId.(*token.Token).Lit),
		pattern: lexPattern.(*LexPattern),
	}
	return tokDef, nil
}

func NewLexStringLitTokDef(tokId string) *LexTokDef {
	runes := bytes.Runes([]byte(tokId))
	alt, _ := NewLexAlt(newLexCharLitFromRune(runes[0]))
	for i := 1; i < len(runes); i++ {
		alt, _ = AppendLexTerm(alt, newLexCharLitFromRune(runes[i]))
	}
	ptrn, _ := NewLexPattern(alt)
	return &LexTokDef{
		id:      tokId,
		pattern: ptrn,
	}
}

func (*LexTokDef) RegDef() bool {
	return false
}

func (this *LexTokDef) String() string {
	buf := new(strings.Builder)
	fmt.Fprintf(buf, "%s : %s", this.id, this.pattern.String())
	return buf.String()
}

func (this *LexTokDef) Id() string {
	return this.id
}

func (this *LexTokDef) LexPattern() *LexPattern {
	return this.pattern
}
