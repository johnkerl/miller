// ================================================================
// Handlers for non-DSL statements like ':open foo.dat' or ':help'.
// ================================================================

package repl

import (
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/colorizer"
	"mlr/internal/pkg/dsl"
	"mlr/internal/pkg/dsl/cst"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
// Types for the lookup table.

// Handlers should return false if they want their usage function to be called.
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
		{verbNames: []string{":reopen"}, handlerFunc: handleReopen, usageFunc: usageReopen},
		{verbNames: []string{":r", ":read"}, handlerFunc: handleRead, usageFunc: usageRead},
		{verbNames: []string{":w", ":write"}, handlerFunc: handleWrite, usageFunc: usageWrite},
		{verbNames: []string{":rw"}, handlerFunc: handleReadWrite, usageFunc: usageReadWrite},
		{verbNames: []string{":c", ":context"}, handlerFunc: handleContext, usageFunc: usageContext},
		{verbNames: []string{":s", ":skip"}, handlerFunc: handleSkip, usageFunc: usageSkip},
		{verbNames: []string{":p", ":process"}, handlerFunc: handleProcess, usageFunc: usageProcess},
		{verbNames: []string{":w", ":>"}, handlerFunc: handleRedirectWrite, usageFunc: usageRedirectWrite},
		{verbNames: []string{":w", ":>>"}, handlerFunc: handleRedirectAppend, usageFunc: usageRedirectAppend},
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
func (repl *Repl) findHandler(verbName string) *handlerInfo {
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
func (repl *Repl) handleNonDSLLine(trimmedLine string) bool {
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
	handler := repl.findHandler(verbName)
	if handler == nil {
		fmt.Printf("REPL verb %s not found.\n", verbName)
		return true
	}

	if !handler.handlerFunc(repl, args) {
		handler.usageFunc(repl)
	}
	return true
}

// ----------------------------------------------------------------
func usageLoad(repl *Repl) {
	fmt.Println(":load {one or more filenames containing Miller DSL statements}")
	fmt.Println("If a filename is a directory, all \"*.mlr\" files will be loaded from within it.")
	fmt.Print(
		`Any 'begin {...}' / 'end{...}' blocks are parsed and saved. (You can then type
':begin' or ':end', respectively, to execute them.) User-defined functions and
subroutines ('func' and 'subr') are parsed and saved. Other statements are
saved in a 'main' block.  (You can then type ':main' to execute them on any
given record. See :open and :read for more on how to do repl.)
`)
}

func handleLoad(repl *Repl, args []string) bool {
	args = args[1:] // strip off verb
	if len(args) < 1 {
		return false
	}
	for _, filename := range args {
		dslStrings, err := lib.LoadStringsFromFileOrDir(filename, ".mlr")
		if err != nil {
			fmt.Printf("Cannot load DSL expression file \"%s\": ",
				filename)
			fmt.Println(err)
			return true
		}

		for _, dslString := range dslStrings {
			err = repl.handleDSLStringBulk(dslString, repl.doWarnings)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return true
}

// ----------------------------------------------------------------
func usageOpen(repl *Repl) {
	fmt.Printf(
		":open {one or more data-file names in the format specifed by %s %s}.\n",
		repl.exeName, repl.replName,
	)
	fmt.Print(
		`Then you can type :read to load the next record. Then any interactive
DSL commands will use that record. Also you can type ':main' to invoke any
main-block statements from multi-line input or :load.

If zero data-file names are supplied (i.e. ':open' with no file names), then
each record will be taken from standard input when you type :read.
`)

}

func handleOpen(repl *Repl, args []string) bool {
	args = args[1:] // strip off verb
	if openFilesPreCheck(repl, args) {
		repl.openFiles(args)
	}
	return true
}

// Using the record-reader API, if filenames are presented one or more of which
// are not accessible, then the 'no such file' error isn't encountered until
// the first record-read is attempted. For non-REPL uxe, this is fine. For REPL
// use, if the user types ':open nonesuch' then we want to proactively say
// something instead of waiting to show them an error only when they type
// ':read'.
func openFilesPreCheck(repl *Repl, args []string) bool {
	if len(args) == 0 {
		// Zero file names is stdin, which is readable
	}
	for _, arg := range args {
		fileInfo, err := os.Stat(arg)
		if err != nil {
			fmt.Printf("%s %s: could not open \"%s\"\n",
				repl.exeName, repl.replName, arg,
			)
			return false
		}
		if fileInfo.IsDir() {
			fmt.Printf("%s %s: \"%s\" is a directory.\n",
				repl.exeName, repl.replName, arg,
			)
			return false
		}
	}
	return true
}

// Also invoked from the main entry-point, hence split out as a separate method.
func (repl *Repl) openFiles(filenames []string) {
	// Remember for :reopen
	repl.options.FileNames = filenames

	repl.inputChannel = make(chan *types.RecordAndContext, 10)
	repl.errorChannel = make(chan error, 1)
	repl.downstreamDoneChannel = make(chan bool, 1)

	go repl.recordReader.Read(
		filenames,
		*repl.runtimeState.Context,
		repl.inputChannel,
		repl.errorChannel,
		repl.downstreamDoneChannel,
	)
}

// ----------------------------------------------------------------
func usageReopen(repl *Repl) {
	fmt.Println(":reopen with no arguments.")
	fmt.Println("Like :open with the same filenames you provided at the time you typed :open.")
}

func handleReopen(repl *Repl, args []string) bool {
	args = args[1:] // strip off verb
	if len(args) != 0 {
		return false
	}

	if openFilesPreCheck(repl, repl.options.FileNames) {
		repl.openFiles(repl.options.FileNames)
	}
	return true
}

// ----------------------------------------------------------------
func usageRead(repl *Repl) {
	fmt.Println(":read with no arguments.")
	fmt.Printf(
		"Reads in the next record from file(s) listed by :open, or by %s %s.\n",
		repl.exeName, repl.replName,
	)
	fmt.Println("Then you can operate on it with interactive statements, or :main, and you can")
	fmt.Println("write it out using :write.")
}

func handleRead(repl *Repl, args []string) bool {
	args = args[1:] // strip off verb
	if len(args) != 0 {
		return false
	}
	if repl.inputChannel == nil {
		fmt.Println("No open files")
		return true
	}

	var recordAndContext *types.RecordAndContext = nil
	var err error = nil

	select {
	case recordAndContext = <-repl.inputChannel:
		break
	case err = <-repl.errorChannel:
		break
	}

	if err != nil {
		fmt.Println(err)
		repl.inputChannel = nil
		repl.errorChannel = nil
		return true
	}

	if recordAndContext != nil {
		skipOrProcessRecord(
			repl,
			recordAndContext,
			false, // processingNotSkipping -- since we will let the user interact with this record
			false, // testingByFilterExpression -- since we're just stepping by 1
		)
	}

	return true
}

// ----------------------------------------------------------------
func usageContext(repl *Repl) {
	fmt.Println(":context with no arguments.")
	fmt.Println("Displays the current context variables: NR, FNR, FILENUM, FILENAME.")
}

func handleContext(repl *Repl, args []string) bool {
	args = args[1:] // strip off verb
	if len(args) != 0 {
		return false
	}
	fmt.Println(repl.runtimeState.Context.GetStatusString())
	return true
}

// ----------------------------------------------------------------
func usageSkip(repl *Repl) {
	fmt.Println(":skip {n} to read n records without invoking :main statements or printing the records.")
	fmt.Printf(
		"Reads in the next record from file(s) listed by :open, or by %s %s.\n",
		repl.exeName, repl.replName,
	)
	fmt.Println("Then you can operate on it with interactive statements, or :main, and you can")
	fmt.Println("write it out using :write.")
	fmt.Println("Or: :skip until {some DSL expression}. You can use 'u' as shorthand for 'until'.")
	fmt.Println("Example: :skip until NR == 30")
	fmt.Println("Example: :skip until $status_code != 200")
	fmt.Println("Or: ':skip until intr' which means keep skipping until you type control-C to interrupt.")
}

func handleSkip(repl *Repl, args []string) bool {
	if repl.inputChannel == nil {
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
			handleSkipOrProcessN(repl, n, false)
		}
		return true
	} else if args[0] != "until" && args[0] != "u" {
		return false
	} else {
		args := args[1:]
		dslString := strings.Join(args, " ")
		// If they say ':skip until intr' then we use a DSL string of 'false',
		// i.e. skip until they type control-C.
		if len(args) == 1 && args[0] == "intr" {
			dslString = "false"
		}
		handleSkipOrProcessUntil(repl, dslString, false)
		return true
	}
}

// ----------------------------------------------------------------
func usageProcess(repl *Repl) {
	fmt.Println(":process {n} to read n records, invoking :main statements on them, and printing the records.")
	fmt.Printf(
		"Reads in the next record from file(s) listed by :open, or by %s %s.\n",
		repl.exeName, repl.replName,
	)
	fmt.Println("Then you can operate on it with interactive statements, or :main, and you can")
	fmt.Println("write it out using :write.")
	fmt.Println("Or: :process until {some DSL expression}. You can use 'u' as shorthand for 'until'.")
	fmt.Println("Example: :process until NR == 30")
	fmt.Println("Example: :process until $status_code != 200")
	fmt.Println("Or: ':process until intr' which means keep processing until you type control-C to interrupt.")
}

func handleProcess(repl *Repl, args []string) bool {
	if repl.inputChannel == nil {
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
			handleSkipOrProcessN(repl, n, true)
		}
		return true
	} else if args[0] != "until" && args[0] != "u" {
		return false
	} else {
		args := args[1:]
		dslString := strings.Join(args, " ")
		// If they say ':process until intr' then we use a DSL string of 'false',
		// i.e. skip until they type control-C.
		if len(args) == 1 && args[0] == "intr" {
			dslString = "false"
		}
		handleSkipOrProcessUntil(repl, dslString, true)
		return true
	}
}

// ----------------------------------------------------------------
func handleSkipOrProcessN(repl *Repl, n int, processingNotSkipping bool) {
	var recordAndContext *types.RecordAndContext = nil
	var err error = nil

	for i := 1; i <= n; i++ {
		select {
		case recordAndContext = <-repl.inputChannel:
			break
		case err = <-repl.errorChannel:
			break
		case _ = <-repl.appSignalNotificationChannel: // user typed control-C
			break
		}

		if err != nil {
			fmt.Println(err)
			repl.inputChannel = nil
			repl.errorChannel = nil
			return
		}

		if recordAndContext != nil {
			shouldBreak := skipOrProcessRecord(
				repl,
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

func handleSkipOrProcessUntil(repl *Repl, dslString string, processingNotSkipping bool) {

	err := repl.cstRootNode.Build(
		[]string{dslString},
		cst.DSLInstanceTypeREPL,
		true, // isReplImmediate
		repl.doWarnings,
		false, // warningsAreFatal
		func(dslString string, astNode *dsl.AST) {
			if repl.astPrintMode == ASTPrintParex {
				astNode.PrintParex()
			} else if repl.astPrintMode == ASTPrintParexOneLine {
				astNode.PrintParexOneLine()
			} else if repl.astPrintMode == ASTPrintIndent {
				astNode.Print()
			}
		},
	)
	if err != nil {
		// Error message already printed out
		//TODO: check this
		return
	}

	var recordAndContext *types.RecordAndContext = nil

	for {
		doubleBreak := false
		select {
		case recordAndContext = <-repl.inputChannel:
			break
		case err = <-repl.errorChannel:
			break
		case _ = <-repl.appSignalNotificationChannel: // user typed control-C
			doubleBreak = true
			break
		}
		if doubleBreak {
			break
		}

		if err != nil {
			fmt.Println(err)
			repl.inputChannel = nil
			repl.errorChannel = nil
			return
		}

		if recordAndContext != nil {
			shouldBreak := skipOrProcessRecord(
				repl,
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
	repl *Repl,
	recordAndContext *types.RecordAndContext,
	processingNotSkipping bool, // TODO: make this an enum
	testingByFilterExpression bool, // TODO: make this an enum
) bool { // TODO: make this an enum

	// Acquire incremented NR/FNR/FILENAME/etc
	repl.runtimeState.Update(recordAndContext.Record, &recordAndContext.Context)

	// End-of-stream marker
	if recordAndContext.EndOfStream == true {
		fmt.Println("End of record stream")
		repl.inputChannel = nil
		repl.errorChannel = nil
		return true
	}

	// Strings to be printed from put/filter DSL print/dump/etc statements.
	if recordAndContext.Record == nil {
		if processingNotSkipping {
			fmt.Fprint(repl.outputStream, recordAndContext.OutputString)
		}
		return false
	}

	// Non-nil record to be printed
	if processingNotSkipping {
		outrec, err := repl.cstRootNode.ExecuteMainBlock(repl.runtimeState)
		if err != nil {
			fmt.Println(err)
			return true
		}
		repl.runtimeState.Inrec = outrec
		writeRecord(repl, repl.runtimeState.Inrec)
	}

	if testingByFilterExpression {
		_, err := repl.cstRootNode.ExecuteREPLImmediate(repl.runtimeState)
		if err != nil {
			fmt.Println(err)
			return true
		}

		filterBool, isBool := repl.runtimeState.FilterExpression.GetBoolValue()
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
func usageWrite(repl *Repl) {
	fmt.Println(":write with no arguments.")
	fmt.Println("Sends the current record (maybe modifed by statements you enter)")
	fmt.Printf("to standard output, with format as specified by %s %s.\n",
		repl.exeName, repl.replName)
}
func handleWrite(repl *Repl, args []string) bool {
	if len(args) != 1 {
		return false
	}
	writeRecord(repl, repl.runtimeState.Inrec)
	return true
}

func writeRecord(repl *Repl, outrec *types.Mlrmap) {
	if outrec != nil {
		// E.g. '{"req": {"method": "GET", "path": "/api/check"}}' becomes
		// req.method=GET,req.path=/api/check.
		if repl.options.WriterOptions.AutoFlatten {
			outrec.Flatten(repl.options.WriterOptions.FLATSEP)
		}
		// E.g.  req.method=GET,req.path=/api/check becomes
		// '{"req": {"method": "GET", "path": "/api/check"}}'
		if repl.options.WriterOptions.AutoUnflatten {
			outrec.Unflatten(repl.options.WriterOptions.FLATSEP)
		}
	}
	repl.recordWriter.Write(outrec, repl.outputStream, true /*outputIsStdout*/)
}

// ----------------------------------------------------------------
func usageReadWrite(repl *Repl) {
	fmt.Println(":rw with no arguments.")
	fmt.Println("Same as ':r' followed by ':w'.")
}
func handleReadWrite(repl *Repl, args []string) bool {
	if !handleRead(repl, args) {
		return false
	}
	if !handleWrite(repl, args) {
		return false
	}
	return true
}

// ----------------------------------------------------------------
func usageRedirectWrite(repl *Repl) {
	fmt.Println(":> {filename} sends record-write output to the specified file.")
	fmt.Println(":> with no arguments sends record-write output to stdout.")
}
func handleRedirectWrite(repl *Repl, args []string) bool {
	args = args[1:] // strip off verb
	if len(args) == 0 {
		// TODO: fclose old if not already os.Stdout
		repl.outputStream = os.Stdout
		return true
	}

	if len(args) != 1 {
		return false
	}

	filename := args[0]
	handle, err := os.OpenFile(
		filename,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0644, // TODO: let users parameterize this
	)
	if err != nil {
		fmt.Printf(
			"%s %s: couldn't open \"%s\" for write.\n",
			repl.exeName, repl.replName, filename,
		)
	}
	fmt.Printf("Redirecting record output to \"%s\"\n", filename)

	// TODO: fclose old if not already os.Stdout
	repl.outputStream = handle

	return true
}

// ----------------------------------------------------------------
func usageRedirectAppend(repl *Repl) {
	fmt.Println(":>> {filename}")
	fmt.Println("Appends record-write output to the specified file.")
}
func handleRedirectAppend(repl *Repl, args []string) bool {
	args = args[1:] // strip off verb
	if len(args) != 1 {
		return false
	}

	filename := args[0]
	handle, err := os.OpenFile(
		filename,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644, // TODO: let users parameterize this
	)
	if err != nil {
		fmt.Printf(
			"%s %s: couldn't open \"%s\" for write.\n",
			repl.exeName, repl.replName, filename,
		)
	}
	fmt.Printf("Redirecting record output to \"%s\"\n", filename)

	// TODO: fclose old if not already os.Stdout
	repl.outputStream = handle

	return true
}

// ----------------------------------------------------------------
func usageBegin(repl *Repl) {
	fmt.Println(":begin with no arguments.")
	fmt.Println("Executes any begin {...} blocks which have been entered.")
}
func handleBegin(repl *Repl, args []string) bool {
	args = args[1:] // strip off verb
	if len(args) != 0 {
		return false
	}
	err := repl.cstRootNode.ExecuteBeginBlocks(repl.runtimeState)
	if err != nil {
		fmt.Println(err)
	}
	return true
}

// ----------------------------------------------------------------
func usageMain(repl *Repl) {
	fmt.Println(":main with no arguments.")
	fmt.Println("Executes any statements outside of begin/end/func/subr which have been entered")
	fmt.Println("with :load or multi-line input.")
}
func handleMain(repl *Repl, args []string) bool {
	args = args[1:] // strip off verb
	if len(args) != 0 {
		return false
	}
	_, err := repl.cstRootNode.ExecuteMainBlock(repl.runtimeState)
	if err != nil {
		fmt.Println(err)
	}
	return true
}

// ----------------------------------------------------------------
func usageEnd(repl *Repl) {
	fmt.Println(":end with no arguments.")
	fmt.Println("Executes any end {...} blocks which have been entered.")
}
func handleEnd(repl *Repl, args []string) bool {
	args = args[1:] // strip off verb
	if len(args) != 0 {
		return false
	}
	err := repl.cstRootNode.ExecuteEndBlocks(repl.runtimeState)
	if err != nil {
		fmt.Println(err)
	}
	return true
}

// ----------------------------------------------------------------
func usageASTPrint(repl *Repl) {
	fmt.Println(":astprint {format option}")
	fmt.Println("Shows the AST (abstract syntax tree) associated with DSL statements entered in.")
	fmt.Println("Format options:")
	fmt.Println("parex  Prints AST as a parenthesized multi-line expression.")
	fmt.Println("parex1 Prints AST as a parenthesized single-line expression.")
	fmt.Println("indent Prints AST as an indented tree expression.")
	fmt.Println("none   Disables AST printing. (This is the default.)")
}
func handleASTPrint(repl *Repl, args []string) bool {
	args = args[1:] // strip off verb
	if len(args) != 1 {
		return false
	}
	style := args[0]
	if style == "parex" {
		repl.astPrintMode = ASTPrintParex
	} else if style == "parex1" {
		repl.astPrintMode = ASTPrintParexOneLine
	} else if style == "indent" {
		repl.astPrintMode = ASTPrintIndent
	} else if style == "none" {
		repl.astPrintMode = ASTPrintNone
	} else {
		fmt.Printf("Unrecognized style %s: see ':help :astprint'.\n", style)
	}
	return true
}

// ----------------------------------------------------------------
func usageBlocks(repl *Repl) {
	fmt.Println(":blocks with no arguments.")
	fmt.Println("Shows the number of begin{...} blocks that have been loaded, the number")
	fmt.Println("of main-block statements that have been loaded with :load or non-immediate")
	fmt.Println("multi-line input, and the number of end{...} blocks that have been loaded.")

}
func handleBlocks(repl *Repl, args []string) bool {
	args = args[1:] // strip off verb
	if len(args) != 0 {
		return false
	}
	repl.cstRootNode.ShowBlockReport()
	return true
}

// ----------------------------------------------------------------
func usageQuit(repl *Repl) {
	fmt.Println(":quit with no arguments.")
	fmt.Println("Ends the Miller REPL session.")
}

// The :quit command is handled outside this file; we have a help function,
// though, to expose it for online help.

// ----------------------------------------------------------------
func usageHelp(repl *Repl) {
	fmt.Println("Options:")
	fmt.Println(":help intro")
	fmt.Println(":help examples")
	fmt.Println(":help repl-list")
	fmt.Println(":help repl-details")
	fmt.Println(":help prompt")
	fmt.Println(":help function-names")
	fmt.Println(":help function-details")
	fmt.Println(":help {function name}, e.g. :help sec2gmt")
	fmt.Println(":help {function name}, e.g. :help sec2gmt")
}

func handleHelp(repl *Repl, args []string) bool {
	args = args[1:] // Strip off verb ':help'
	if len(args) == 0 {
		usageHelp(repl)
		return true
	}

	for _, arg := range args {
		handleHelpSingle(repl, arg)
	}

	return true
}

func handleHelpSingle(repl *Repl, arg string) {
	if arg == "intro" {
		showREPLIntro(repl)
		return
	}

	if arg == "examples" {
		showREPLExamples(repl)
		return
	}

	if arg == "repl-list" {
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
			fmt.Println(colorizer.MaybeColorizeHelp(strings.Join(entry.verbNames, " or "), true))
			entry.usageFunc(repl)
		}
		return
	}

	if arg == "prompt" {
		fmt.Printf(
			"You can export the environment variable %s to customize the Miller REPL prompt.\n",
			colorizer.MaybeColorizeHelp(ENV_PRIMARY_PROMPT, true),
		)

		fmt.Printf(
			"Otherwise, it defaults to \"%s\".\n",
			colorizer.MaybeColorizeHelp(DEFAULT_PRIMARY_PROMPT, true),
		)

		fmt.Printf(
			"Likewise you can export the environment variable %s to customize the secondary prompt,\n",
			colorizer.MaybeColorizeHelp(ENV_SECONDARY_PROMPT, true),
		)

		fmt.Printf(
			"which defaults to \"%s\". This is used for multi-line input.\n",
			colorizer.MaybeColorizeHelp(DEFAULT_SECONDARY_PROMPT, true),
		)

		return
	}

	if arg == "function-names" {
		cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionNamesAsParagraph()
		return
	}

	if arg == "function-details" {
		cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionUsages()
		return
	}

	if cst.TryUsageForKeyword(arg) {
		return
	}

	if cst.BuiltinFunctionManagerInstance.TryListBuiltinFunctionUsage(arg, true) {
		return
	}

	nonDSLHandler := repl.findHandler(arg)
	if nonDSLHandler != nil {
		nonDSLHandler.usageFunc(repl)
		return
	}

	fmt.Printf("No help available for %s\n", arg)
}

func showREPLIntro(repl *Repl) {
	fmt.Println(colorizer.MaybeColorizeHelp("What the Miller REPL is", true))
	fmt.Println(
		`The Miller REPL (read-evaluate-print loop) is an interactive counterpart
to record-processing using the put/filter DSL (domain-specific language).`)
	fmt.Println()

	fmt.Println(colorizer.MaybeColorizeHelp("Using Miller without the REPL:", true))
	fmt.Printf(
		`Using put and filter, you can do the following:
* Specify input format (e.g. --icsv), output format (e.g. --ojson), etc. using
  command-line flags.
* Specify filenames on the command line.
* Define begin {...} blocks which are executed before the first record is read.
* Define end {...} blocks which are executed after the last record is read.
* Define user-defined functions/subroutines using func and subr.
* Specify statements to be executed on each record -- which are anything outside of begin/end/func/subr.
* Example:
  %s --icsv --ojson put 'begin {print "HELLO"} $z = $x + $y; end {print "GOODBYE"}`,
		repl.exeName)
	fmt.Println()
	fmt.Println()

	fmt.Println(colorizer.MaybeColorizeHelp("Using Miller with the REPL:", true))
	fmt.Println(
		`Using the REPL, by contrast, you get interactive control over those same steps:
* Specify input format (e.g. --icsv), output format (e.g. --ojson), etc. using
  command-line flags.
* REPL-only statements (non-DSL statements) start with ':', such as ':help' or ':quit'
  or ':open'.
* Specify filenames either on the command line or via ':open' at the Miller REPL.
* Read records one at a time using ':read'.
* Write the current record (maybe after you've modified it with things like '$z = $x + $y')
  using ':write'. This goes to the terminal; you can use ':> {filename}' to make writes
  go to a file, or ':>> {filename}' to append.
* You can type ':reopen' to go back to the start of the same file(s) you specified
  with ':open'.
* Skip ahead using statements ':skip 10' or ':skip until NR == 100' or
  ':skip until $status_code != 200'.
* Similarly, but processing records rather than skipping past them, using
  ':process' rather than ':skip'. Like ':write', these go to the screen;
  use ':> {filename}' or ':>> {filename}' to log to a file instead.
* Define begin {...} blocks; invoke them at will using ':begin'.
* Define end {...} blocks; invoke them at will using ':end'.
* Define user-defined functions/subroutines using func/subr; call them from other statements.
* Interactively specify statements to be executed immediately on the current record.
* Load any of the above from Miller-script files using ':load'.`)
	fmt.Println()

	fmt.Println(
		`The input "record" by default is the empty map but you can do things like
'$x=3', or 'unset $y', or '$* = {"x": 3, "y": 4}' to populate it. Or, ':open
foo.dat' followed by ':read' to populate it from a data file.

Non-assignment expressions, such as '7' or 'true', operate as filter conditions
in the put DSL: they can be used to specify whether a record will or won't be
included in the output-record stream.  But here in the REPL, they are simply
printed to the terminal, e.g. if you type '1+2', you will see '3'.`)
	fmt.Println()

	fmt.Println(colorizer.MaybeColorizeHelp("Entering multi-line statements", true))
	fmt.Println(
		`* To enter multi-line statements, enter '<' on a line by itself, then the code (taking care
  for semicolons), then ">" on a line by itself. These will be executed immediately.
* If you enter '<<' on a line by itself, then the code, then '>>' on a line by
  itself, the statements will be remembered for executing on records with
  ':main', as if you had done ':load' to load statements from a file.`)
	fmt.Println()

	fmt.Println(colorizer.MaybeColorizeHelp("History-editing:", true))
	fmt.Println(
		`No command-line-history-editing feature is built in but 'rlwrap mlr repl' is a
delight. You may need 'brew install rlwrap', 'sudo apt-get install rlwrap',
etc. depending on your platform.`)
	fmt.Println()

	fmt.Println(colorizer.MaybeColorizeHelp("Online help:", true))
	fmt.Println("Type ':h' or ':help' to see more about your options. In particular, ':help examples'.")
}

// ----------------------------------------------------------------
func showREPLExamples(repl *Repl) {
	fmt.Println(colorizer.MaybeColorizeHelp("Immediately executed statements", true))
	fmt.Println(
		`[mlr] 1+2
3

[mlr] x=3  # These are local variables
[mlr] y=4
[mlr] x+y
7`)
	fmt.Println()
	fmt.Println(colorizer.MaybeColorizeHelp("Defining functions:", true))
	fmt.Println(
		`[mlr] <
func f(a,b) {
  return a**b
}
>
[mlr] f(7,5)
16807`)
	fmt.Println()
	fmt.Println(colorizer.MaybeColorizeHelp("Reading and processing records:", true))
	fmt.Println(
		`[mlr] :open foo.dat
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
[mlr] $z = f($x, $i)
[mlr] :write
a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463,z=0.021160211005187134`)
	fmt.Println()
}
