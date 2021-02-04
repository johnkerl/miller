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
	verbNames   []string
	handlerFunc tHandlerFunc
	usageFunc   tUsageFunc
}

// We get a Golang "initialization loop" if this is defined statically. So, we
// use a "package init" function.
var handlerLookupTable = []handlerInfo{}

func init() {
	handlerLookupTable = []handlerInfo{
		{verbNames: []string{":l", ":load"}, handlerFunc: handleLoad, usageFunc: usageLoad},
		{verbNames: []string{":o", ":open"}, handlerFunc: handleOpen, usageFunc: usageOpen},
		{verbNames: []string{":r", ":read"}, handlerFunc: handleRead, usageFunc: usageRead},
		{verbNames: []string{":w", ":write"}, handlerFunc: handleWrite, usageFunc: usageWrite},
		{verbNames: []string{":b", ":begin"}, handlerFunc: handleBegin, usageFunc: usageBegin},
		{verbNames: []string{":m", ":main"}, handlerFunc: handleMain, usageFunc: usageMain},
		{verbNames: []string{":e", ":end"}, handlerFunc: handleEnd, usageFunc: usageEnd},
		{verbNames: []string{":astprint"}, handlerFunc: handleASTPrint, usageFunc: usageASTPrint},
		{verbNames: []string{":blocks"}, handlerFunc: handleBlocks, usageFunc: usageBlocks},
		{verbNames: []string{":h", ":help"}, handlerFunc: handleHelp, usageFunc: usageHelp},
	}
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
// TODO: comment
func (this *Repl) handleNonDSLLine(trimmedLine string) bool {
	args := strings.Fields(trimmedLine)
	if len(args) == 0 {
		return false
	}
	verbName := args[0]

	// We handle all single lines starting with a colon.  Anything that starts
	// with a semicolon but which we don't recognize, we should say so here --
	// rather than deferring to the DSL parser which will say "cannot parse DSL
	// expression".
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
	fmt.Println(":load {one or more filenames containing Miller DSL statements}")
	fmt.Print(
		`'begin {...}' / 'end{...}' blocks are parsed and saved. Type ':begin' or
':end', respectively, to execute them. User-defined functions and subroutines
('func' and 'subr') are parsed and saved. Other statements are saved in a
'main' block.  Type ':main' to execute them on a given record. (See :open and
:read for more on how to do this.)
`)
}

func handleLoad(this *Repl, args []string) bool {
	args = args[1:] // strip off verb
	if len(args) < 1 {
		return false
	}
	for _, filename := range args {
		dslBytes, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Printf("Cannot load DSL expression from file \"%s\": ",
				filename)
			fmt.Println(err)
			return true
		}
		dslString := string(dslBytes)

		err = this.handleDSLStringBulk(dslString)
		if err != nil {
			fmt.Println(err)
		}
	}
	return true
}

// ----------------------------------------------------------------
func usageOpen(this *Repl) {
	fmt.Printf(
		":open {one or more data-file names, in the format specifed by %s %s\n",
		this.exeName, this.replName,
	)
	fmt.Print(
		`Then you can type :read to load the next record. Then any interactive
DSL commands will use that record. Also you can type ':main' to invoke any
main-block statements from multiline input or :load.

If zero data-file names are supplied, then standard input will be read when
you type :read.
`)

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
	fmt.Println(":read with no arguments.")
	fmt.Printf(
		"Reads in the next record from file(s) listed by :open, or by %s %s.\n",
		this.exeName, this.replName,
	)
	fmt.Println("Then you can operate on it with interactive statements, or :main, and you can")
	fmt.Println("write it out using :write.")
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
	fmt.Println(":write with no arguments.")
	fmt.Println("Sends the current record (maybe modifed by statements you enter)")
	fmt.Printf("to the output with format as specified by %s %s.\n",
		this.exeName, this.replName)
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
	fmt.Println(":begin with no arguments.")
	fmt.Println("Executes any begin {...} which have been entered.")
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
	fmt.Println(":main with no arguments.")
	fmt.Println("Executes any statements outside of begin/end/func/subr which have been.")
	fmt.Println("with :load or multi-line input.")
}
func handleMain(this *Repl, args []string) bool {
	_, err := this.cstRootNode.ExecuteMainBlock(this.runtimeState)
	if err != nil {
		fmt.Println(err)
	}
	return true
}

// ----------------------------------------------------------------
func usageEnd(this *Repl) {
	fmt.Println(":end with no arguments.")
	fmt.Println("Executes any end {...} which have been entered.")
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
	fmt.Println(":astprint {format option}")
	fmt.Println("Shows the AST (abstract syntax tree) associated with DSL statements entered in.")
	fmt.Println("Format options:")
	fmt.Println("parex  Prints AST as a parenthesized multi-line expression.")
	fmt.Println("parex1 Prints AST as a parenthesized single-line expression.")
	fmt.Println("indent Prints AST as an indented tree expression.")
	fmt.Println("none   Disables AST printing.")
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
	fmt.Println(":blocks with no arguments.")
	fmt.Println("Shows the number of begin{...} blocks that have been loaded, the number")
	fmt.Println("of main-block statements that have been loaded with :load or multi-line input,")
	fmt.Println("and the number of end{...} blocks that have been loaded.")

}
func handleBlocks(this *Repl, args []string) bool {
	this.cstRootNode.ShowBlockReport()
	return true
}

// ----------------------------------------------------------------
func usageHelp(this *Repl) {
	fmt.Println("TODO: metahelp is TBD.")
}

// PLAN:
// * :help
// * :help invocation (CLI flags ...)
//   mlr repl -h ...
//   -f/-e/-s ...
// * :help functions
// * :help keywords
//   --> sort them ...
// * :help {function}
// * :help {keyword}
// * :help repl
// * :help repl commands
// * :help repl intro
// * :help :foo

func handleHelp(this *Repl, args []string) bool {
	args = args[1:] // Strip off verb ':help'
	if len(args) == 0 {
		fmt.Println("Options:")
		fmt.Println(":help intro")
		fmt.Println(":help repl")
		fmt.Println(":help repl-details")
		fmt.Println(":help functions")
		fmt.Println(":help {function name}, e.g. :help sec2gmt")
		fmt.Println(":help {function name}, e.g. :help sec2gmt")
		return true
	}

	for _, arg := range args {
		handleHelpSingle(this, arg)
	}

	return true
}

func handleHelpSingle(this *Repl, arg string) {
	if arg == "intro" {
		showREPLIntro()
		return
	}

	if arg == "repl" {
		for _, entry := range handlerLookupTable {
			names := strings.Join(entry.verbNames, " or ")
			fmt.Println(names)
		}
		return
	}

	if arg == "repl-details" {
		for i, entry := range handlerLookupTable {
			if i > 0 {
				fmt.Println()
			}
			fmt.Printf("%s: \n", strings.Join(entry.verbNames, " or "))
			entry.usageFunc(this)
		}
		return
	}

	if arg == "functions" {
		cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionUsages(os.Stdout)
		return
	}

	if cst.BuiltinFunctionManagerInstance.TryListBuiltinFunctionUsage(arg, os.Stdout) {
		return
	}

	nonDSLHandler := this.findHandler(arg)
	if nonDSLHandler != nil {
		nonDSLHandler.usageFunc(this)
		return
	}

	fmt.Printf("No help available for %s\n", arg)
}

func showREPLIntro() {
	fmt.Print(
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
