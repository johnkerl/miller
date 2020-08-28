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

type LexImport struct {
	Id      string
	ExtFunc string
}

func NewLexImport(regDefId, extFunc interface{}) (*LexImport, error) {
	return &LexImport{
		Id:      string(regDefId.(*token.Token).Lit),
		ExtFunc: getExtFunc(extFunc),
	}, nil
}

func (this *LexImport) IsTerminal() bool {
	return true
}

func (this *LexImport) String() string {
	return fmt.Sprintf("%s \"%s\"", this.Id, this.ExtFunc)
}

func getExtFunc(strLit interface{}) string {
	lit := strLit.(*token.Token).Lit
	return string(lit[1 : len(lit)-1])
}
