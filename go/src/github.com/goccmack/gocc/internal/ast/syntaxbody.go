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

type SyntaxBody struct {
	Error   bool
	Symbols SyntaxSymbols
	SDT     string
}

func NewSyntaxBody(symbols, sdtLit interface{}) (*SyntaxBody, error) {
	syntaxBody := &SyntaxBody{
		Error: false,
	}
	if symbols != nil {
		syntaxBody.Symbols = symbols.(SyntaxSymbols)
	}
	if sdtLit != nil {
		syntaxBody.SDT = sdtLit.(*token.Token).SDTVal()
	}
	return syntaxBody, nil
}

func NewErrorBody(symbols, sdtLit interface{}) (*SyntaxBody, error) {
	body, _ := NewSyntaxBody(symbols, sdtLit)
	body.Error = true
	return body, nil
}

func NewEmptyBody() (*SyntaxBody, error) {
	return NewSyntaxBody(nil, nil)
}

func (this *SyntaxBody) Empty() bool {
	return len(this.Symbols) == 0
}

func (this *SyntaxBody) String() string {
	return fmt.Sprintf("%s\t<< %s >>", this.Symbols.String(), this.SDT)
}
