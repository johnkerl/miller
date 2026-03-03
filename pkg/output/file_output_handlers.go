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
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/types"
)

// Maximum number of file handlers to keep open when writing to many files
// (e.g. mlr split -g field). Evicted handlers are closed; re-opened in append
// mode when needed again. Prevents "too many open files" (issue #1105).
const lruFileHandlerCapacity = 256

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

type lruNode struct {
	key     string
	handler *FileOutputHandler
	prev    *lruNode
	next    *lruNode
}

type MultiOutputHandlerManager struct {
	outputHandlers map[string]OutputHandler

	// For stdout or stderr
	singleHandler *FileOutputHandler

	// TODO: make an enum
	append              bool // True for ">>", false for ">" and "|"
	pipe                bool // True for "|", false for ">" and ">>"
	recordWriterOptions *cli.TWriterOptions

	// LRU cache for file handlers: limit open FDs, re-open evicted files in append mode
	mu               sync.Mutex
	lruNodes         map[string]*lruNode // filename -> node (file mode only)
	lruHead          *lruNode            // MRU
	lruTail          *lruNode            // LRU
	evictedFilenames map[string]bool     // filenames we closed; re-open with append
}

func NewFileOutputHandlerManager(
	recordWriterOptions *cli.TWriterOptions,
	doAppend bool,
) *MultiOutputHandlerManager {
	if doAppend {
		return NewFileAppendHandlerManager(recordWriterOptions)
	}
	return NewFileWritetHandlerManager(recordWriterOptions)
}

func NewFileWritetHandlerManager(
	recordWriterOptions *cli.TWriterOptions,
) *MultiOutputHandlerManager {
	return &MultiOutputHandlerManager{
		outputHandlers:      make(map[string]OutputHandler),
		singleHandler:       nil,
		append:              false,
		pipe:                false,
		recordWriterOptions: recordWriterOptions,
		lruNodes:            make(map[string]*lruNode),
		evictedFilenames:    make(map[string]bool),
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
		lruNodes:            make(map[string]*lruNode),
		evictedFilenames:    make(map[string]bool),
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

func (mgr *MultiOutputHandlerManager) WriteString(
	outputString string,
	filename string,
) error {
	outputHandler, err := mgr.getOutputHandlerFor(filename)
	if err != nil {
		return err
	}
	return outputHandler.WriteString(outputString)
}

func (mgr *MultiOutputHandlerManager) WriteRecordAndContext(
	outrecAndContext *types.RecordAndContext,
	filename string,
) error {
	outputHandler, err := mgr.getOutputHandlerFor(filename)
	if err != nil {
		return err
	}
	return outputHandler.WriteRecordAndContext(outrecAndContext)
}

func (mgr *MultiOutputHandlerManager) getOutputHandlerFor(
	filename string,
) (OutputHandler, error) {
	if mgr.singleHandler != nil {
		return mgr.singleHandler, nil
	}

	mgr.mu.Lock()

	// Pipe mode: no LRU (pipes cannot be re-opened)
	if mgr.pipe {
		outputHandler := mgr.outputHandlers[filename]
		if outputHandler == nil {
			var err error
			outputHandler, err = NewPipeWriteOutputHandler(filename, mgr.recordWriterOptions)
			if err != nil {
				mgr.mu.Unlock()
				return nil, err
			}
			mgr.outputHandlers[filename] = outputHandler
		}
		mgr.mu.Unlock()
		return outputHandler, nil
	}

	// File mode: LRU cache with eviction and append-mode re-open
	node := mgr.lruNodes[filename]
	if node != nil {
		mgr.lruTouch(node)
		mgr.mu.Unlock()
		return node.handler, nil
	}

	// Cache miss: evict LRU if at capacity
	if len(mgr.lruNodes) >= lruFileHandlerCapacity && mgr.lruTail != nil {
		tail := mgr.lruTail
		mgr.lruRemove(tail)
		mgr.evictedFilenames[tail.key] = true
		handlerToClose := tail.handler
		mgr.mu.Unlock()
		_ = handlerToClose.Close() // may block on channel drain
		mgr.mu.Lock()
	}

	// Create new handler: use append if re-opening after eviction
	useAppend := mgr.append || mgr.evictedFilenames[filename]
	if mgr.evictedFilenames[filename] {
		delete(mgr.evictedFilenames, filename)
	}

	var handler *FileOutputHandler
	var err error
	if useAppend {
		handler, err = NewFileAppendOutputHandler(filename, mgr.recordWriterOptions)
	} else {
		handler, err = NewFileWriteOutputHandler(filename, mgr.recordWriterOptions)
	}
	if err != nil {
		mgr.mu.Unlock()
		return nil, err
	}

	node = &lruNode{key: filename, handler: handler}
	mgr.lruInsert(node)
	mgr.mu.Unlock()
	return handler, nil
}

func (mgr *MultiOutputHandlerManager) lruTouch(node *lruNode) {
	if mgr.lruHead == node {
		return
	}
	mgr.lruRemove(node)
	mgr.lruInsert(node)
}

func (mgr *MultiOutputHandlerManager) lruRemove(node *lruNode) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		mgr.lruHead = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	} else {
		mgr.lruTail = node.prev
	}
	delete(mgr.lruNodes, node.key)
}

