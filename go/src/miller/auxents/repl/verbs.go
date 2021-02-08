// ================================================================
// Handlers for non-DSL statements like ':open foo.dat' or ':help'.
// ================================================================

package repl

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"miller/dsl/cst"
	"miller/lib"
	"miller/types"
)

// ----------------------------------------------------------------
// Types for the lookup table
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
		{verbNames: []string{":c", ":context"}, handlerFunc: handleContext, usageFunc: usageContext},
		{verbNames: []string{":s", ":skip"}, handlerFunc: handleSkip, usageFunc: usageSkip},
		{verbNames: []string{":p", ":process"}, handlerFunc: handleProcess, usageFunc: usageProcess},
		{verbNames: []string{":w", ":write"}, handlerFunc: handleWrite, usageFunc: usageWrite},
		{verbNames: []string{":b", ":begin"}, handlerFunc: handleBegin, usageFunc: usageBegin},
		{verbNames: []string{":m", ":main"}, handlerFunc: handleMain, usageFunc: usageMain},
		{verbNames: []string{":e", ":end"}, handlerFunc: handleEnd, usageFunc: usageEnd},
		{verbNames: []string{":astprint"}, handlerFunc: handleASTPrint, usageFunc: usageASTPrint},
		{verbNames: []string{":blocks"}, handlerFunc: handleBlocks, usageFunc: usageBlocks},
		{verbNames: []string{":q", ":quit"}, handlerFunc: nil, usageFunc: usageQuit},
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
// Handles a single non-DSL statement like ':open foo.dat' or ':help'.
func (this *Repl) handleNonDSLLine(trimmedLine string) bool {
	args := strings.Fields(trimmedLine)
	if len(args) == 0 {
		return false
	}
	verbName := args[0]

	// We handle all single lines starting with a colon.  Anything that starts
	// with a semicolon but which we don't recognize, we should say so here --
	// rather than deferring to the DSL parser which would only say "cannot
	// parse DSL expression", which would only be more confusing.
	if !strings.HasPrefix(verbName, ":") {
		return false
	}
	handler := this.findHandler(verbName)
	if handler == nil {
		fmt.Printf("REPL verb %s not found.\n", verbName)
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
		`Any 'begin {...}' / 'end{...}' blocks are parsed and saved. (You can then type
':begin' or ':end', respectively, to execute them.) User-defined functions and
subroutines ('func' and 'subr') are parsed and saved. Other statements are
saved in a 'main' block.  (You can then type ':main' to execute them on any
given record. See :open and :read for more on how to do this.)
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
		":open {one or more data-file names in the format specifed by %s %s}.\n",
		this.exeName, this.replName,
	)
	fmt.Print(
		`Then you can type :read to load the next record. Then any interactive
DSL commands will use that record. Also you can type ':main' to invoke any
main-block statements from multiline input or :load.

If zero data-file names are supplied (i.e. ':open' with no file names), then
each record will be taken from standard input when you type :read.
`)

}

func handleOpen(this *Repl, args []string) bool {
	args = args[1:] // strip off verb
	if openFilesPreCheck(this, args) {
		this.openFiles(args)
	}
	return true
}

// Using the record-reader API, if filenames are presented one or more of which
// are not accessible, then the 'no such file' error isn't encountered until
// the first record-read is attempted. For non-REPL uxe, this is fine. For REPL
// use, if the user types ':open nonesuch' then we want to proactively say
// something instead of waiting to show them an error only when they type
// ':read'.
func openFilesPreCheck(this *Repl, args []string) bool {
	if len(args) == 0 {
		// Zero file names is stdin, which is readable
	}
	for _, arg := range args {
		fileInfo, err := os.Stat(arg)
		if err != nil {
			fmt.Printf("%s %s: could not open \"%s\"\n",
				this.exeName, this.replName, arg,
			)
			return false
		}
		if fileInfo.IsDir() {
			fmt.Printf("%s %s: \"%s\" is a directory.\n",
				this.exeName, this.replName, arg,
			)
			return false
		}
	}
	return true
}

// Also invoked from the main entry-point, hence split out as a separate method.
func (this *Repl) openFiles(filenames []string) {
	this.inputChannel = make(chan *types.RecordAndContext, 10)
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
	args = args[1:] // strip off verb
	if len(args) != 0 {
		return false
	}
	if this.inputChannel == nil {
		fmt.Println("No open files")
		return true
	}

	var recordAndContext *types.RecordAndContext = nil
	var err error = nil

	select {
	case recordAndContext = <-this.inputChannel:
		break
	case err = <-this.errorChannel:
		break
	}

	if err != nil {
		fmt.Println(err)
		this.inputChannel = nil
		this.errorChannel = nil
		return true
	}

	if recordAndContext != nil {
		skipOrProcessRecord(
			this,
			recordAndContext,
			false, // processingNotSkipping -- since we will let the user interact with this record
			false, // testingByFilterExpression -- since we're just stepping by 1
		)
	}

	return true
}

// ----------------------------------------------------------------
func usageContext(this *Repl) {
	fmt.Println(":context with no arguments.")
	fmt.Println("Displays the current context variables: NR, FNR, FILENUM, FILENAME.")
}

func handleContext(this *Repl, args []string) bool {
	args = args[1:] // strip off verb
	if len(args) != 0 {
		return false
	}
	fmt.Println(this.runtimeState.Context.GetStatusString())
	return true
}

// ----------------------------------------------------------------
func usageSkip(this *Repl) {
	fmt.Println(":skip {n} to read n records without invoking :main statements or printing the records.")
	fmt.Printf(
		"Reads in the next record from file(s) listed by :open, or by %s %s.\n",
		this.exeName, this.replName,
	)
	fmt.Println("Then you can operate on it with interactive statements, or :main, and you can")
	fmt.Println("write it out using :write.")
	fmt.Println("Or: :skip until {some DSL expression}. You can use 'u' as shorthand for 'until'.")
	fmt.Println("Example: :skip until NR == 30")
	fmt.Println("Example: :skip until $status_code != 200")
}

func handleSkip(this *Repl, args []string) bool {
	if this.inputChannel == nil {
		fmt.Println("No open files")
		return true
	}

	args = args[1:] // strip off verb
	if len(args) < 1 {
		return false
	}

	if len(args) == 1 {
		n, ok := lib.TryIntFromString(args[0])
		if !ok {
			fmt.Printf("Could not parse \"%s\" as integer.\n", args[0])
		} else {
			handleSkipOrProcessN(this, n, false)
		}
		return true
	} else if args[0] != "until" && args[0] != "u" {
		return false
	} else {
		dslString := strings.Join(args[1:], " ")
		handleSkipOrProcessUntil(this, dslString, false)
		return true
	}
}

// ----------------------------------------------------------------
func usageProcess(this *Repl) {
	fmt.Println(":process {n} to read n records, invoking :main statements on them, and printing the records.")
	fmt.Printf(
		"Reads in the next record from file(s) listed by :open, or by %s %s.\n",
		this.exeName, this.replName,
	)
	fmt.Println("Then you can operate on it with interactive statements, or :main, and you can")
	fmt.Println("write it out using :write.")
	fmt.Println("Or: :process until {some DSL expression}. You can use 'u' as shorthand for 'until'.")
	fmt.Println("Example: :process until NR == 30")
	fmt.Println("Example: :process until $status_code != 200")
}

func handleProcess(this *Repl, args []string) bool {
	if this.inputChannel == nil {
		fmt.Println("No open files")
		return true
	}

	args = args[1:] // strip off verb
	if len(args) < 1 {
		return false
	}

	if len(args) == 1 {
		n, ok := lib.TryIntFromString(args[0])
		if !ok {
			fmt.Printf("Could not parse \"%s\" as integer.\n", args[0])
		} else {
			handleSkipOrProcessN(this, n, true)
		}
		return true
	} else if args[0] != "until" && args[0] != "u" {
		return false
	} else {
		dslString := strings.Join(args[1:], " ")
		handleSkipOrProcessUntil(this, dslString, true)
		return true
	}
}

// ----------------------------------------------------------------
func handleSkipOrProcessN(this *Repl, n int, processingNotSkipping bool) {
	var recordAndContext *types.RecordAndContext = nil
	var err error = nil

	for i := 1; i <= n; i++ {
		select {
		case recordAndContext = <-this.inputChannel:
			break
		case err = <-this.errorChannel:
			break
		}

		if err != nil {
			fmt.Println(err)
			this.inputChannel = nil
			this.errorChannel = nil
			return
		}

		if recordAndContext != nil {
			shouldBreak := skipOrProcessRecord(
				this,
				recordAndContext,
				processingNotSkipping,
				false, // testingByFilterExpression -- since we're counting to N
			)
			if shouldBreak {
				break
			}
		}
	}
}

func handleSkipOrProcessUntil(this *Repl, dslString string, processingNotSkipping bool) {
	astRootNode, err := this.BuildASTFromStringWithMessage(dslString)
	if err != nil {
		// Error message already printed out
		return
	}

	err = this.cstRootNode.IngestAST(
		astRootNode,
		false, /*isFilter*/
		true,  /*isReplImmediate*/
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = this.cstRootNode.Resolve()
	if err != nil {
		fmt.Println(err)
		return
	}

	var recordAndContext *types.RecordAndContext = nil

	for {
		select {
		case recordAndContext = <-this.inputChannel:
			break
		case err = <-this.errorChannel:
			break
		}

		if err != nil {
			fmt.Println(err)
			this.inputChannel = nil
			this.errorChannel = nil
			return
		}

		if recordAndContext != nil {
			shouldBreak := skipOrProcessRecord(
				this,
				recordAndContext,
				processingNotSkipping,
				true, // testingByFilterExpression -- since we're continuing until the filter expresssion is true
			)
			if shouldBreak {
				break
			}
		}
	}
}

// Three things can come through:
//
// * End-of-stream marker
// * Non-nil record to be printed
// * Strings to be printed from put/filter DSL print/dump/etc statements. They
//   are handled here rather than fmt.Println directly in the put/filter
//   handlers since we want all print statements and record-output to be in the
//   same goroutine, for deterministic output ordering.
//
// The first two are passed to the transformer. The third we send along the
// output channel without involving the record-transformer, since there is no
// record to be transformed.

// Return value is true if an end-of-loop condition has been detected.
func skipOrProcessRecord(
	this *Repl,
	recordAndContext *types.RecordAndContext,
	processingNotSkipping bool, // TODO: make this an enum
	testingByFilterExpression bool, // TODO: make this an enum
) bool { // TODO: make this an enum

	// Acquire incremented NR/FNR/FILENAME/etc
	this.runtimeState.Update(recordAndContext.Record, &recordAndContext.Context)

	// End-of-stream marker
	if recordAndContext.EndOfStream == true {
		fmt.Println("End of record stream")
		this.inputChannel = nil
		this.errorChannel = nil
		return true
	}

	// Strings to be printed from put/filter DSL print/dump/etc statements.
	if recordAndContext.Record == nil {
		if processingNotSkipping {
			fmt.Print(recordAndContext.OutputString)
		}
		return false
	}

	// Non-nil record to be printed
	if processingNotSkipping {
		outrec, err := this.cstRootNode.ExecuteMainBlock(this.runtimeState)
		if err != nil {
			fmt.Println(err)
			return true
		}
		this.runtimeState.Inrec = outrec
		this.recordWriter.Write(outrec, os.Stdout)
	}

	if testingByFilterExpression {
		_, err := this.cstRootNode.ExecuteREPLImmediate(this.runtimeState)
		if err != nil {
			fmt.Println(err)
			return true
		}

		filterBool, isBool := this.runtimeState.FilterExpression.GetBoolValue()
		if !isBool {
			filterBool = false
		}
		if filterBool {
			return true
		}
	}
	return false
}

// ----------------------------------------------------------------
func usageWrite(this *Repl) {
	fmt.Println(":write with no arguments.")
	fmt.Println("Sends the current record (maybe modifed by statements you enter)")
	fmt.Printf("to standard output, with format as specified by %s %s.\n",
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
	fmt.Println("Executes any begin {...} blocks which have been entered.")
}
func handleBegin(this *Repl, args []string) bool {
	args = args[1:] // strip off verb
	if len(args) != 0 {
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
	fmt.Println("Executes any statements outside of begin/end/func/subr which have been entered")
	fmt.Println("with :load or multi-line input.")
}
func handleMain(this *Repl, args []string) bool {
	args = args[1:] // strip off verb
	if len(args) != 0 {
		return false
	}
	_, err := this.cstRootNode.ExecuteMainBlock(this.runtimeState)
	if err != nil {
		fmt.Println(err)
	}
	return true
}

// ----------------------------------------------------------------
func usageEnd(this *Repl) {
	fmt.Println(":end with no arguments.")
	fmt.Println("Executes any end {...} blocks which have been entered.")
}
func handleEnd(this *Repl, args []string) bool {
	args = args[1:] // strip off verb
	if len(args) != 0 {
		return false
	}
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
	fmt.Println("none   Disables AST printing. (This is the default.)")
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
	args = args[1:] // strip off verb
	if len(args) != 0 {
		return false
	}
	this.cstRootNode.ShowBlockReport()
	return true
}

// ----------------------------------------------------------------
func usageQuit(this *Repl) {
	fmt.Println(":quit with no arguments.")
	fmt.Println("Ends the Miller REPL session.")
}

// The :quit command is handled outside this file; we have a help function,
// though, to expose it for on-line help.

// ----------------------------------------------------------------
func usageHelp(this *Repl) {
	fmt.Println("Options:")
	fmt.Println(":help intro")
	fmt.Println(":help repl-list")
	fmt.Println(":help repl-details")
	fmt.Println(":help prompt")
	fmt.Println(":help function-names")
	fmt.Println(":help function-detailss")
	fmt.Println(":help {function name}, e.g. :help sec2gmt")
	fmt.Println(":help {function name}, e.g. :help sec2gmt")
}

func handleHelp(this *Repl, args []string) bool {
	args = args[1:] // Strip off verb ':help'
	if len(args) == 0 {
		usageHelp(this)
		return true
	}

	for _, arg := range args {
		handleHelpSingle(this, arg)
	}

	return true
}

func handleHelpSingle(this *Repl, arg string) {
	if arg == "intro" {
		showREPLIntro(this)
		return
	}

	if arg == "repl-commands" {
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

	if arg == "prompt" {
		fmt.Printf(
			"You can export the environment variable %s to customize the Miller REPL prompt.\n",
			ENV_PRIMARY_PROMPT,
		)
		fmt.Printf(
			"Otherwise, it defaults to \"%s\".\n",
			DEFAULT_PRIMARY_PROMPT,
		)
		fmt.Printf(
			"Likewise you can export the environment variable %s to customize the secondary prompt,\n",
			ENV_SECONDARY_PROMPT,
		)
		fmt.Printf(
			"which defaults to \"%s\". This is used for multiline input.\n",
			DEFAULT_SECONDARY_PROMPT,
		)
		return
	}

	if arg == "function-names" {
		cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionsRaw(os.Stdout)
		return
	}

	if arg == "function-details" {
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

func showREPLIntro(this *Repl) {
	fmt.Printf(
`The Miller REPL is an interactive counterpart to record-processing using the
put/filter DSL.

Using put and filter, you can do the following:
* Specify input format (e.g. --icsv), output format (e.g. --ojson), etc. using
  command-line flags.
* Specify filenames on the command line.
* Define begin {...} blocks which are executed before the first record is read.
* Define end {...} blocks which are executed after the last record is read.
* Define user-defined functions/subroutines using func/subr.
* Specify statements to be executed on each record.
* Example:
  %s --icsv --ojson put 'begin {print "HELLO"} $z = $x + $y end {print "GOODBYE"}

Using the REPL, by contrast, you get interactive control over those steps:
* Specify input format (e.g. --icsv), output format (e.g. --ojson), etc. using
  command-line flags.
* Specify filenames either on the command line or via ':open' at the Miller
  REPL.
* Read records one at a time using ':read'.
* Skip ahead using statements ':skip 10' or ':skip until NR == 100' or ':skip
  until $status_code != 200'.
* Similarly, but processing records rather than skipping past them, using
  ':process' rather than ':skip'.
* Skip ahead using statements ':skip 10' or ':skip until NR == 100' or ':skip
  until $status_code != 200'.
* Define begin {...} blocks; invoke them at will using ':begin'.
* Define end {...} blocks; invoke them at will using ':end'.
* Define user-defined functions/subroutines using func/subr; call them from other statements.
* Interactively specify statements to be executed on the current record.
* Load any of the above from Miller-script files using ':load'.
* Furthermore, any DSL statements other than begin/end/func/subr loaded using
  ':load' -- or from "multiline input mode" where you type '<' on a line by
  itself, enter the code, then type '>' on a line by itself -- will be
  remembered and can be invoked on a given record using ':main'.  In multi-line
  mode, semicolons are required between statements; otherwise they are not
  needed.

At this REPL prompt you can enter any Miller DSL expression.  REPL-only
statements (non-DSL statements) start with ':', such as ':help' or ':quit'.
Type ':help' to see more about your options.

No command-line-history-editing feature is built in but 'rlwrap mlr repl' is a
delight. You may need 'brew install rlwrap', 'sudo apt-get install rlwrap',
etc. depending on your platform.

The input "record" by default is the empty map but you can do things like
'$x=3', or 'unset $y', or '$* = {"x": 3, "y": 4}' to populate it. Or, ':open
foo.dat' and ':read' to populat it from a data file.

Non-assignment expressions, such as '7' or 'true', operate as filter conditions
in the put DSL: they can be used to specify whether a record will or won't be
included in the output-record stream.  But here in the REPL, they are simply
printed to the terminal, e.g. if you type '1+2', you will see '3'.

Examples:

[mlr] 1+2
3

[mlr] x=3  # These are local variables
[mlr] y=4
[mlr] x+y
7

[mlr] <
func f(a,b) {
  return a**b
}
>
[mlr] f(7,5)
16807

[mlr] :open foo.dat
[mlr] :read
[mlr] :context
FILENAME="foo.dat",FILENUM=1,NR=1,FNR=1
[mlr] $*
{
  "a": "eks",
  "b": "wye",
  "i": 4,
  "x": 0.38139939387114097,
  "y": 0.13418874328430463
}
[mlr] f($x,$i)
0.021160211005187134
`, this.exeName)
}
