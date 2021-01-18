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

// TODO: many more comments
// TODO: rename classes
// TODO: pipe and write/append need different maps. else, '> "cat"' would be weird.

package cst

import (
	"errors"
	"fmt"
	"io"
	"os"

	"miller/lib"
)

// ================================================================
type OutputHandlerManager interface {
	Print(outputString string, filename string) error
	Close() []error
}

type OutputHandler interface {
	Print(outputString string) error
	Close() error
}

// ================================================================
type MultiOutputHandlerManager struct {
	outputHandlers map[string]OutputHandler
	// TOOD: make an enum
	append bool
	pipe   bool
}

func NewFileWritetHandlerManager() *MultiOutputHandlerManager {
	return &MultiOutputHandlerManager{
		outputHandlers: make(map[string]OutputHandler),
		append:         false,
		pipe:           false,
	}
}

func NewFileAppendHandlerManager() *MultiOutputHandlerManager {
	return &MultiOutputHandlerManager{
		outputHandlers: make(map[string]OutputHandler),
		append:         true,
		pipe:           false,
	}
}

func NewPipeWriteHandlerManager() *MultiOutputHandlerManager {
	return &MultiOutputHandlerManager{
		outputHandlers: make(map[string]OutputHandler),
		append:         false,
		pipe:           true,
	}
}

func (this *MultiOutputHandlerManager) Print(
	outputString string,
	filename string,
) error {
	// TODO: LRU with close and re-open in case too many files are open
	outputHandler := this.outputHandlers[filename]
	if outputHandler == nil {
		var err error = nil
		if this.pipe {
			outputHandler, err = NewPipeWriteOutputHandler(filename)
			if err != nil {
				return err
			}
			if outputHandler != nil {
			}
		} else if this.append {
			outputHandler, err = NewFileAppendOutputHandler(filename)
			if err != nil {
				return err
			}
		} else {
			outputHandler, err = NewFileWriteOutputHandler(filename)
			if err != nil {
				return err
			}
		}
		this.outputHandlers[filename] = outputHandler
	}
	return outputHandler.Print(outputString)
}

func (this *MultiOutputHandlerManager) Close() []error {
	errs := make([]error, 0)
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
}

func (this FileOutputHandler) Print(outputString string) error {
	_, err := this.handle.Write([]byte(outputString))
	return err
}
func (this *FileOutputHandler) Close() error {
	if this.closeable {
		return this.handle.Close()
	} else {
		return nil
	}
}

// ----------------------------------------------------------------
func NewFileWriteOutputHandler(
	filename string,
) (*FileOutputHandler, error) {
	handle, err := os.OpenFile(
		filename,
		os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil, err
	}
	return &FileOutputHandler{
		filename:  filename,
		handle:    handle,
		closeable: true,
	}, nil
}

func NewFileAppendOutputHandler(
	filename string,
) (*FileOutputHandler, error) {
	handle, err := os.OpenFile(
		filename,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, err
	}
	return &FileOutputHandler{
		filename: filename,
		handle:   handle,
	}, nil
}

func NewPipeWriteOutputHandler(
	commandString string,
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

	return &FileOutputHandler{
		filename:  "| " + commandString,
		handle:    writePipe,
		closeable: true,
	}, nil
}
