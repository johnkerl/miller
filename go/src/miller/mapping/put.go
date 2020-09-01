package mapping

import (
	"fmt"
	"os"

	"miller/clitypes"
	"miller/containers"
	"miller/dsl"
	"miller/parsing/lexer"
	"miller/parsing/parser"
	"miller/runtime"
)

var MapperPutSetup = MapperSetup{
	Verb:         "put",
	ParseCLIFunc: mapperPutParseCLIFunc,
	UsageFunc:    mapperPutUsageFunc,
	IgnoresInput: false,
}

func mapperPutParseCLIFunc(
	pargi *int,
	argc int,
	args []string,
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) IRecordMapper {
	if argc - *pargi < 2 {
		return nil
	}
	// xxx temp hack
	dslString := args[*pargi + 1]
	*pargi += 2

	mapper, _ := NewMapperPut(dslString)
	return mapper
}

func mapperPutUsageFunc(
	o *os.File,
	argv0 string,
	verb string,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o, "TODO: un-stub this help function.\n");
}

// ----------------------------------------------------------------
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
