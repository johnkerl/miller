// ================================================================
// Just playing around -- nothing serious here.
// ================================================================

package auxents

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

//	"errors"
//	"fmt"
//	"io/ioutil"
//	"os"
//	"strings"
//
//	"miller/cliutil"
//	"miller/dsl"
//	"miller/dsl/cst"
//	"miller/lib"
//	"miller/parsing/lexer"
//	"miller/parsing/parser"
//	"miller/transforming"
//	"miller/types"

func replUsage(verbName string, o *os.File, exitCode int) {
	fmt.Fprintf(o, "Usage: %s %s with no arguments\n", mlrExeName(), verbName)
	os.Exit(exitCode)
}

func replMain(args []string) int {
	repl := NewRepl()
	repl.HandleSession(os.Stdin)
	return 0
}

// ================================================================
type Repl struct {
	// Idea: astPrint (parex or trindent)

	// recordReaderOptions *cliutil.TReaderOptions
	// recordWriterOptions *cliutil.TWriterOptions
	// astRootNode         *dsl.AST
	// cstRootNode         *cst.RootNode
	// cstState            *cst.State
	// outputChannel       chan<-*types.RecordAndContext

}

func NewRepl() *Repl {
	return &Repl{}
}

func (this *Repl) HandleSession(istream *os.File) {
	lineReader := bufio.NewReader(istream)

	for {
		line, err := lineReader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			// TODO: lib.MlrExeName()
			fmt.Fprintln(os.Stderr, "mlr repl:", err)
			os.Exit(1)
		}

		// This is how to do a chomp:
		line = strings.TrimRight(line, "\n")

		this.HandleLine(line)
	}
}

func (this *Repl) HandleLine(line string) {

	//	if printASTOnly {
	//		astRootNode, err := BuildASTFromStringWithMessage(dslString, false)
	//		if err == nil {
	//			if printASTSingleLine {
	//				astRootNode.PrintParexOneLine()
	//			} else {
	//				astRootNode.PrintParex()
	//			}
	//		// xxx astRootNode.Print()
	//			os.Exit(0)
	//		} else {
	//			// error message already printed out
	//			os.Exit(1)
	//		}
	//	}

}

//	astRootNode, err := BuildASTFromStringWithMessage(dslString, verbose)
//	if err != nil {
//		// Error message already printed out
//		return nil, err
//	}

//
//	cstRootNode, err := cst.Build(astRootNode, isFilter, recordWriterOptions)
//	cstState := cst.NewEmptyState()
//	if err != nil {
//		fmt.Fprintln(os.Stderr, err)
//		return nil, err
//	}

//func BuildASTFromStringWithMessage(dslString string, verbose bool) (*dsl.AST, error) {
//	astRootNode, err := BuildASTFromString(dslString)
//	if err != nil {
//		// Leave this out until we get better control over the error-messaging.
//		// At present it's overly parser-internal, and confusing. :(
//		// fmt.Fprintln(os.Stderr, err)
//		fmt.Fprintf(os.Stderr, "%s: cannot parse DSL expression.\n",
//			lib.MlrExeName())
//		if verbose {
//			fmt.Fprintln(os.Stderr, dslString)
//		}
//		fmt.Fprintln(os.Stderr, err)
//		return nil, err
//	} else {
//		return astRootNode, nil
//	}
//}

//func BuildASTFromString(dslString string) (*dsl.AST, error) {
//	theLexer := lexer.NewLexer([]byte(dslString))
//	theParser := parser.NewParser()
//	interfaceAST, err := theParser.Parse(theLexer)
//	if err != nil {
//		return nil, err
//	}
//	astRootNode := interfaceAST.(*dsl.AST)
//	return astRootNode, nil
//}

//	this.cstState.OutputChannel = outputChannel
//
//	inrec := inrecAndContext.Record
//	context := inrecAndContext.Context
//	if !inrecAndContext.EndOfStream {
//
//		if this.callCount == 1 {
//			this.cstState.Update(nil, &context)
//			err := this.cstRootNode.ExecuteBeginBlocks(this.cstState)
//			if err != nil {
//				fmt.Fprintln(os.Stderr, err)
//				os.Exit(1)
//			}
//			this.executedBeginBlocks = true
//		}
//
//		this.cstState.Update(inrec, &context)

//		// Execute the main block on the current input record
//		outrec, err := this.cstRootNode.ExecuteMainBlock(this.cstState)
//		if err != nil {
//			fmt.Fprintln(os.Stderr, err)
//			os.Exit(1)
//		}

//			outputChannel <- types.NewRecordAndContext(
//				outrec,
//				&context,
//			)
//
//		// Send all registered OutputHandlerManager instances the end-of-stream
//		// indicator.
//		this.cstRootNode.ProcessEndOfStream()
//
//		outputChannel <- types.NewEndOfStreamMarker(&context)