func (mgr *MultiOutputHandlerManager) lruInsert(node *lruNode) {
	mgr.lruNodes[node.key] = node
	node.prev = nil
	node.next = mgr.lruHead
	if mgr.lruHead != nil {
		mgr.lruHead.prev = node
	} else {
		mgr.lruTail = node
	}
	mgr.lruHead = node
}

func (mgr *MultiOutputHandlerManager) Close() []error {
	errs := []error{}
	if mgr.singleHandler != nil {
		err := mgr.singleHandler.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}
	mgr.mu.Lock()
	handlersToClose := make([]OutputHandler, 0, len(mgr.outputHandlers)+len(mgr.lruNodes))
	for _, oh := range mgr.outputHandlers {
		handlersToClose = append(handlersToClose, oh)
	}
	for _, node := range mgr.lruNodes {
		handlersToClose = append(handlersToClose, node.handler)
	}
	mgr.outputHandlers = make(map[string]OutputHandler)
	mgr.lruNodes = make(map[string]*lruNode)
	mgr.lruHead, mgr.lruTail = nil, nil
	mgr.mu.Unlock()
	for _, outputHandler := range handlersToClose {
		err := outputHandler.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

type FileOutputHandler struct {
	filename             string
	handle               io.WriteCloser
	bufferedOutputStream *bufio.Writer
	closeable            bool

	// This will be nil if WriteRecordAndContext has never been called. It's
	// lazily created on WriteRecord. The record-writer / channel parts are
	// called only by WriteRecrod which is called by emit and tee variants;
	// print and dump variants call WriteString.
	recordWriterOptions  *cli.TWriterOptions
	recordWriter         IRecordWriter
	recordOutputChannel  chan []*types.RecordAndContext // list of *types.RecordAndContext
	recordDoneChannel    chan bool
	recordErroredChannel chan bool
}

func newOutputHandlerCommon(
	filename string,
	handle io.WriteCloser,
	closeable bool,
	recordWriterOptions *cli.TWriterOptions,
) *FileOutputHandler {
	return &FileOutputHandler{
		filename:             filename,
		handle:               handle,
		bufferedOutputStream: bufio.NewWriter(handle),
		closeable:            closeable,

		recordWriterOptions:  recordWriterOptions,
		recordWriter:         nil,
		recordOutputChannel:  nil,
		recordDoneChannel:    nil,
		recordErroredChannel: nil,
	}
}

func NewFileOutputHandler(
	filename string,
	recordWriterOptions *cli.TWriterOptions,
	doAppend bool,
) (*FileOutputHandler, error) {
	if doAppend {
		return NewFileAppendOutputHandler(filename, recordWriterOptions)
	}
	return NewFileWriteOutputHandler(filename, recordWriterOptions)
}

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
		return nil, fmt.Errorf(`could not launch command "%s" for pipe-to`, commandString)
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

func (handler *FileOutputHandler) WriteString(outputString string) error {
	_, err := handler.bufferedOutputStream.WriteString(outputString)
	return err
}

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

	// TODO: myybe refactor to batch better
	handler.recordOutputChannel <- []*types.RecordAndContext{outrecAndContext}
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

	handler.recordOutputChannel = make(chan []*types.RecordAndContext, 1) // list of *types.RecordAndContext
	handler.recordDoneChannel = make(chan bool, 1)
	handler.recordErroredChannel = make(chan bool, 1)

	go ChannelWriter(
		handler.recordOutputChannel,
		handler.recordWriter,
		handler.recordWriterOptions,
		handler.recordDoneChannel,
		handler.recordErroredChannel,
		handler.bufferedOutputStream,
		false, // outputIsStdout
	)

	return nil
}

func (handler *FileOutputHandler) Close() (retval error) {
	retval = nil

	if handler.recordOutputChannel != nil {
		// TODO: see if we need a real context
		emptyContext := types.Context{}
		handler.recordOutputChannel <- types.NewEndOfStreamMarkerList(&emptyContext)

		// Wait for the output channel to drain
		done := false
		for !done {
			select {
			case <-handler.recordErroredChannel:
				done = true
				retval = errors.New("exiting due to data error") // details already printed
			case <-handler.recordDoneChannel:
				done = true
			}
		}
	}

	if retval != nil {
		return retval
	}

	handler.bufferedOutputStream.Flush()
	if handler.closeable {
		return handler.handle.Close()
	} // e.g. stdout
	return nil
}
