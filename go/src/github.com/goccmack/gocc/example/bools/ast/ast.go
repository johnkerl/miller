//Copyright 2012 Vastech SA (PTY) LTD
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
	"strconv"
	"strings"

	"github.com/goccmack/gocc/example/bools/token"
	"github.com/goccmack/gocc/example/bools/util"
)

type Attrib interface{}

type Val interface {
	Attrib
	Eval() bool
	String() string
}

type Op int

func (o Op) String() string {
	switch o {
	case OR:
		return "|"
	case AND:
		return "&"
	}
	return ""
}

const (
	NOOP = Op(0)
	OR   = Op(1)
	AND  = Op(2)
)

type BoolExpr struct {
	A  Val
	B  Val
	Op Op
}

func NewBoolAndExpr(a, b Attrib) (*BoolExpr, error) {
	return &BoolExpr{a.(Val), b.(Val), AND}, nil
}

func NewBoolOrExpr(a, b Attrib) (*BoolExpr, error) {
	return &BoolExpr{a.(Val), b.(Val), OR}, nil
}

func NewBoolGroupExpr(a Attrib) (*BoolExpr, error) {
	return &BoolExpr{a.(Val), nil, NOOP}, nil
}

func (this *BoolExpr) Eval() bool {
	switch this.Op {
	case OR:
		return this.A.Eval() || this.B.Eval()
	case AND:
		return this.A.Eval() && this.B.Eval()
	}
	return this.A.Eval()
}

func (this *BoolExpr) String() string {
	return this.A.String() + " " + this.Op.String() + " " + this.B.String()
}

type val bool

func (v val) Eval() bool {
	return bool(v)
}

func (v val) String() string {
	if v {
		return "true"
	}
	return "false"
}

var (
	TRUE  = val(true)
	FALSE = val(false)
)

type LessThanExpr struct {
	A int64
	B int64
}

func NewLessThanExpr(a, b Attrib) (*LessThanExpr, error) {
	aint, err := util.IntValue(a.(*token.Token).Lit)
	if err != nil {
		return nil, err
	}
	bint, err := util.IntValue(b.(*token.Token).Lit)
	if err != nil {
		return nil, err
	}
	return &LessThanExpr{aint, bint}, nil
}

func (this *LessThanExpr) Eval() bool {
	return this.A < this.B
}

func (this *LessThanExpr) String() string {
	return strconv.FormatInt(this.A, 10) + " < " + strconv.FormatInt(this.B, 10)
}

type SubStringExpr struct {
	A string
	B string
}

func NewSubStringExpr(a, b Attrib) (*SubStringExpr, error) {
	astr, err := strconv.Unquote(string(a.(*token.Token).Lit))
	if err != nil {
		return nil, err
	}
	bstr, err := strconv.Unquote(string(b.(*token.Token).Lit))
	if err != nil {
		return nil, err
	}
	return &SubStringExpr{astr, bstr}, nil
}

func (this *SubStringExpr) Eval() bool {
	return strings.Contains(this.B, this.A)
}

func (this *SubStringExpr) String() string {
	return this.A + " in " + this.B
}
