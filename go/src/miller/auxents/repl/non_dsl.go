// ================================================================
// Just playing around -- nothing serious here.
// ================================================================

package repl

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"miller/dsl/cst"
	"miller/types"
)

// ----------------------------------------------------------------
type tHandlerFunc func(repl *Repl, args []string) bool
type tUsageFunc func(repl *Repl)
type handlerInfo struct {
	handlerFunc tHandlerFunc
	usageFunc   tUsageFunc
	verbNames   []string
}

// TODO: tabularize arg-count expected and common usage if not.

var handlerLookupTable = []handlerInfo{
	{handlerFunc: handleLoad, usageFunc: usageLoad, verbNames: []string{":l", ":load"}},
	{handlerFunc: handleOpen, usageFunc: usageOpen, verbNames: []string{":o", ":open"}},
	{handlerFunc: handleRead, usageFunc: usageRead, verbNames: []string{":r", ":read"}},
	{handlerFunc: handleWrite, usageFunc: usageWrite, verbNames: []string{":w", ":write"}},
	{handlerFunc: handleBegin, usageFunc: usageBegin, verbNames: []string{":b", ":begin"}},
	{handlerFunc: handleMain, usageFunc: usageMain, verbNames: []string{":m", ":main"}},
	{handlerFunc: handleEnd, usageFunc: usageEnd, verbNames: []string{":e", ":end"}},
	{handlerFunc: handleASTPrint, usageFunc: usageASTPrint, verbNames: []string{":astprint"}},
	{handlerFunc: handleBlocks, usageFunc: usageBlocks, verbNames: []string{":blocks"}},
	{handlerFunc: handleHelp, usageFunc: usageHelp, verbNames: []string{":h", ":help"}},
}

// ----------------------------------------------------------------
// No hash-table acceleration; things here are small, and the tool is interactive.
func (this *Repl) findHandler(verbName string) *handlerInfo {
	for _, entry := range handlerLookupTable {
		for _, entryVerbName := range entry.verbNames {
			if entryVerbName == verbName {
				return &entry
			}
		}
	}
	return nil
}

// ----------------------------------------------------------------
func (this *Repl) handleNonDSLLine(trimmedLine string) bool {
	args := strings.Fields(trimmedLine)
	if len(args) == 0 {
		return false
	}
	verbName := args[0]
	// We handle all single lines starting with a colon.
	if !strings.HasPrefix(verbName, ":") {
		return false
	}

	handler := this.findHandler(verbName)
	if handler == nil {
		fmt.Printf("Non-DSL verb %s not found.\n", verbName)
		return true
	}

	if !handler.handlerFunc(this, args) {
		handler.usageFunc(this)
	}
	return true
}

// ----------------------------------------------------------------
func usageLoad(this *Repl) {
	fmt.Println("Usage: :load {one or more filenames containing Miller DSL statements}")
}

func handleLoad(this *Repl, args []string) bool {
	args = args[1:] // strip off verb
	if len(args) < 1 {
		return false
	}
	for _, filename := range args {
		dslBytes, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot load DSL expression from file \"%s\": ",
				filename)
			fmt.Println(err)
			return true
		}
		dslString := string(dslBytes)

		err = this.handleDSLStringBulk(dslString)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	return true
}

// ----------------------------------------------------------------
func usageOpen(this *Repl) {
	fmt.Printf(
		"Usage: :load {one or more data-file names, in the format specifed by %s %s\n",
		this.exeName, this.replName,
	)
}

func handleOpen(this *Repl, args []string) bool {
	this.openFiles(args[1:]) // strip off verb
	return true
}

// Also invoked from the main entry-point, hence split out as a separate method.
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

// ----------------------------------------------------------------
func usageRead(this *Repl) {
	fmt.Println("Usage: :read with no arguments.")
}

func handleRead(this *Repl, args []string) bool {
	if this.inputChannel == nil {
		fmt.Println("No open files")
		return true
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
	return true
}

// ----------------------------------------------------------------
func usageWrite(this *Repl) {
	fmt.Println("Usage: :write with no arguments.")
}
func handleWrite(repl *Repl, args []string) bool {
	if len(args) != 1 {
		return false
	}
	repl.recordWriter.Write(repl.runtimeState.Inrec, os.Stdout)
	return true
}

// ----------------------------------------------------------------
func usageBegin(this *Repl) {
	fmt.Println("Usage: :begin with no arguments.")
}
func handleBegin(this *Repl, args []string) bool {
	if len(args) != 1 {
		return false
	}
	err := this.cstRootNode.ExecuteBeginBlocks(this.runtimeState)
	if err != nil {
		fmt.Println(err)
	}
	return true
}

// ----------------------------------------------------------------
func usageMain(this *Repl) {
	fmt.Println("Usage: :main with no arguments.")
}
func handleMain(this *Repl, args []string) bool {
	_, err := this.cstRootNode.ExecuteMainBlock(this.runtimeState)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return true
}

// ----------------------------------------------------------------
func usageEnd(this *Repl) {
	fmt.Println("Usage: :end with no arguments.")
}
func handleEnd(this *Repl, args []string) bool {
	err := this.cstRootNode.ExecuteEndBlocks(this.runtimeState)
	if err != nil {
		fmt.Println(err)
	}
	return true
}

// ----------------------------------------------------------------
func usageASTPrint(this *Repl) {
	fmt.Println("Need argument: see ':help :astprint'.")
}
func handleASTPrint(this *Repl, args []string) bool {
	args = args[1:] // strip off verb
	if len(args) != 1 {
		return false
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
	}
	return true
}

// ----------------------------------------------------------------
func usageBlocks(this *Repl) {
	fmt.Println("Usage: :blocks with no arguments.")
}
func handleBlocks(this *Repl, args []string) bool {
	this.cstRootNode.ShowBlockReport()
	return true
}

// ----------------------------------------------------------------
func usageHelp(this *Repl) {
	fmt.Println("TODO: metahelp is TBD.")
}

func handleHelp(this *Repl, args []string) bool {
	args = args[1:] // strip off verb
	if len(args) == 0 {
		fmt.Println("Options:")
		fmt.Println(":help repl")
		fmt.Println(":help functions")
		fmt.Println(":help {function name}, e.g. :help sec2gmt")
	} else {
		for _, arg := range args {
			if arg == "repl" {
				showREPLHelp()
			} else if arg == "functions" {
				cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionUsages(os.Stdout)
			} else {
				cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionUsage(arg, os.Stdout)
			}
		}
	}
	return true
}

// TODO: make this more like ':help repl explain' or somesuch
func showREPLHelp() {
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
