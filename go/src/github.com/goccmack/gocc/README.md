# New
Have a look at [https://github.com/goccmack/gogll](https://github.com/goccmack/gogll) for scannerless GLL parser generation.
# Gocc

[![Build Status](https://travis-ci.org/goccmack/gocc.svg?branch=master)](https://travis-ci.org/goccmack/gocc)

## Introduction

Gocc is a compiler kit for Go written in Go.

Gocc generates lexers and parsers or stand-alone DFAs or parsers from a BNF.

Lexers are DFAs, which recognise regular languages. Gocc lexers accept UTF-8 input.

Gocc parsers are PDAs, which recognise LR-1 languages. Optional LR1 conflict handling automatically resolves shift / reduce and reduce / reduce conflicts.

Generating a lexer and parser starts with creating a bnf file. Action expressions embedded in the BNF allows the user to specify semantic actions for syntax productions.

For complex applications the user typically uses an abstract syntax tree (AST) to represent the derivation of the input. The user provides a set of functions to construct the AST, which are called from the action expressions specified in the BNF.

See the [README](example/bools/README) for an included example.

[User Guide (PDF): Learn You a gocc for Great Good](https://raw.githubusercontent.com/goccmack/gocc/master/doc/gocc_user_guide.pdf) (gocc3 user guide will be published shortly)

## Installation

* First download and Install Go From http://golang.org/
* Setup your GOPATH environment variable.
* Next in your command line run: go get github.com/goccmack/gocc (go get will git clone gocc into GOPATH/src/github.com/goccmack/gocc and run go install)
* Alternatively clone the source: https://github.com/goccmack/gocc . Followed by go install github.com/goccmack/gocc
* Finally make sure that the bin folder where the gocc binary is located is in your PATH environment variable.

## Getting Started

Once installed start by creating your BNF in a package folder.

For example GOPATH/src/foo/bar.bnf:

```
/* Lexical Part */

id : 'a'-'z' {'a'-'z'} ;

!whitespace : ' ' | '\t' | '\n' | '\r' ;

/* Syntax Part */

<< import "foo/ast" >>

Hello:  "hello" id << ast.NewWorld($1) >> ;
```

Next to use gocc, run:

```sh
cd $GOPATH/src/foo
gocc bar.bnf
```

This will generate a scanner, parser and token package inside GOPATH/src/foo Following times you might only want to run gocc without the scanner flag, since you might want to start making the scanner your own. Gocc is after all only a parser generator even if the default scanner is quite useful.

Next create ast.go file at $GOPATH/src/foo/ast with the following contents:

```go
package ast

import (
    "foo/token"
)

type Attrib interface {}

type World struct {
    Name string
}

func NewWorld(id Attrib) (*World, error) {
    return &World{string(id.(*token.Token).Lit)}, nil
}

func (this *World) String() string {
    return "hello " + this.Name
}
```

Finally we want to parse a string into the ast, so let us write a test at $GOPATH/src/foo/test/parse_test.go with the following contents:

```go
package test

import (
    "foo/ast"
    "foo/lexer"
    "foo/parser"
    "testing"
)

func TestWorld(t *testing.T) {
    input := []byte(`hello gocc`)
    lex := lexer.NewLexer(input)
    p := parser.NewParser()
    st, err := p.Parse(lex)
    if err != nil {
        panic(err)
    }
    w, ok := st.(*ast.World)
    if !ok {
        t.Fatalf("This is not a world")
    }
    if w.Name != `gocc` {
        t.Fatalf("Wrong world %v", w.Name)
    }
}
```

Finally run the test:

```sh
cd $GOPATH/src/foo/test
go test -v
```

You have now created your first grammar with gocc. This should now be relatively easy to change into the grammar you actually want to create or an existing LR1 grammar you would like to parse.

## BNF

The Gocc BNF is specified [here](spec/gocc2.ebnf)

An example bnf with action expressions can be found [here](example/bools/example.bnf)

## Action Expressions and AST

An action expression is specified as "<", "<", goccExpressionList , ">", ">" . The goccExpressionList is equivalent to a [goExpressionList](https://golang.org/ref/spec#ExpressionList). This expression list should return an Attrib and an error. Where Attrib is:

```go
type Attrib interface {}
```

Also parsed elements of the corresponding bnf rule can be represented in the expressionList as "$", digit.

Some action expression examples:

```
<< $0, nil >>
<< ast.NewFoo($1) >>
<< ast.NewBar($3, $1) >>
<< ast.TRUE, nil >>
```

Contants, functions, etc. that are returned or called should be programmed by the user in his ast (Abstract Syntax Tree) package. The ast package requires that you define your own Attrib interface as shown above. All parameters passed to functions will be of this type.

Some example of functions:

```go
func NewFoo(a Attrib) (*Foo, error) { ... }
func NewBar(a, b Attrib) (*Bar, error) { ... }
```

An example of an ast can be found [here](example/bools/ast/ast.go)

## Release Notes for gocc 2.1

### Changes

1. no_lexer option added to suppress generation of lexer. See the user guide.

2. Unreachable code removed from generated code.

### Bugs fixed:

1. gocc 2.1 does not support string_lit symbols with the same value as production names of the BNF. E.g. (t2.bnf):

```
A : "a" | "A" ;
```

string_lit "A" is not allowed.

Previously gocc silently ignored the conflicting string_lit. Now it generates an ugly panic:

```
$ gocc t2.bnf
panic: string_lit "A" conflicts with production name A
```

This issue will be properly resolved in a future release.

## Users

These projects use gocc:

* [gogo](https://github.com/shivansh/gogo) - [BNF file](https://github.com/shivansh/gogo/blob/master/src/lang.bnf) - a Go to MIPS compiler written in Go
* [gonum/gonum](https://github.com/gonum/gonum) - [BNF file](https://github.com/gonum/gonum/blob/master/graph/formats/dot/internal/dot.bnf) - DOT decoder (part of the graph library of Gonum)
* [llir/llvm](https://github.com/llir/llvm) - [BNF file](https://github.com/llir/llvm/blob/master/asm/internal/ll.bnf) - LLVM IR library in pure Go
* [mewmew/uc](https://github.com/mewmew/uc) - [BNF file](https://github.com/mewmew/uc/blob/master/gocc/uc.bnf) - A compiler for the µC language
* [gographviz](https://github.com/awalterschulze/gographviz) - [BNF file](https://github.com/awalterschulze/gographviz/blob/master/dot.bnf) - Parses the Graphviz DOT language in golang 
* [katydid/relapse](http://katydid.github.io/) - [BNF file](https://github.com/katydid/katydid/blob/master/relapse/bnf/all.bnf) - Encoding agnostic validation language
