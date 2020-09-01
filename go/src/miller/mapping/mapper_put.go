package mapping

import (
	"fmt"

	"miller/containers"
	"miller/dsl"
	"miller/parsing/lexer"
	"miller/parsing/parser"
	"miller/runtime"
)

type MapperPut struct {
	ast         *dsl.AST
	interpreter *dsl.Interpreter
}

func NewMapperPut(dslString string) (*MapperPut, error) {
	ast, err := NewASTFromString(dslString)
	if err != nil {
		return nil, err
	}
	return &MapperPut{
		ast: ast,
		interpreter: dsl.NewInterpreter(),
	}, nil
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

func (this *MapperPut) Map(
	inrec *containers.Lrec,
	context *runtime.Context,
	outrecs chan<- *containers.Lrec,
) {
	if inrec != nil {
		// xxx maybe ast -> interpreter ctor
		outrec, err := this.interpreter.InterpretOnInputRecord(inrec, context, this.ast)
		if err != nil {
			// need echan or what?
			fmt.Println(err)
		} else {
			outrecs <- outrec
		}
	} else {
		outrecs <- nil
	}
}
