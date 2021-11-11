// ================================================================
// Top-level handler for a REPL session, including setup/construction, and
// ingesting command-lines. Command-line strings are triaged and send off to
// the appropriate handlers: DSL parse/execute if the comand is a DSL statement
// (like '$z = $x + $y'); REPL-command-line parse/execute otherwise (like
// ':open foo.dat' or ':help').
//
// No command-line-history-editing feature is built in but
//
//   rlwrap mlr repl
//
// is a delight. :)
//
// ================================================================

package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/dsl/cst"
	"mlr/internal/pkg/input"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/output"
	"mlr/internal/pkg/runtime"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
func NewRepl(
	exeName string,
	replName string,
	showStartupBanner bool,
	showPrompts bool,
	astPrintMode ASTPrintMode,
	doWarnings bool,
	options *cli.TOptions,
) (*Repl, error) {

	recordReader, err := input.Create(&options.ReaderOptions)
	if err != nil {
		return nil, err
	}

	recordWriter, err := output.Create(&options.WriterOptions)
	if err != nil {
		return nil, err
	}
	outputStream := os.Stdout

	// $* is the empty map {} until/unless the user opens a file and reads records from it.
	inrec := types.NewMlrmapAsRecord()
	// NR is 0, etc until/unless the user opens a file and reads records from it.
	context := types.NewContext(
		options.ReaderOptions.IPS,
		options.ReaderOptions.IFS,
		options.ReaderOptions.IRS,
		options.WriterOptions.OPS,
		options.WriterOptions.OFS,
		options.WriterOptions.ORS,
		options.WriterOptions.FLATSEP,
	)
	runtimeState := runtime.NewEmptyState()
	runtimeState.Update(inrec, context)
	// The filter expression for the main Miller DSL is any non-assignment
	// statment like 'true' or '$x > 0.5' etc. For the REPL, we re-use this for
	// interactive expressions to be printed to the terminal. For the main DSL,
	// the default is types.MlrvalFromTrue(); for the REPL, the default is
	// types.MLRVAL_VOID.
	runtimeState.FilterExpression = types.MLRVAL_VOID

	// For control-C handling
	sysToSignalHandlerChannel := make(chan os.Signal, 1) // Our signal handler reads system notification here
	appSignalNotificationChannel := make(chan bool, 1)   // Our signal handler writes this for our app to poll
	signal.Notify(sysToSignalHandlerChannel, os.Interrupt, syscall.SIGTERM)
	go controlCHandler(sysToSignalHandlerChannel, appSignalNotificationChannel)

	cstRootNode := cst.NewEmptyRoot(&options.WriterOptions, cst.DSLInstanceTypeREPL).WithRedefinableUDFUDS()

	// TODO

	// If there was a --load/--mload on the command line, load those DSL strings here (e.g.
	// someone's local function library).
	dslStrings := make([]string, 0)
	for _, filename := range options.DSLPreloadFileNames {
		theseDSLStrings, err := lib.LoadStringsFromFileOrDir(filename, ".mlr")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s %s: cannot load DSL expression from \"%s\": ",
				exeName, replName, filename)
			return nil, err
		}
		dslStrings = append(dslStrings, theseDSLStrings...)
	}

	repl := &Repl{
		exeName:  exeName,
		replName: replName,

		inputIsTerminal:   getInputIsTerminal(),
		showStartupBanner: showStartupBanner,
		showPrompts:       showPrompts,
		prompt1:           getPrompt1(),
		prompt2:           getPrompt2(),

		astPrintMode: astPrintMode,
		doWarnings:   doWarnings,
		cstRootNode:  cstRootNode,

		options:      options,
		inputChannel: nil,
		errorChannel: nil,
		recordReader: recordReader,
		recordWriter: recordWriter,
		outputStream: outputStream,

		runtimeState:                 runtimeState,
		sysToSignalHandlerChannel:    sysToSignalHandlerChannel,
		appSignalNotificationChannel: appSignalNotificationChannel,
	}

	for _, dslString := range dslStrings {
		err := repl.handleDSLStringBulk(dslString, doWarnings)
		if err != nil {
			// Error message already printed out
			return nil, err
		}
	}

	return repl, nil
}

// When the user types control-C, immediately print something to the screen,
// then also write to a channel while long-running things like :skip and
// :process can check it.
func controlCHandler(sysToSignalHandlerChannel chan os.Signal, appSignalNotificationChannel chan bool) {
	for {
		<-sysToSignalHandlerChannel          // Block until control-C notification is sent by system
		fmt.Println("^C")                    // Acknowledge for user reassurance
		appSignalNotificationChannel <- true // Let our app poll this
	}
}

// ----------------------------------------------------------------
func (repl *Repl) handleSession(istream *os.File) {
	if repl.showStartupBanner {
		repl.printStartupBanner()
	}

	lineReader := bufio.NewReader(istream)

	for {
		repl.printPrompt1()

		line, err := lineReader.ReadString('\n')
		if err == io.EOF {
			if repl.inputIsTerminal {
				fmt.Println()
			}
			break
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s %s: %v", repl.exeName, repl.replName, err)
			os.Exit(1)
		}

		// Acknowledge any control-C's, even if typed at a ready prompt. We
		// need to drain them all out since they're in a channel from the
		// signal-handler goroutine.
		doneDraining := false
		for {
			select {
			case _ = <-repl.appSignalNotificationChannel:
				line = "" // Ignore any partially-entered line -- a ^C should do that
			default:
				doneDraining = true
			}
			if doneDraining {
				break
			}
		}

		// This trims the trailing newline, as well as leading/trailing whitespace:
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "<" {
			repl.handleMultiLine(lineReader, ">", true) // multi-line immediate
		} else if trimmedLine == "<<" {
			repl.handleMultiLine(lineReader, ">>", false) // multi-line non-immediate
		} else if trimmedLine == ":quit" || trimmedLine == ":q" {
			break
		} else if repl.handleNonDSLLine(trimmedLine) {
			// Handled in that method.
		} else {
			// We need the non-trimmed line here since the DSL syntax for comments is '#.*\n'.
			err = repl.handleDSLStringImmediate(line, repl.doWarnings)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

// ----------------------------------------------------------------
// Context: the "<" or "<<" has already been seen. we read until ">" or ">>".

func (repl *Repl) handleMultiLine(
	lineReader *bufio.Reader,
	terminator string,
	doImmediate bool,
) {
	var buffer bytes.Buffer
	for {
		repl.printPrompt2()

		line, err := lineReader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s %s: %v\n", repl.exeName, repl.replName, err)
			os.Exit(1)
		}

		if strings.TrimSpace(line) == terminator {
			break
		}
		buffer.WriteString(line)
	}
	dslString := buffer.String()

	var err error = nil
	if doImmediate {
		err = repl.handleDSLStringImmediate(dslString, repl.doWarnings)
	} else {
		err = repl.handleDSLStringBulk(dslString, repl.doWarnings)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
