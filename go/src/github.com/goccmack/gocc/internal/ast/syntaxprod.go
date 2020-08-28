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

	"github.com/goccmack/gocc/internal/frontend/token"
)

type SyntaxProd struct {
	Id   string
	Body *SyntaxBody
}

func NewSyntaxProd(prodId, alts interface{}) ([]*SyntaxProd, error) {
	pid := string(prodId.(*token.Token).Lit)
	alts1 := alts.(SyntaxAlts)
	prods := make([]*SyntaxProd, len(alts1))
	for i, body := range alts1 {
		prods[i] = &SyntaxProd{
			Id:   pid,
			Body: body,
		}
	}
	return prods, nil
}

func (this *SyntaxProd) String() string {
	return fmt.Sprintf("%s : %s", this.Id, this.Body.String())
}
