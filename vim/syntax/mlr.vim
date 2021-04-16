" Copyright 2021 John Kerl. All rights reserved.
" Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.
"
" mlr.vim: Vim syntax file for the Miller DSL.

if exists("b:current_syntax")
  finish
endif

syn case match

" ----------------------------------------------------------------
" Goal: map the lexical elements of the Miller DSL grammar
"   https://github.com/johnkerl/miller/blob/main/go/src/parsing/mlr.bnf
" to Vim syntax options
"   http://vimdoc.sourceforge.net/htmldoc/syntax.html
" ----------------------------------------------------------------

" ----------------------------------------------------------------
syn region mlrComment          start="#" end="$"
syn region mlrString           start=+"+ skip=+\\\\\|\\"+ end=+"+ contains=@mlrStringGroup
syn match   mlrDecimalInt      "\<-\=\(0\|[1-9]_\?\(\d\|\d\+_\?\d\+\)*\)\%([Ee][-+]\=\d\+\)\=\>"
syn match   mlrHexadecimalInt  "\<-\=0[xX]_\?\(\x\+_\?\)\+\>"
syn match   mlrBinaryInt       "\<-\=0[bB]_\?\([01]\+_\?\)\+\>"
syn match   mlrFloat           "\<-\=\d\+\.\d*\%([Ee][-+]\=\d\+\)\=\>"
syn match   mlrFloat           "\<-\=\.\d\+\%([Ee][-+]\=\d\+\)\=\>"
syn keyword mlrConstant        M_PI M_E
syn keyword mlrBoolean         true false
syn keyword mlrContextVariable IPS IFS IRS OPS OFS ORS OFLATSEP NF NR FNR FILENAME FILENUM
syn keyword mlrENV             ENV
syn match   mlrOperator        /[-+%<>!&|^*=]=\?/
syn match   mlrOperator        /\/\%(=\|\ze[^/*]\)/
syn match   mlrOperator        /\%(<<\|>>\|&^\)=\?/
syn match   mlrOperator        /:=\|||\|<-\|++\|--/
" TODO: more operators
syn keyword mlrKeyword         begin do elif else end filter for if in while break continue return
syn keyword mlrKeyword         func subr call dump edump
syn keyword mlrKeyword         emit emitp emitf eprint eprintn print printn tee stdout stderr unset null
syn keyword mlrType            arr bool float int map num str var
syn match   mlrFieldName       /\(\$[_a-zA-Z0-9]+\)/
syn match   mlrFieldName       /\(\$\*\)/
syn match   mlrOosvarName      /\(@[_a-zA-Z0-9]+\)/
syn match   mlrOosvarName      /\(@\*\)/
syn match   mlrIdentifier      /\(\<[a-zA-Z_][a-zA-Z_0-9]*\>\)/
syn match   mlrFunctionCall    /\w\+\ze(/ contains=mlrBuiltins,mlrDeclaration

syn region  mlrParen           start='(' end=')' transparent
syn region  mlrBlock           start="{" end="}" transparent

" Trailing whitespace; space-tab
syn match   mlrSpaceError      display excludenl "\s\+$"
syn match   mlrSpaceError      display " \+\t"me=e-1

" ----------------------------------------------------------------

hi def link mlrComment         Comment
hi def link mlrString          String
hi def link mlrDecimalInt      Integer
hi def link mlrHexadecimalInt  Integer
hi def link mlrBinaryInt       Integer
hi def link Integer            Number
hi def link mlrFloat           Float
hi def link mlrConstant        Constant
hi def link mlrBoolean         Boolean
hi def link mlrContextVariable Keyword
hi def link mlrENV             Keyword
hi def link mlrOperator        Operator
hi def link mlrKeyword         Keyword
hi def link mlrType            Type
hi def link mlrFieldName       Special
hi def link mlrOosvarName      Special
hi def link mlrIdentifier      Identifier
hi def link mlrFunctionCall    Type
hi def link mlrSpaceError      Error

" ----------------------------------------------------------------
syn sync minlines=200

let b:current_syntax = "mlr"

" vim: sw=2 ts=2 et
