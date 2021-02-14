// ================================================================
// This is the interface between the REPL and the DSL-to-AST parser.
// ================================================================

package repl

import (
	"fmt"
	"os"

	"miller/src/dsl"
	"miller/src/lib"
	"miller/src/parsing/lexer"
	"miller/src/parsing/parser"
)

// ----------------------------------------------------------------
func (this *Repl) BuildASTFromStringWithMessage(
	dslString string,
) (*dsl.AST, error) {
	astRootNode, err := this.BuildASTFromString(dslString)
	if err != nil {
		// Leave this out until we get better control over the error-messaging.
		// At present it's overly parser-internal, and confusing. :(
		// fmt.Fprintln(os.Stderr, err)
		fmt.Fprintf(os.Stderr, "%s: cannot parse DSL expression.\n",
			lib.MlrExeName())
		if this.astPrintMode != ASTPrintNone {
			fmt.Fprintln(os.Stderr, dslString)
		}
		return nil, err
	} else {
		if this.astPrintMode == ASTPrintParex {
			astRootNode.PrintParex()
		} else if this.astPrintMode == ASTPrintParexOneLine {
			astRootNode.PrintParexOneLine()
		} else if this.astPrintMode == ASTPrintIndent {
			astRootNode.Print()
		}

		return astRootNode, nil
	}
}

// ----------------------------------------------------------------
func (this *Repl) BuildASTFromString(dslString string) (*dsl.AST, error) {
	theLexer := lexer.NewLexer([]byte(dslString))
	theParser := parser.NewParser()
	interfaceAST, err := theParser.Parse(theLexer)
	if err != nil {
		return nil, err
	}
	astRootNode := interfaceAST.(*dsl.AST)
	return astRootNode, nil
}
