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

package first

import (
	"fmt"
	"sort"
	"strings"
)

/*
key: symbol string
*/
type SymbolSet map[string]bool

func (this SymbolSet) AddSet(that SymbolSet) {
	for id := range that {
		this[id] = true
	}
}

func (this SymbolSet) Equal(that SymbolSet) bool {
	if len(this) != len(that) {
		return false
	}
	for symbol := range this {
		if _, contain := that[symbol]; !contain {
			return false
		}
	}
	return true
}

func (this SymbolSet) String() string {
	buf := new(strings.Builder)
	fmt.Fprintf(buf, "{\n")
	var keys []string
	for key := range this {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, str := range keys {
		fmt.Fprintf(buf, "\t%s\n", str)
	}
	fmt.Fprintf(buf, "}")
	return buf.String()
}
