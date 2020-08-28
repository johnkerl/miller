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

type LexAlt struct {
	Terms []LexTerm
}

func NewLexAlt(lexTerm interface{}) (*LexAlt, error) {
	return &LexAlt{
		Terms: []LexTerm{lexTerm.(LexTerm)},
	}, nil
}

func AppendLexTerm(lexAlt, lexTerm interface{}) (*LexAlt, error) {
	la := lexAlt.(*LexAlt)
	la.Terms = append(la.Terms, lexTerm.(LexTerm))
	return la, nil
}

func (this *LexAlt) Contain(term LexTerm) bool {
	for _, thisTerm := range this.Terms {
		if thisTerm == term {
			return true
		}
	}
	return false
}

func (this *LexAlt) String() string {
	buf := new(strings.Builder)
	for i, term := range this.Terms {
		if i > 0 {
			fmt.Fprintf(buf, " ")
		}
		fmt.Fprintf(buf, "%s", term.String())
	}
	return buf.String()
}
