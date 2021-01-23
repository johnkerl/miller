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

	"miller/clitypes"
	"miller/lib"
	"miller/types"
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
	recordWriterOptions *clitypes.TWriterOptions
}

// ----------------------------------------------------------------
func NewFileWritetHandlerManager(
	recordWriterOptions *clitypes.TWriterOptions,
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
	recordWriterOptions *clitypes.TWriterOptions,
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
	recordWriterOptions *clitypes.TWriterOptions,
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
	recordWriterOptions *clitypes.TWriterOptions,
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
	recordWriterOptions *clitypes.TWriterOptions,
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
func (this *MultiOutputHandlerManager) WriteString(
	outputString string,
	filename string,
) error {
	outputHandler, err := this.getOutputHandlerFor(filename)
	if err != nil {
		return err
	}
	return outputHandler.WriteString(outputString)
}

func (this *MultiOutputHandlerManager) WriteRecordAndContext(
	outrecAndContext *types.RecordAndContext,
	filename string,
) error {
	outputHandler, err := this.getOutputHandlerFor(filename)
	if err != nil {
		return err
	}
	return outputHandler.WriteRecordAndContext(outrecAndContext)
}

func (this *MultiOutputHandlerManager) getOutputHandlerFor(
	filename string,
) (OutputHandler, error) {
	if this.singleHandler != nil {
		return this.singleHandler, nil
	}

	// TODO: LRU with close and re-open in case too many files are open
	outputHandler := this.outputHandlers[filename]
	if outputHandler == nil {
		var err error = nil
		if this.pipe {
			outputHandler, err = NewPipeWriteOutputHandler(
				filename,
				this.recordWriterOptions,
			)
			if err != nil {
				return nil, err
			}
			if outputHandler != nil {
			}
		} else if this.append {
			outputHandler, err = NewFileAppendOutputHandler(
				filename,
				this.recordWriterOptions,
			)
			if err != nil {
				return nil, err
			}
		} else {
			outputHandler, err = NewFileWriteOutputHandler(
				filename,
				this.recordWriterOptions,
			)
			if err != nil {
				return nil, err
			}
		}
		this.outputHandlers[filename] = outputHandler
	}
	return outputHandler, nil
}

func (this *MultiOutputHandlerManager) Close() []error {
	errs := make([]error, 0)
	if this.singleHandler != nil {
		err := this.singleHandler.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}
	for _, outputHandler := range this.outputHandlers {
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
	recordWriterOptions *clitypes.TWriterOptions
	recordWriter        IRecordWriter
	recordOutputChannel chan *types.RecordAndContext
	recordDoneChannel   chan bool
}

func newOutputHandlerCommon(
	filename string,
	handle io.WriteCloser,
	closeable bool,
	recordWriterOptions *clitypes.TWriterOptions,
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
	recordWriterOptions *clitypes.TWriterOptions,
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
	recordWriterOptions *clitypes.TWriterOptions,
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
	recordWriterOptions *clitypes.TWriterOptions,
) (*FileOutputHandler, error) {
	writePipe, err := lib.OpenOutboundHalfPipe(commandString)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf(
				"%s: could not launch command \"%s\" for pipe-to.",
				os.Args[0],
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
	recordWriterOptions *clitypes.TWriterOptions,
) *FileOutputHandler {
	return newOutputHandlerCommon(
		"(stdout)",
		os.Stdout,
		false,
		recordWriterOptions,
	)
}

func newStderrOutputHandler(
	recordWriterOptions *clitypes.TWriterOptions,
) *FileOutputHandler {
	return newOutputHandlerCommon(
		"(stderr)",
		os.Stderr,
		false,
		recordWriterOptions,
	)
}

// ----------------------------------------------------------------
func (this *FileOutputHandler) WriteString(outputString string) error {
	_, err := this.handle.Write([]byte(outputString))
	return err
}

// ----------------------------------------------------------------
func (this *FileOutputHandler) WriteRecordAndContext(
	outrecAndContext *types.RecordAndContext,
) error {
	// Lazily create the record-writer and output channel.
	if this.recordWriter == nil {
		err := this.setUpRecordWriter()
		if err != nil {
			return err
		}
	}

	this.recordOutputChannel <- outrecAndContext
	return nil
}

func (this *FileOutputHandler) setUpRecordWriter() error {
	if this.recordWriter != nil {
		return nil
	}

	recordWriter := Create(this.recordWriterOptions)
	if recordWriter == nil {
		return errors.New(
			"Output format not found: " + this.recordWriterOptions.OutputFileFormat,
		)
	}
	this.recordWriter = recordWriter

	this.recordOutputChannel = make(chan *types.RecordAndContext, 1)
	this.recordDoneChannel = make(chan bool, 1)

	go ChannelWriter(
		this.recordOutputChannel,
		this.recordWriter,
		this.recordDoneChannel,
		this.handle,
	)

	return nil
}

// ----------------------------------------------------------------
func (this *FileOutputHandler) Close() error {
	if this.recordOutputChannel != nil {
		// TODO: see if we need a real context
		emptyContext := types.Context{}
		this.recordOutputChannel <- types.NewEndOfStreamMarker(&emptyContext)

		// Wait for the output channel to drain
		done := false
		for !done {
			select {
			case _ = <-this.recordDoneChannel:
				done = true
				break
			}
		}
	}

	if this.closeable {
		return this.handle.Close()
	} else {
		return nil
	}
}
