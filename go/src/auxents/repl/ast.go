// ================================================================
// This is the interface between the REPL and the DSL-to-AST parser.
// ================================================================

package repl

import (
	"fmt"
	"os"

	"miller/src/dsl"
	"miller/src/parsing/lexer"
	"miller/src/parsing/parser"
)

// ----------------------------------------------------------------
func (repl *Repl) BuildASTFromStringWithMessage(
	dslString string,
) (*dsl.AST, error) {
	astRootNode, err := repl.BuildASTFromString(dslString)
	if err != nil {
		// Leave this out until we get better control over the error-messaging.
		// At present it's overly parser-internal, and confusing. :(
		// fmt.Fprintln(os.Stderr, err)
		fmt.Fprintf(os.Stderr, "%s: cannot parse DSL expression.\n",
			"mlr")
		if repl.astPrintMode != ASTPrintNone {
			fmt.Fprintln(os.Stderr, dslString)
		}
		return nil, err
	} else {
		if repl.astPrintMode == ASTPrintParex {
			astRootNode.PrintParex()
		} else if repl.astPrintMode == ASTPrintParexOneLine {
			astRootNode.PrintParexOneLine()
		} else if repl.astPrintMode == ASTPrintIndent {
			astRootNode.Print()
		}

		return astRootNode, nil
	}
}

// ----------------------------------------------------------------
func (repl *Repl) BuildASTFromString(dslString string) (*dsl.AST, error) {
	theLexer := lexer.NewLexer([]byte(dslString))
	theParser := parser.NewParser()
	interfaceAST, err := theParser.Parse(theLexer)
	if err != nil {
		return nil, err
	}
	astRootNode := interfaceAST.(*dsl.AST)
	return astRootNode, nil
}
