// ================================================================
// Top-level handler for a REPL session, including setup/construction, and
// ingesting command-lines. Command-line strings are triaged and send off to
// the appropriate handlers: DSL parse/execute if the command is a DSL statement
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

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/dsl/cst"
	"github.com/johnkerl/miller/internal/pkg/input"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/output"
	"github.com/johnkerl/miller/internal/pkg/runtime"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ----------------------------------------------------------------
func NewRepl(
	exeName string,
	replName string,
	showStartupBanner bool,
	showPrompts bool,
	astPrintMode ASTPrintMode,
	doWarnings bool,
	strictMode bool,
	options *cli.TOptions,
	recordOutputFileName string,
	recordOutputStream *os.File,
) (*Repl, error) {

	recordReader, err := input.Create(&options.ReaderOptions, 1) // recordsPerBatch
	if err != nil {
		return nil, err
	}

	recordWriter, err := output.Create(&options.WriterOptions)
	if err != nil {
		return nil, err
	}

	// $* is the empty map {} until/unless the user opens a file and reads records from it.
	inrec := mlrval.NewMlrmapAsRecord()
	// NR is 0, etc until/unless the user opens a file and reads records from it.
	context := types.NewContext()

	runtimeState := runtime.NewEmptyState(options, strictMode)
	runtimeState.Update(inrec, context)
	// The filter expression for the main Miller DSL is any non-assignment
	// statement like 'true' or '$x > 0.5' etc. For the REPL, we re-use this for
	// interactive expressions to be printed to the terminal. For the main DSL,
	// the default is mlrval.FromTrue(); for the REPL, the default is
	// mlrval.NULL.
	runtimeState.FilterExpression = mlrval.NULL

	// For control-C handling
	sysToSignalHandlerChannel := make(chan os.Signal, 1) // Our signal handler reads system notification here
	appSignalNotificationChannel := make(chan bool, 1)   // Our signal handler writes this for our app to poll
	signal.Notify(sysToSignalHandlerChannel, os.Interrupt, syscall.SIGTERM)
	go controlCHandler(sysToSignalHandlerChannel, appSignalNotificationChannel)

	cstRootNode := cst.NewEmptyRoot(
		&options.WriterOptions, cst.DSLInstanceTypeREPL,
	).WithRedefinableUDFUDS().WithStrictMode(strictMode)

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

		options:       options,
		readerChannel: nil,
		errorChannel:  nil,
		recordReader:  recordReader,
		recordWriter:  recordWriter,

		runtimeState:                 runtimeState,
		sysToSignalHandlerChannel:    sysToSignalHandlerChannel,
		appSignalNotificationChannel: appSignalNotificationChannel,
	}

	repl.setBufferedOutputStream(recordOutputFileName, recordOutputStream)

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
func (repl *Repl) handleSession(istream *os.File) error {
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
			return err
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
			err = repl.handleMultiLine(lineReader, ">", true) // multi-line immediate
			if err != nil {
				return err
			}
		} else if trimmedLine == "<<" {
			err = repl.handleMultiLine(lineReader, ">>", false) // multi-line non-immediate
			if err != nil {
				return err
			}
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
	return nil
}

// ----------------------------------------------------------------
// Context: the "<" or "<<" has already been seen. we read until ">" or ">>".

func (repl *Repl) handleMultiLine(
	lineReader *bufio.Reader,
	terminator string,
	doImmediate bool,
) error {
	var buffer bytes.Buffer
	for {
		repl.printPrompt2()

		line, err := lineReader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
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

	return nil
}

func (repl *Repl) setBufferedOutputStream(
	recordOutputFileName string,
	recordOutputStream *os.File,
) {
	repl.recordOutputFileName = recordOutputFileName
	repl.recordOutputStream = recordOutputStream
	repl.bufferedRecordOutputStream = bufio.NewWriter(recordOutputStream)
}

func (repl *Repl) closeBufferedOutputStream() error {
	if repl.recordOutputStream != os.Stdout {
		err := repl.recordOutputStream.Close()
		if err != nil {
			return fmt.Errorf("mlr repl: error on redirect close of %s: %v\n",
				repl.recordOutputFileName, err,
			)
		}
	}
	return nil
}
