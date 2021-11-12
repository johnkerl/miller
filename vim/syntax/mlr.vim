" " Copyright 2021 John Kerl. All rights reserved.
" " Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.
" "
" " mlr.vim: Vim syntax file for the Miller DSL.
" 
if exists("b:current_syntax")
  finish
endif

syn case match

" ----------------------------------------------------------------
" Goal: map the lexical elements of the Miller DSL grammar in mlr.bnf
" to Vim syntax options
"   http://vimdoc.sourceforge.net/htmldoc/syntax.html
" ----------------------------------------------------------------

" ----------------------------------------------------------------
syn region mlrComment          start="#" end="$"
syn region mlrString           start=+"+ skip=+\\\\\|\\"+ end=+"+
syn keyword mlrBoolean         true false
syn match   mlrDecimalInt      "\(\<[0-9][0-9]*\>\)"
syn match   mlrHexadecimalInt  "\(\<0x[0-9a-fA-F][0-9a-fA-F]*\>\)"
syn match   mlrBinaryInt       "\(\<0b[01][01]*\>\)"

syn match   mlrFloat           "\.[0-9][0-9]*"
syn match   mlrFloat           "[0-9][0-9]*\."
syn match   mlrFloat           "[0-9][0-9]*\.[0-9][0-9]*"
syn match   mlrFloat           "[0-9][0-9]*\.[0-9]*"
syn match   mlrFloat           "\([0-9][0-9]*[eE][0-9][0-9]*\)"
syn match   mlrFloat           "\([0-9][0-9]*[eE]-[0-9][0-9]*\)"
syn match   mlrFloat           "\([0-9][0-9]*\.[0-9]*[eE][0-9][0-9]*\)"
syn match   mlrFloat           "\([0-9][0-9]*\.[0-9]*[eE]-[0-9][0-9]*\)"
syn match   mlrFloat           "\([0-9]*\.[0-9][0-9]*[eE][0-9][0-9]*\)"
syn match   mlrFloat           "\([0-9]*\.[0-9][0-9]*[eE]-[0-9][0-9]*\)"

syn keyword mlrConstant        M_PI M_E
syn keyword mlrContextVariable IPS IFS IRS OPS OFS ORS OFLATSEP NF NR FNR FILENAME FILENUM
syn keyword mlrENV             ENV

syn keyword mlrType            arr bool float int map num str var
syn keyword mlrKeyword         begin do elif else end filter for if in while break continue return
syn keyword mlrKeyword         func subr call dump edump
syn keyword mlrKeyword         emit emitp emitf eprint eprintn print printn tee stdout stderr unset null
syn match   mlrFieldName       /\(\$\<[a-zA-Z_][a-zA-Z_0-9]*\>\)/
syn match   mlrFieldName       /\(\$\*\)/
syn match   mlrOosvarName      /\(@\<[a-zA-Z_][a-zA-Z_0-9]*\>\)/
syn match   mlrOosvarName      /\(@\*\)/
syn match   mlrIdentifier      /\(\<[a-zA-Z_][a-zA-Z_0-9]*\>\)/ 
syn match   mlrFunctionCall    /\(\<[a-zA-Z_][a-zA-Z_0-9]*\>\)[ \t]*(/ 

syn region  mlrParen           start='(' end=')' transparent
syn region  mlrBlock           start="{" end="}" transparent

" Trailing whitespace; space-tab
syn match   mlrSpaceError      display excludenl "\s\+$"
syn match   mlrSpaceError      display " \+\t"me=e-1

syn match   mlrOperator         /!/
syn match   mlrOperator         /%/
syn match   mlrOperator         /&/
syn match   mlrOperator         /\*/
syn match   mlrOperator         /+/
syn match   mlrOperator         /-/
syn match   mlrOperator         /\./
syn match   mlrOperator         /\//
syn match   mlrOperator         /:/
syn match   mlrOperator         /</
syn match   mlrOperator         />/
syn match   mlrOperator         /?/
syn match   mlrOperator         /\^/
syn match   mlrOperator         /|/
syn match   mlrOperator         /\~/
syn match   mlrOperator         /!=/
syn match   mlrOperator         /&&/
syn match   mlrOperator         /\*\*/
syn match   mlrOperator         /\.\*/

syn match   mlrOperator         /\.+/
syn match   mlrOperator         /\.+/
syn match   mlrOperator         /\.-/
syn match   mlrOperator         /\.-/
syn match   mlrOperator         /\.\//
syn match   mlrOperator         /\/\//
syn match   mlrOperator         /<</
syn match   mlrOperator         /<=/
syn match   mlrOperator         /==/
syn match   mlrOperator         /=\~/
syn match   mlrOperator         />=/
syn match   mlrOperator         />>/
syn match   mlrOperator         /??/
syn match   mlrOperator         /^^/
syn match   mlrOperator         /||/
syn match   mlrOperator         /!=\~/
syn match   mlrOperator         /\.\/\//
syn match   mlrOperator         />>>/
syn match   mlrOperator         /???/
syn match   mlrOperator         /%=/
syn match   mlrOperator         /&=/
syn match   mlrOperator         /\*=/
syn match   mlrOperator         /+=/
syn match   mlrOperator         /-=/
syn match   mlrOperator         /\.=/
syn match   mlrOperator         /\/=/
syn match   mlrOperator         /^=/
syn match   mlrOperator         /|=/
syn match   mlrOperator         /&&=/
syn match   mlrOperator         /\*\*=/
syn match   mlrOperator         /\/\/=/
syn match   mlrOperator         /<<=/
syn match   mlrOperator         />>=/
syn match   mlrOperator         /??=/
syn match   mlrOperator         /^^=/
syn match   mlrOperator         /||=/
syn match   mlrOperator         />>>=/
syn match   mlrOperator         /???=/

" ----------------------------------------------------------------
hi def link mlrConstant        Constant
hi def link mlrString          String
hi def link mlrBoolean         Boolean
hi def link mlrNumber          Number
hi def link mlrFloat           Float
hi def link mlrInteger         Number

hi def link mlrComment         Comment
hi def link mlrKeyword         Keyword
hi def link mlrContextVariable Keyword
hi def link mlrENV             Keyword
hi def link mlrError           Error
hi def link mlrSpecial         Special
hi def link mlrType            Type

hi def link mlrDecimalInt      Number
hi def link mlrHexadecimalInt  Number
hi def link mlrBinaryInt       Number
hi def link mlrFieldName       Special
hi def link mlrOosvarName      Special
hi def link mlrIdentifier      Identifier
hi def link mlrFunctionCall    Type
hi def link mlrSpaceError      Error

hi def link mlrOperator        Operator

" ----------------------------------------------------------------
syn sync minlines=200

let b:current_syntax = "mlr"
