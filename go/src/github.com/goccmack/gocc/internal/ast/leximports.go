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
	"errors"
	"fmt"
	"strings"
)

type LexImports struct {
	Imports map[string]*LexImport
}

func NewLexImports(lexImport interface{}) (*LexImports, error) {
	imports, err := newLexImports().Add(lexImport.(*LexImport))
	return imports, err
}

func newLexImports() *LexImports {
	return &LexImports{
		Imports: make(map[string]*LexImport),
	}
}

func AddLexImport(imports, lexImport interface{}) (*LexImports, error) {
	return imports.(*LexImports).Add(lexImport.(*LexImport))
}

/*
Return true if a new lex import has been added.
Return false if lexImport is a duplicate.
*/
func (this *LexImports) Add(lexImport *LexImport) (*LexImports, error) {
	if _, exist := this.Imports[lexImport.Id]; exist {
		return nil, errors.New(fmt.Sprintf("Duplicate builtin declaration: %s", lexImport.String()))
	}
	this.Imports[lexImport.Id] = lexImport
	return this, nil
}

func (this *LexImports) String() string {
	w := new(strings.Builder)
	fmt.Fprintf(w, "import(\n")
	for _, imp := range this.Imports {
		fmt.Fprintf(w, "\t%s\n", imp.String())
	}
	fmt.Fprintf(w, ")")
	return w.String()
}
