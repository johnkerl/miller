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
	"fmt"
	"strings"

	"github.com/goccmack/gocc/internal/frontend/token"
)

type LexIgnoredTokDef struct {
	id      string
	pattern *LexPattern
}

func NewLexIgnoredTokDef(tokId, lexPattern interface{}) (*LexIgnoredTokDef, error) {
	tokDef := &LexIgnoredTokDef{
		id:      string(tokId.(*token.Token).Lit),
		pattern: lexPattern.(*LexPattern),
	}
	return tokDef, nil
}

func (*LexIgnoredTokDef) RegDef() bool {
	return false
}

func (this *LexIgnoredTokDef) String() string {
	buf := new(strings.Builder)
	fmt.Fprintf(buf, "%s : %s", this.id, this.pattern.String())
	return buf.String()
}

func (this *LexIgnoredTokDef) Id() string {
	return this.id
}

func (this *LexIgnoredTokDef) LexPattern() *LexPattern {
	return this.pattern
}
