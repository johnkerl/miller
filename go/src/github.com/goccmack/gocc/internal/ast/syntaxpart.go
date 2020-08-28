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

type SyntaxPart struct {
	Header   *FileHeader
	ProdList SyntaxProdList
}

func NewSyntaxPart(header, prodList interface{}) (*SyntaxPart, error) {
	sp := &SyntaxPart{}
	if header != nil {
		sp.Header = header.(*FileHeader)
	} else {
		sp.Header = new(FileHeader)
	}
	if prodList != nil {
		sp.ProdList = prodList.(SyntaxProdList)
	}

	return sp, nil
}

func (this *SyntaxPart) augment() *SyntaxPart {
	startProd := &SyntaxProd{
		Id: "S'",
		Body: &SyntaxBody{
			Symbols: []SyntaxSymbol{SyntaxProdId(this.ProdList[0].Id)},
		},
	}
	newProdList := SyntaxProdList{startProd}
	return &SyntaxPart{
		Header:   this.Header,
		ProdList: append(newProdList, this.ProdList...),
	}
}
