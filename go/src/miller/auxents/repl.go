// ================================================================
// Just playing around -- nothing serious here.
// ================================================================

package auxents

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"miller/cliutil"
	"miller/dsl"
	"miller/dsl/cst"
	"miller/lib"
	"miller/output"
	"miller/parsing/lexer"
	"miller/parsing/parser"
	"miller/runtime"
	"miller/types"
)

func replUsage(verbName string, o *os.File, exitCode int) {
	fmt.Fprintf(o, "Usage: %s %s with no arguments\n", mlrExeName(), verbName)
	os.Exit(exitCode)
}

func replMain(args []string) int {
	repl, err := NewRepl()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	repl.HandleSession(os.Stdin)
	return 0
}

// ================================================================
type Repl struct {
	prompt1             string
	prompt2             string
	doingMultilineInput bool

	verboseASTParse bool
	isFilter        bool
	cstRootNode     *cst.RootNode

	options      cliutil.TOptions
	context      *types.Context
	recordWriter output.IRecordWriter

	runtimeState *runtime.State
}

// ----------------------------------------------------------------
func NewRepl() (*Repl, error) {
	// https://pkg.go.dev/golang.org/x/term#IsTerminal
	prompt1 := "[mlr] "
	prompt2 := ""
	doingMultilineInput := false

	options := cliutil.DefaultOptions()
	inrec := types.NewMlrmapAsRecord()
	context := types.NewContext(&options)
	recordWriter := output.Create(&options.WriterOptions)
	if recordWriter == nil {
		return nil, errors.New("Output format not found: " + options.WriterOptions.OutputFileFormat)
	}

	runtimeState := runtime.NewEmptyState()
	runtimeState.Update(inrec, context)

	// TODO: empty record

	return &Repl{
		prompt1:             prompt1,
		prompt2:             prompt2,
		doingMultilineInput: doingMultilineInput,

		verboseASTParse: false,
		isFilter:        false,
		cstRootNode:     cst.NewEmptyRoot(&options.WriterOptions),

		options:      options,
		context:      context,
		recordWriter: recordWriter,

		runtimeState: runtimeState,
	}, nil
}

// ----------------------------------------------------------------
func (this *Repl) HandleSession(istream *os.File) {
	lineReader := bufio.NewReader(istream)
	dslString := ""

	for {
		if !this.doingMultilineInput {
			fmt.Print(this.prompt1)
		} else {
			fmt.Print(this.prompt2)
		}

		line, err := lineReader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			// TODO: lib.MlrExeName()
			fmt.Fprintln(os.Stderr, "mlr repl:", err)
			os.Exit(1)
		}

		// This trims the trailing newline, as well as leading/trailing
		// whitespace:
		trimmedLine := strings.TrimSpace(line)

		if !this.doingMultilineInput {
			if trimmedLine == "<" {
				this.doingMultilineInput = true
				dslString = ""
				continue
			} else if trimmedLine == ":quit" {
				break
			} else if this.handleNonDSLLine(trimmedLine) {
				continue
			} else {
				dslString = line
			}

		} else {
			if trimmedLine == ">" {
				this.doingMultilineInput = false
			} else {
				dslString += line
				continue
			}
		}

		err = this.HandleDSLString(dslString)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// TODO:
		// UTs
		// take filenames & reader/writer options from CLI -- ?
		// ?/:help -- w/o and w/ function name
		// :some flag settings ... ? :--o json -- ?
		// :open file (or stdin?)
		// :next ....
		// :close
		// :continue
		// :break conditions -- ?
		// ^C -- ?!?
		// auto-start/auto-end multiline mode w/ func/subr defs?
		// [DONE] < start multiline
		// [DONE] > end multiline
	}
}

// ----------------------------------------------------------------
func (this *Repl) handleNonDSLLine(trimmedLine string) bool {
	args := strings.Fields(trimmedLine)
	if len(args) == 0 {
		return false
	}
	verb := args[0]
	// Make a lookup-table maybe
	if verb == ":help" || verb == "?" {
		this.handleHelp(args)
		return true
	}

	// TODO: :help/? {funcname/keyword}

	return false
}

// ----------------------------------------------------------------
func (this *Repl) handleHelp(args []string) {
	args = args[1:] // strip off verb
	if len(args) == 0 {
		cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionUsages(os.Stdout)
	} else {
		for _, arg := range args {
			cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionUsage(arg, os.Stdout)
		}
	}
	//verb := args[0]
}

// ----------------------------------------------------------------
func (this *Repl) HandleDSLString(dslString string) error {

	astRootNode, err := this.BuildASTFromStringWithMessage(dslString, this.verboseASTParse)
	if err != nil {
		// Error message already printed out
		return err
	}

	this.cstRootNode.ResetForREPL()
	err = this.cstRootNode.Build(astRootNode, this.isFilter)
	if err != nil {
		return err
	}

	outrec, err := this.cstRootNode.ExecuteREPLExperimental(this.runtimeState)
	if err != nil {
		return err
	}

	if false { // not interesting ... maybe with a CLI flag ...
		this.recordWriter.Write(outrec, os.Stdout)
	}

	return nil
}

// ----------------------------------------------------------------
func (this *Repl) BuildASTFromStringWithMessage(
	dslString string,
	verbose bool,
) (*dsl.AST, error) {
	astRootNode, err := this.BuildASTFromString(dslString)
	if err != nil {
		// Leave this out until we get better control over the error-messaging.
		// At present it's overly parser-internal, and confusing. :(
		// fmt.Fprintln(os.Stderr, err)
		fmt.Fprintf(os.Stderr, "%s: cannot parse DSL expression.\n",
			lib.MlrExeName())
		if verbose {
			fmt.Fprintln(os.Stderr, dslString)
		}
		return nil, err
	} else {

		//	if printASTSingleLine {
		//		astRootNode.PrintParexOneLine()
		//	} else if xxx {
		//		astRootNode.PrintParex()
		//	} else {
		//		astRootNode.Print()
		//	}

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
