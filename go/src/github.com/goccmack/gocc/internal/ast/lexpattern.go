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
)

type LexPattern struct {
	Alternatives []*LexAlt
}

func NewLexPattern(lexAlt interface{}) (*LexPattern, error) {
	return &LexPattern{
		Alternatives: []*LexAlt{lexAlt.(*LexAlt)},
	}, nil
}

func AppendLexAlt(lexPattern, lexAlt interface{}) (*LexPattern, error) {
	lp := lexPattern.(*LexPattern)
	lp.Alternatives = append(lp.Alternatives, lexAlt.(*LexAlt))
	return lp, nil
}

func (this *LexPattern) String() string {
	buf := new(strings.Builder)
	for i, alt := range this.Alternatives {
		if i > 0 {
			fmt.Fprintf(buf, " | ")
		}
		fmt.Fprintf(buf, "%s", alt.String())
	}
	return buf.String()
}
