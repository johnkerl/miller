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
	"sort"
	"strings"
)

// All maps are indexed by production id
type LexPart struct {
	Header *FileHeader
	*LexImports
	TokDefsList        []*LexTokDef
	TokDefs            map[string]*LexTokDef
	stringLitToks      map[string]*LexTokDef
	RegDefsList        []*LexRegDef
	RegDefs            map[string]*LexRegDef
	IgnoredTokDefsList []*LexIgnoredTokDef
	IgnoredTokDefs     map[string]*LexIgnoredTokDef
	ProdList           *LexProductions
	ProdMap            *LexProdMap
}

func NewLexPart(header, imports, prodList interface{}) (*LexPart, error) {
	lexPart := &LexPart{
		TokDefsList:        make([]*LexTokDef, 0, 16),
		TokDefs:            make(map[string]*LexTokDef, 16),
		stringLitToks:      make(map[string]*LexTokDef, 16),
		RegDefsList:        make([]*LexRegDef, 0, 16),
		RegDefs:            make(map[string]*LexRegDef, 16),
		IgnoredTokDefsList: make([]*LexIgnoredTokDef, 0, 16),
		IgnoredTokDefs:     make(map[string]*LexIgnoredTokDef, 16),
	}
	if header != nil {
		lexPart.Header = header.(*FileHeader)
	} else {
		lexPart.Header = new(FileHeader)
	}
	if imports != nil {
		lexPart.LexImports = imports.(*LexImports)
	} else {
		lexPart.LexImports = newLexImports()
	}
	if prodList != nil {
		lexPart.ProdList = prodList.(*LexProductions)
		lexPart.ProdMap = NewLexProdMap(prodList.(*LexProductions))
		for _, p := range prodList.(*LexProductions).Productions {
			pid := p.Id()

			switch p1 := p.(type) {
			case *LexTokDef:
				//TODO: decide whether to handle in separate semantic check
				if _, exist := lexPart.TokDefs[pid]; exist {
					return nil, errors.New(fmt.Sprintf("Duplicate token def: %s", pid))
				}
				lexPart.TokDefs[pid] = p1
				lexPart.TokDefsList = append(lexPart.TokDefsList, p1)
			case *LexRegDef:
				//TODO: decide whether to handle in separate semantic check
				if _, exist := lexPart.RegDefs[pid]; exist {
					return nil, errors.New(fmt.Sprintf("Duplicate token def: %s", pid))
				}
				lexPart.RegDefs[pid] = p1
				lexPart.RegDefsList = append(lexPart.RegDefsList, p1)
			case *LexIgnoredTokDef:
				if _, exist := lexPart.IgnoredTokDefs[pid]; exist {
					return nil, errors.New(fmt.Sprintf("Duplicate ignored token def: %s", pid))
				}
				lexPart.IgnoredTokDefs[pid] = p1
				lexPart.IgnoredTokDefsList = append(lexPart.IgnoredTokDefsList, p1)
			}
		}
	} else {
		lexPart.ProdList = newLexProductions()
		lexPart.ProdMap = newLexProdMap()
	}
	return lexPart, nil
}

func (this *LexPart) StringLitTokDef(id string) *LexTokDef {
	tokDef, _ := this.stringLitToks[id]
	return tokDef
}

func (this *LexPart) Production(id string) LexProduction {
	idx, exist := this.ProdMap.idMap[id]
	if !exist {
		panic(fmt.Sprintf("Unknown production: %s", id))
	}
	return this.ProdList.Productions[idx]
}

func (this *LexPart) ProdIndex(id string) LexProdIndex {
	idx, exist := this.ProdMap.idMap[id]
	if !exist {
		panic(fmt.Sprintf("Unknown production %s", id))
	}
	return idx
}

func (this *LexPart) TokenIds() []string {
	tids := make([]string, 0, len(this.TokDefs))
	for tid := range this.TokDefs {
		tids = append(tids, tid)
	}
	sort.Strings(tids)
	return tids
}

func (this *LexPart) UpdateStringLitTokens(tokens []string) {
	for _, strLit := range tokens {
		tokDef := NewLexStringLitTokDef(strLit)
		this.ProdMap.Add(tokDef)
		this.TokDefsList = append(this.TokDefsList, tokDef)
		this.TokDefs[strLit] = tokDef
		this.stringLitToks[strLit] = tokDef
		this.ProdList.Productions = append(this.ProdList.Productions, tokDef)
	}
}

func (this *LexPart) String() string {
	buf := new(strings.Builder)
	fmt.Fprintf(buf, "LexPart.ProdMap:\n")
	if this.ProdList != nil {
		for i, p := range this.ProdList.Productions {
			fmt.Fprintf(buf, "\t%d: %s\n", i, p)
		}
	}
	fmt.Fprintf(buf, "LexPart.TokDefs:\n")
	for t := range this.TokDefs {
		fmt.Fprintf(buf, "\t%s\n", t)
	}
	return buf.String()
}
