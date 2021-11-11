// ================================================================
// These are handlers for print, dump, emit, etc in the put/filter verbs.
//
// * For "> filename" ">> filename", these handle the open/write/close file operations.
// * For "| command", these handle open/write/close pipe operations.
// * For stderr, these write to stderr immediately.
// * For stdout, these write to the main record-output Go channel.
//   The reason for this is since we want all print statements and
//   record-output to be in the same goroutine, for deterministic output
//   ordering. (Main record-writer output is also to stdout.)
//   ================================================================

package output

import (
	"errors"
	"fmt"
	"io"
	"os"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

// ================================================================
type OutputHandlerManager interface {

	// For print-variants and dump-variants
	WriteString(outputString string, filename string) error

	// For emit-variants and tee
	WriteRecordAndContext(outrecAndContext *types.RecordAndContext, filename string) error

	Close() []error
}

type OutputHandler interface {
	WriteString(outputString string) error
	WriteRecordAndContext(outrecAndContext *types.RecordAndContext) error
	Close() error
}

// ================================================================
type MultiOutputHandlerManager struct {
	outputHandlers map[string]OutputHandler

	// For stdout or stderr
	singleHandler *FileOutputHandler

	// TOOD: make an enum
	append              bool // True for ">>", false for ">" and "|"
	pipe                bool // True for "|", false for ">" and ">>"
	recordWriterOptions *cli.TWriterOptions
}

// ----------------------------------------------------------------
func NewFileWritetHandlerManager(
	recordWriterOptions *cli.TWriterOptions,
) *MultiOutputHandlerManager {
	return &MultiOutputHandlerManager{
		outputHandlers:      make(map[string]OutputHandler),
		singleHandler:       nil,
		append:              false,
		pipe:                false,
		recordWriterOptions: recordWriterOptions,
	}
}

func NewFileAppendHandlerManager(
	recordWriterOptions *cli.TWriterOptions,
) *MultiOutputHandlerManager {
	return &MultiOutputHandlerManager{
		outputHandlers:      make(map[string]OutputHandler),
		singleHandler:       nil,
		append:              true,
		pipe:                false,
		recordWriterOptions: recordWriterOptions,
	}
}

func NewPipeWriteHandlerManager(
	recordWriterOptions *cli.TWriterOptions,
) *MultiOutputHandlerManager {
	return &MultiOutputHandlerManager{
		outputHandlers:      make(map[string]OutputHandler),
		singleHandler:       nil,
		append:              false,
		pipe:                true,
		recordWriterOptions: recordWriterOptions,
	}
}

func NewStdoutWriteHandlerManager(
	recordWriterOptions *cli.TWriterOptions,
) *MultiOutputHandlerManager {
	return &MultiOutputHandlerManager{
		outputHandlers:      make(map[string]OutputHandler),
		singleHandler:       newStdoutOutputHandler(recordWriterOptions),
		append:              false,
		pipe:                false,
		recordWriterOptions: recordWriterOptions,
	}
}

func NewStderrWriteHandlerManager(
	recordWriterOptions *cli.TWriterOptions,
) *MultiOutputHandlerManager {
	return &MultiOutputHandlerManager{
		outputHandlers:      make(map[string]OutputHandler),
		singleHandler:       newStderrOutputHandler(recordWriterOptions),
		append:              false,
		pipe:                false,
		recordWriterOptions: recordWriterOptions,
	}
}

// ----------------------------------------------------------------
func (manager *MultiOutputHandlerManager) WriteString(
	outputString string,
	filename string,
) error {
	outputHandler, err := manager.getOutputHandlerFor(filename)
	if err != nil {
		return err
	}
	return outputHandler.WriteString(outputString)
}

func (manager *MultiOutputHandlerManager) WriteRecordAndContext(
	outrecAndContext *types.RecordAndContext,
	filename string,
) error {
	outputHandler, err := manager.getOutputHandlerFor(filename)
	if err != nil {
		return err
	}
	return outputHandler.WriteRecordAndContext(outrecAndContext)
}

func (manager *MultiOutputHandlerManager) getOutputHandlerFor(
	filename string,
) (OutputHandler, error) {
	if manager.singleHandler != nil {
		return manager.singleHandler, nil
	}

	// TODO: LRU with close and re-open in case too many files are open
	outputHandler := manager.outputHandlers[filename]
	if outputHandler == nil {
		var err error = nil
		if manager.pipe {
			outputHandler, err = NewPipeWriteOutputHandler(
				filename,
				manager.recordWriterOptions,
			)
			if err != nil {
				return nil, err
			}
			if outputHandler != nil {
			}
		} else if manager.append {
			outputHandler, err = NewFileAppendOutputHandler(
				filename,
				manager.recordWriterOptions,
			)
			if err != nil {
				return nil, err
			}
		} else {
			outputHandler, err = NewFileWriteOutputHandler(
				filename,
				manager.recordWriterOptions,
			)
			if err != nil {
				return nil, err
			}
		}
		manager.outputHandlers[filename] = outputHandler
	}
	return outputHandler, nil
}

func (manager *MultiOutputHandlerManager) Close() []error {
	errs := make([]error, 0)
	if manager.singleHandler != nil {
		err := manager.singleHandler.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}
	for _, outputHandler := range manager.outputHandlers {
		err := outputHandler.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// ================================================================
type FileOutputHandler struct {
	filename  string
	handle    io.WriteCloser
	closeable bool

	// This will be nil if WriteRecordAndContext has never been called. It's
	// lazily created on WriteRecord. The record-writer / channel parts are
	// called only by WriteRecrod which is called by emit and tee variants;
	// print and dump variants call WriteString.
	recordWriterOptions *cli.TWriterOptions
	recordWriter        IRecordWriter
	recordOutputChannel chan *types.RecordAndContext
	recordDoneChannel   chan bool
}

func newOutputHandlerCommon(
	filename string,
	handle io.WriteCloser,
	closeable bool,
	recordWriterOptions *cli.TWriterOptions,
) *FileOutputHandler {
	return &FileOutputHandler{
		filename:  filename,
		handle:    handle,
		closeable: closeable,

		recordWriterOptions: recordWriterOptions,
		recordWriter:        nil,
		recordOutputChannel: nil,
		recordDoneChannel:   nil,
	}
}

// ----------------------------------------------------------------
func NewFileWriteOutputHandler(
	filename string,
	recordWriterOptions *cli.TWriterOptions,
) (*FileOutputHandler, error) {
	handle, err := os.OpenFile(
		filename,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0644, // TODO: let users parameterize this
	)
	if err != nil {
		return nil, err
	}
	return newOutputHandlerCommon(
		filename,
		handle,
		true,
		recordWriterOptions,
	), nil
}

func NewFileAppendOutputHandler(
	filename string,
	recordWriterOptions *cli.TWriterOptions,
) (*FileOutputHandler, error) {
	handle, err := os.OpenFile(
		filename,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644, // TODO: let users parameterize this
	)
	if err != nil {
		return nil, err
	}
	return newOutputHandlerCommon(
		filename,
		handle,
		true,
		recordWriterOptions,
	), nil
}

func NewPipeWriteOutputHandler(
	commandString string,
	recordWriterOptions *cli.TWriterOptions,
) (*FileOutputHandler, error) {
	writePipe, err := lib.OpenOutboundHalfPipe(commandString)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf(
				"%s: could not launch command \"%s\" for pipe-to.",
				"mlr",
				commandString,
			),
		)
	}

	return newOutputHandlerCommon(
		"| "+commandString,
		writePipe,
		true,
		recordWriterOptions,
	), nil
}

func newStdoutOutputHandler(
	recordWriterOptions *cli.TWriterOptions,
) *FileOutputHandler {
	return newOutputHandlerCommon(
		"(stdout)",
		os.Stdout,
		false,
		recordWriterOptions,
	)
}

func newStderrOutputHandler(
	recordWriterOptions *cli.TWriterOptions,
) *FileOutputHandler {
	return newOutputHandlerCommon(
		"(stderr)",
		os.Stderr,
		false,
		recordWriterOptions,
	)
}

// ----------------------------------------------------------------
func (handler *FileOutputHandler) WriteString(outputString string) error {
	_, err := handler.handle.Write([]byte(outputString))
	return err
}

// ----------------------------------------------------------------
func (handler *FileOutputHandler) WriteRecordAndContext(
	outrecAndContext *types.RecordAndContext,
) error {
	// Lazily create the record-writer and output channel.
	if handler.recordWriter == nil {
		err := handler.setUpRecordWriter()
		if err != nil {
			return err
		}
	}

	handler.recordOutputChannel <- outrecAndContext
	return nil
}

func (handler *FileOutputHandler) setUpRecordWriter() error {
	if handler.recordWriter != nil {
		return nil
	}

	recordWriter, err := Create(handler.recordWriterOptions)
	if err != nil {
		return err
	}
	handler.recordWriter = recordWriter

	handler.recordOutputChannel = make(chan *types.RecordAndContext, 1)
	handler.recordDoneChannel = make(chan bool, 1)

	go ChannelWriter(
		handler.recordOutputChannel,
		handler.recordWriter,
		handler.recordDoneChannel,
		handler.handle,
		false, // outputIsStdout
	)

	return nil
}

// ----------------------------------------------------------------
func (handler *FileOutputHandler) Close() error {
	if handler.recordOutputChannel != nil {
		// TODO: see if we need a real context
		emptyContext := types.Context{}
		handler.recordOutputChannel <- types.NewEndOfStreamMarker(&emptyContext)

		// Wait for the output channel to drain
		done := false
		for !done {
			select {
			case _ = <-handler.recordDoneChannel:
				done = true
				break
			}
		}
	}

	if handler.closeable {
		return handler.handle.Close()
	} else {
		return nil
	}
}
