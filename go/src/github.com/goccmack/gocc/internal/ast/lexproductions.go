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

type LexProductions struct {
	Productions []LexProduction
}

func NewLexProductions(lexProd interface{}) (*LexProductions, error) {
	return &LexProductions{
		Productions: []LexProduction{lexProd.(LexProduction)},
	}, nil
}

func newLexProductions() *LexProductions {
	return &LexProductions{
		Productions: []LexProduction{},
	}
}

func AppendLexProduction(lexProds, prod interface{}) (*LexProductions, error) {
	lp := lexProds.(*LexProductions)
	lp.Productions = append(lp.Productions, prod.(LexProduction))
	return lp, nil
}

func (this *LexProductions) String() string {
	w := new(strings.Builder)
	for _, prod := range this.Productions {
		fmt.Fprintf(w, "%s ;", prod.String())
	}
	return w.String()
}
