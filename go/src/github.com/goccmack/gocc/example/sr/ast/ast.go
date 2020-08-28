package ast

import (
	"github.com/goccmack/gocc/example/sr/token"
)

type (
	Attrib interface{}

	Stmt interface {
		Equals(Stmt) bool
		MatchIf(cnd string, s Stmt) bool
		MatchIfElse(cnd string, s1, s2 Stmt) bool
		MatchId(s string) bool
		String() string
	}

	If struct {
		C string
		S Stmt
	}

	IfElse struct {
		C      string
		S1, S2 Stmt
	}

	IdStmt string
)

func NewIf(cnd, stmt Attrib) (ifs *If) {
	ifs = &If{
		C: string(cnd.(*token.Token).Lit),
		S: stmt.(Stmt),
	}
	return
}

func NewIfElse(cnd, stmt1, stmt2 Attrib) (ies *IfElse) {
	ies = &IfElse{
		C:  string(cnd.(*token.Token).Lit),
		S1: stmt1.(Stmt),
		S2: stmt2.(Stmt),
	}
	return
}

func NewIdStmt(s Attrib) (is IdStmt) {
	is = IdStmt(string(s.(*token.Token).Lit))
	return
}

func (this *If) Equals(that Stmt) bool {
	if this == that {
		return true
	}
	that1, ok := that.(*If)
	if !ok {
		return false
	}
	return this.C == that1.C && this.S.Equals(that1.S)
}

func (this *IfElse) Equals(that Stmt) bool {
	if this == that {
		return true
	}
	that1, ok := that.(*IfElse)
	if !ok {
		return false
	}
	return this.C == that1.C && this.S1.Equals(that1.S1) && this.S2.Equals(that1.S2)
}

func (this IdStmt) Equals(that Stmt) bool {
	if this == "" {
		return false
	}
	that1, ok := that.(IdStmt)
	if !ok {
		return false
	}
	return this == that1
}

func (this *If) MatchIf(cnd string, s Stmt) bool {
	return this.C == cnd && this.S.Equals(s)
}

func (this *IfElse) MatchIf(cnd string, s Stmt) bool {
	return false
}

func (this IdStmt) MatchIf(cnd string, s Stmt) bool {
	return false
}

func (this *If) MatchIfElse(cnd string, s1, s2 Stmt) bool {
	return false
}

func (this *IfElse) MatchIfElse(cnd string, s1, s2 Stmt) bool {
	return this.C == cnd && this.S1.Equals(s1) && this.S2.Equals(s2)
}

func (this IdStmt) MatchIfElse(cnd string, s1, s2 Stmt) bool {
	return false
}

func (this *If) MatchId(s string) bool {
	return false
}

func (this *IfElse) MatchId(s string) bool {
	return false
}

func (this IdStmt) MatchId(s string) bool {
	return this == IdStmt(s)
}

func (this *If) String() string {
	return "*If: " + string(this.C) + this.S.String()
}

func (this *IfElse) String() string {
	return "*IfElse: " + this.C + this.S1.String() + this.S2.String()
}

func (this IdStmt) String() string {
	return "IdStmt: " + string(this)
}
