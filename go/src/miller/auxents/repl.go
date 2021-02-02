// ================================================================
// Just playing around -- nothing serious here.
// ================================================================

// TODO:
// oFlatSep = inrecAndContext.Context.OFLATSEP
// inrec.Flatten(oFlatSep)
// iFlatSep = inrecAndContext.Context.IFLATSEP
// inrec.Unflatten(iFlatSep)

package auxents

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/term"

	"miller/cliutil"
	"miller/dsl"
	"miller/dsl/cst"
	"miller/input"
	"miller/lib"
	"miller/output"
	"miller/parsing/lexer"
	"miller/parsing/parser"
	"miller/runtime"
	"miller/types"
)

// ================================================================
func replUsage(verbName string, o *os.File, exitCode int) {
	fmt.Fprintf(o, "Usage: %s %s with no arguments\n", mlrExeName(), verbName)
	os.Exit(exitCode)
}

// args are the full Miller command line: "mlr repl foo bar".
func replMain(args []string) int {
	//exeName := args[0]
	//replName := args[1]
	argc := len(args)
	argi := 2

	astPrintMode := ASTPrintNone
	options := cliutil.DefaultOptions()

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if !strings.HasPrefix(args[argi], "-") {
			break // No more flag options to process
		}

		if args[argi] == "-h" || args[argi] == "--help" {
			fmt.Println("help stub")
			os.Exit(0)
			// transformerPutUsage(os.Stdout, true, 0)

		} else if args[argi] == "-v" {
			astPrintMode = ASTPrintIndent
			argi++
		} else if args[argi] == "-d" {
			astPrintMode = ASTPrintParex
			argi++
		} else if args[argi] == "-D" {
			astPrintMode = ASTPrintParexOneLine
			argi++

		} else if cliutil.ParseReaderWriterOptions(args, argc, &argi, &options.ReaderOptions, &options.WriterOptions) {

		} else if cliutil.ParseReaderOptions(args, argc, &argi, &options.ReaderOptions) {

		} else if cliutil.ParseWriterOptions(args, argc, &argi, &options.WriterOptions) {

		} else {
			fmt.Println("help stub")
			os.Exit(1)
			// transformerPutUsage(os.Stderr, true, 1)
		}
	}

	repl, err := NewRepl(
		astPrintMode,
		&options,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	filenames := args[argi:]
	if len(filenames) > 0 {
		repl.openFiles(filenames)
	}

	repl.HandleSession(os.Stdin)
	return 0
}

// ================================================================
type ASTPrintMode int

const (
	ASTPrintNone ASTPrintMode = iota
	ASTPrintParex
	ASTPrintParexOneLine
	ASTPrintIndent
)

// ================================================================
type Repl struct {
	inputIsTerminal     bool
	prompt1             string
	prompt2             string
	doingMultilineInput bool

	astPrintMode ASTPrintMode
	isFilter     bool
	cstRootNode  *cst.RootNode

	options *cliutil.TOptions

	inputChannel chan *types.RecordAndContext
	errorChannel chan error
	recordReader input.IRecordReader
	recordWriter output.IRecordWriter

	runtimeState *runtime.State
}

// ----------------------------------------------------------------
func NewRepl(
	astPrintMode ASTPrintMode,
	options *cliutil.TOptions,
) (*Repl, error) {
	inputIsTerminal := term.IsTerminal(int(os.Stdin.Fd()))
	prompt1 := os.Getenv("MLR_REPL_PS1")
	if prompt1 == "" {
		prompt1 = "[mlr] "
	}
	prompt2 := os.Getenv("MLR_REPL_PS2")
	if prompt2 == "" {
		prompt2 = ""
	}
	doingMultilineInput := false

	recordReader := input.Create(&options.ReaderOptions)
	if recordReader == nil {
		return nil, errors.New("Input format not found: " + options.ReaderOptions.InputFileFormat)
	}

	recordWriter := output.Create(&options.WriterOptions)
	if recordWriter == nil {
		return nil, errors.New("Output format not found: " + options.WriterOptions.OutputFileFormat)
	}

	inrec := types.NewMlrmapAsRecord()
	context := types.NewContext(options)
	runtimeState := runtime.NewEmptyState()
	runtimeState.Update(inrec, context)
	runtimeState.FilterExpression = types.MlrvalFromVoid() // xxx comment

	return &Repl{
		inputIsTerminal:     inputIsTerminal,
		prompt1:             prompt1,
		prompt2:             prompt2,
		doingMultilineInput: doingMultilineInput,

		astPrintMode: ASTPrintNone,
		isFilter:     false,
		cstRootNode:  cst.NewEmptyRoot(&options.WriterOptions).WithRedefinableUDFS(),

		options:      options,
		inputChannel: nil,
		errorChannel: nil,
		recordReader: recordReader,
		recordWriter: recordWriter,

		runtimeState: runtimeState,
	}, nil
}

