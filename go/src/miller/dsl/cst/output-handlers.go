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

package cst

import (
	"fmt"
	"os"
)

// ----------------------------------------------------------------
type OutputHandler interface {
	Print(outputString string)
	Close() error
}

// ----------------------------------------------------------------
type StdoutOutputHandler struct {
}

func NewStdoutOutputHandler() (*StdoutOutputHandler, error) {
	return &StdoutOutputHandler{}, nil
}
func (this *StdoutOutputHandler) Print(outputString string) {
	fmt.Fprint(os.Stdout, outputString)
}
func (this *StdoutOutputHandler) Close() {
	// No-op
}

// ----------------------------------------------------------------
type StderrOutputHandler struct {
}

func NewStderrOutputHandler() (*StderrOutputHandler, error) {
	return &StderrOutputHandler{}, nil
}
func (this *StderrOutputHandler) Print(outputString string) {
	fmt.Fprint(os.Stderr, outputString)
}
func (this *StderrOutputHandler) Close() {
	// No-op
}

// ----------------------------------------------------------------
type FileWriteOutputHandler struct {
	filename string
}

func NewFileWriteOutputHandler(
	filename string,
) (*FileWriteOutputHandler, error) {
	// TODO: stub
	return &FileWriteOutputHandler{
		filename: filename,
	}, nil
}

func (this *FileWriteOutputHandler) Print(outputString string) {
	// TODO: stub
}
func (this *FileWriteOutputHandler) Close() {
	// TODO: stub
}

// ----------------------------------------------------------------
type FileAppendOutputHandler struct {
	filename string
}

func NewFileAppendOutputHandler(
	filename string,
) (*FileAppendOutputHandler, error) {
	// TODO: stub
	return &FileAppendOutputHandler{
		filename: filename,
	}, nil
}

func (this *FileAppendOutputHandler) Print(outputString string) {
	// TODO: stub
}
func (this *FileAppendOutputHandler) Close() {
	// TODO: stub
}

// ----------------------------------------------------------------
type PipeToCommandOutputHandler struct {
	command string
}

func NewPipeToCommandOutputHandler(
	command string,
) (*PipeToCommandOutputHandler, error) {
	// TODO: stub
	return &PipeToCommandOutputHandler{
		command: command,
	}, nil
}

func (this *PipeToCommandOutputHandler) Print(outputString string) {
	// TODO: stub
}
func (this *PipeToCommandOutputHandler) Close() {
	// TODO: stub
}
