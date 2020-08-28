package ast

import (
	"github.com/goccmack/gocc/example/errorrecovery/token"
)

type (
	StmtList []interface{}
	Stmt     string
)

func NewStmtList(stmt interface{}) (StmtList, error) {
	return StmtList{stmt}, nil
}

func AppendStmt(stmtList, stmt interface{}) (StmtList, error) {
	return append(stmtList.(StmtList), stmt), nil
}

func NewStmt(stmtList interface{}) (Stmt, error) {
	return Stmt(stmtList.(*token.Token).Lit), nil
}
