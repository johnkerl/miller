package mapping

import (
	"fmt"
	"os"

	"miller/containers"
	"miller/dsl"
	"miller/parsing/lexer"
	"miller/parsing/parser"
	"miller/runtime"
)

type MapperPut struct {
	ast *dsl.AST
}

func NewMapperPut(dslString string) *MapperPut {
	theLexer := lexer.NewLexer([]byte(dslString))
	theParser := parser.NewParser()
	interfaceAST, err := theParser.Parse(theLexer)
	if err != nil {
		fmt.Println(err) // xxx error propagate to caller -- for all mapper constructors
		os.Exit(1)
	}
	ast := interfaceAST.(*dsl.AST)
	return &MapperPut{
		ast,
	}
}

func (this *MapperPut) Name() string {
	return "put"
}

func (this *MapperPut) Map(
	inrec *containers.Lrec,
	context *runtime.Context,
	outrecs chan<- *containers.Lrec,
) {
	if inrec != nil {
		outrecs <- inrec
	} else {
		outrecs <- nil
	}
}
