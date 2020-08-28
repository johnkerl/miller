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

type LexNode interface {
	LexTerminal() bool
	String() string
}

func (LexAlt) LexTerminal() bool {
	return false
}

func (LexGroupPattern) LexTerminal() bool {
	return false
}

func (*LexIgnoredTokDef) LexTerminal() bool {
	return false
}

func (LexImports) LexTerminal() bool {
	return false
}

func (LexOptPattern) LexTerminal() bool {
	return false
}

func (LexPattern) LexTerminal() bool {
	return false
}

func (LexProductions) LexTerminal() bool {
	return false
}

func (*LexRegDef) LexTerminal() bool {
	return false
}

func (LexRepPattern) LexTerminal() bool {
	return false
}

func (*LexTokDef) LexTerminal() bool {
	return false
}

func (*LexCharLit) LexTerminal() bool {
	return true
}

func (*LexCharRange) LexTerminal() bool {
	return true
}

func (LexDot) LexTerminal() bool {
	return true
}

func (LexRegDefId) LexTerminal() bool {
	return true
}
