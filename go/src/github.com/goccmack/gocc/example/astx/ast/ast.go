package ast

import (
	"github.com/goccmack/gocc/example/astx/token"
)

type (
	StmtList []Stmt
	Stmt     string
)

func NewStmtList(stmt interface{}) (StmtList, error) {
	return StmtList{stmt.(Stmt)}, nil
}

func AppendStmt(stmtList, stmt interface{}) (StmtList, error) {
	return append(stmtList.(StmtList), stmt.(Stmt)), nil
}

func NewStmt(stmtList interface{}) (Stmt, error) {
	return Stmt(stmtList.(*token.Token).Lit), nil
}
