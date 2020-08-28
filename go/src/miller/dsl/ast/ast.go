package ast

import (
	"miller/parsing/token"
)

type (
	StatementList []Statement
	Statement     string
)

func NewStatementList(stmt interface{}) (StatementList, error) {
	return StatementList{stmt.(Statement)}, nil
}

func AppendStatement(stmtList, stmt interface{}) (StatementList, error) {
	return append(stmtList.(StatementList), stmt.(Statement)), nil
}

func NewStatement(stmtList interface{}) (Statement, error) {
	return Statement(stmtList.(*token.Token).Lit), nil
}
