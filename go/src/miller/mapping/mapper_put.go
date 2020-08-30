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
	ast, err := NewASTFromString(dslString)
	if err != nil {
		fmt.Println(err) // xxx error propagate to caller -- for all mapper constructors
		os.Exit(1)
	}
	return &MapperPut{
		ast,
	}
}

// xxx note (package cycle) why not a dsl.AST constructor :(
// xxx maybe split out dsl into two package ... and/or put the ast.go into miller/parsing -- ?
//   depends on TBD split-out of AST and CST ...
func NewASTFromString(dslString string) (*dsl.AST, error) {
	theLexer := lexer.NewLexer([]byte(dslString))
	theParser := parser.NewParser()
	interfaceAST, err := theParser.Parse(theLexer)
	if err != nil {
		return nil, err
	}
	ast := interfaceAST.(*dsl.AST)
	return ast, nil
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
