// ================================================================
// Just playing around -- nothing serious here.
// ================================================================

package repl

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/term"

	"miller/cliutil"
	"miller/dsl/cst"
	"miller/input"
	"miller/output"
	"miller/runtime"
	"miller/types"
	"miller/version"
)

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
		inputIsTerminal: inputIsTerminal,
		prompt1:         prompt1,
		prompt2:         prompt2,

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

func (this *Repl) printPrompt1() {
	if this.inputIsTerminal {
		fmt.Print(this.prompt1)
	}
}

func (this *Repl) printPrompt2() {
	if this.inputIsTerminal {
		fmt.Print(this.prompt2)
	}
}

// ----------------------------------------------------------------
func (this *Repl) HandleSession(istream *os.File) {
	fmt.Printf("Miller %s\n", version.STRING) // TODO: inhibit if mlr repl -q
	lineReader := bufio.NewReader(istream)

	for {
		this.printPrompt1()

		line, err := lineReader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			// TODO: lib.MlrExeName() and maybe %w and also remember verb name in ctor
			fmt.Fprintln(os.Stderr, "mlr repl:", err)
			os.Exit(1)
		}

		// This trims the trailing newline, as well as leading/trailing whitespace:
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "<" {
			this.handleMultiline(lineReader)
		} else if trimmedLine == ":quit" {
			break
		} else if this.handleNonDSLLine(trimmedLine) {
			// handled
		} else {
			// We need the non-trimmed line here since the DSL syntax for comments is '#.*\n'.
			err = this.handleDSLString(line, true) // xxx temp
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

// ----------------------------------------------------------------
// xxx comment "<" has already been seen. we read until ">". xxx "<"
func (this *Repl) handleMultiline(lineReader *bufio.Reader) {
	var buffer bytes.Buffer
	for {
		this.printPrompt2()

		line, err := lineReader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			// TODO: lib.MlrExeName() and maybe %w
			fmt.Fprintln(os.Stderr, "mlr repl:", err)
			os.Exit(1)
		}

		if strings.TrimSpace(line) == ">" {
			break
		}
		buffer.WriteString(line)
	}
	dslString := buffer.String()

	err := this.handleDSLString(dslString, false) // xxx temp
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
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
	if verb == ":help" || verb == "?" || verb == "help" || verb == ":h" {
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

	} else if verb == ":main" {
		_, err := this.cstRootNode.ExecuteMainBlock(this.runtimeState)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

	} else if verb == ":end" {
		err := this.cstRootNode.ExecuteEndBlocks(this.runtimeState)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

	} else if verb == ":blocks" {
		this.cstRootNode.ShowBlockReport()

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

		err = this.handleDSLString(dslString, false)
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
func (this *Repl) handleDSLString(dslString string, isReplImmediate bool) error {
	if strings.TrimSpace(dslString) == "" {
		return nil
	}

	astRootNode, err := this.BuildASTFromStringWithMessage(dslString)
	if err != nil {
		// Error message already printed out
		return err
	}

	this.cstRootNode.ResetForREPL()
	err = this.cstRootNode.IngestAST(astRootNode, this.isFilter, isReplImmediate) // TODO
	if err != nil {
		return err
	}
	err = this.cstRootNode.Resolve()
	if err != nil {
		return err
	}

	outrec, err := this.cstRootNode.ExecuteREPLImmediate(this.runtimeState)
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
