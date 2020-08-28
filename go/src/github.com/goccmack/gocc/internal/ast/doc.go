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

/*
This package contains the Abstract Syntax Tree (AST) elements used by gocc to generate a target lexer and parser.

The top-level node is Grammar in grammar.go

The EBNF accepted by gocc consists of two parts:

1. The lexical part, containing the defintion of tokens.

2. The grammar or syntax part, containing the grammar or syntax of the language. Files containing grammar objects are prefixed with "g", e.g.: galts.go
*/
package ast