// ----------------------------------------------------------------
func (this *Repl) HandleSession(istream *os.File) {
	lineReader := bufio.NewReader(istream)
	dslString := ""

	for {
		if this.inputIsTerminal {
			if !this.doingMultilineInput {
				fmt.Print(this.prompt1)
			} else {
				fmt.Print(this.prompt2)
			}
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
	}
}

// ----------------------------------------------------------------
func (this *Repl) handleNonDSLLine(trimmedLine string) bool {
	args := strings.Fields(trimmedLine)
	if len(args) == 0 {
		return false
	}
	verb := args[0]
	if verb != "?" && verb != "help" && !strings.HasPrefix(verb, ":") {
		return false
	}
	// Make a lookup-table maybe
	if verb == ":help" || verb == "?" || verb == "help" {
		this.handleHelp(args)

	} else if verb == ":astprint" {
		this.handleASTPrint(args)

	} else if verb == ":load" {
		this.handleLoad(args)

	} else if verb == ":begin" {
		err := this.cstRootNode.ExecuteBeginBlocks(this.runtimeState)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

	} else if verb == ":end" {
		err := this.cstRootNode.ExecuteEndBlocks(this.runtimeState)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

	} else if verb == ":open" || verb == ":o" {
		this.handleOpen(args)

	} else if verb == ":read" || verb == ":r" {
		this.handleRead(args)

	} else if verb == ":write" || verb == ":w" {
		this.handleWrite(args)

		// xxx :write

	} else {
		fmt.Printf("Unrecognized command:%s\n", verb)
	}
	return true
}

func (this *Repl) handleHelp(args []string) {
	args = args[1:] // strip off verb
	if len(args) == 0 {
		fmt.Println("Options:")
		fmt.Println(":help repl")
		fmt.Println(":help functions")
		fmt.Println(":help {function name}, e.g. :help sec2gmt")
	} else {
		for _, arg := range args {
			if arg == "repl" {
				this.showREPLHelp()
			} else if arg == "functions" {
				cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionUsages(os.Stdout)
			} else {
				cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionUsage(arg, os.Stdout)
			}
		}
	}
}

// TODO: make this more like ':help repl explain' or somesuch
func (this *Repl) showREPLHelp() {
	fmt.Println(
		`Enter any Miller DSL expression.
Non-DSL commands (REPL-only statements) start with ':', such as ':help' or ':quit'.
Type ':help functions' for help with DSL functions; type ':help repl' for help with non-DSL expressions.

The input "record" by default is the empty map but you can do things like '$x=3',
or 'unset $y', or '$* = {"x": 3, "y": 4}' to populate it.

Enter '<' on a line by itself to enter multi-line mode, e.g. to enter a function definition;
enter '>' on a line by itself to exit multi-line mode.

In multi-line mode, semicolons are required between statements; otherwise they are not needed.
Non-assignment expressions, such as '7' or 'true', in mlr put are filter statements; here, they
are simply printed to the terminal, e.g. if you type '1+2', you will see '3'.

Examples, assuming the prompt is 'mlr: '

mlr: 1+2
3
mlr: x=3
mlr: y=4
mlr: x+y
7
mlr: <
func f(a,b) {
  return a**b
}
>
mlr: f(7,5)
16807
`)
}

func (this *Repl) handleASTPrint(args []string) {
	args = args[1:] // strip off verb
	if len(args) != 1 {
		fmt.Println("Need argument: see ':help :astprint'.")
		return
	}
	style := args[0]
	if style == "parex" {
		this.astPrintMode = ASTPrintParex
	} else if style == "parex1" {
		this.astPrintMode = ASTPrintParexOneLine
	} else if style == "indent" {
		this.astPrintMode = ASTPrintIndent
	} else if style == "none" {
		this.astPrintMode = ASTPrintNone
	} else {
		fmt.Printf("Unrecognized style %s: see ':help :astprint'.\n", style)
		return
	}
}

func (this *Repl) handleLoad(args []string) {
	args = args[1:] // strip off verb
	if len(args) < 1 {
		fmt.Println("Need filenames: see ':help :load'.")
		return
	}
	for _, filename := range args {

		dslBytes, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot load DSL expression from file \"%s\": ",
				filename)
			fmt.Println(err)
			return
		}
		dslString := string(dslBytes)

		err = this.HandleDSLString(dslString)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func (this *Repl) handleOpen(args []string) {
	this.openFiles(args[1:]) // strip off verb
}

func (this *Repl) openFiles(filenames []string) {
	this.inputChannel = make(chan *types.RecordAndContext, 10)
	// TODO: check for use
	this.errorChannel = make(chan error, 1)

	go this.recordReader.Read(
		filenames,
		*this.runtimeState.Context,
		this.inputChannel,
		this.errorChannel,
	)
}

func (this *Repl) handleRead(args []string) {
	if this.inputChannel == nil {
		fmt.Println("No open files")
		return
	}

	recordAndContext := <-this.inputChannel

	// Three things can come through:
	//
	// * End-of-stream marker
	// * Non-nil records to be printed
	// * Strings to be printed from put/filter DSL print/dump/etc
	//   statements. They are handled here rather than fmt.Println directly
	//   in the put/filter handlers since we want all print statements and
	//   record-output to be in the same goroutine, for deterministic
	//   output ordering.
	//
	// The first two are passed to the transformer. The third we send along
	// the output channel without involving the record-transformer, since
	// there is no record to be transformed.

	this.runtimeState.Update(recordAndContext.Record, this.runtimeState.Context)

	if recordAndContext.EndOfStream == true {
		// xxx put to recordwriter
		// xxx temp
		fmt.Println("End of record stream")
		this.inputChannel = nil
		this.errorChannel = nil
	} else if recordAndContext.Record == nil {
		fmt.Print(recordAndContext.OutputString)
	} else {
		fmt.Println(recordAndContext.Context.GetStatusString())
	}
}

func (this *Repl) handleWrite(args []string) {
	this.recordWriter.Write(this.runtimeState.Inrec, os.Stdout)
}

// ----------------------------------------------------------------
func (this *Repl) HandleDSLString(dslString string) error {

	astRootNode, err := this.BuildASTFromStringWithMessage(dslString)
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
	this.runtimeState.Inrec = outrec

	// xxx temp
	filterExpression := this.runtimeState.FilterExpression
	if filterExpression.IsVoid() {
		// nothing
	} else {
		fmt.Println(filterExpression.String())
	}
	this.runtimeState.FilterExpression = types.MlrvalFromVoid()

	return nil
}

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
