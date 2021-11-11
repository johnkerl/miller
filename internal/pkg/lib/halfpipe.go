package lib

import (
	"fmt"
	"os"

	"mlr/internal/pkg/platform"
)

// OpenOutboundHalfPipe returns a handle to a process. Writing to that handle
// writes to the process' stdin. The process' stdout and stderr are the current
// process' stdout and stderr.
//
// This is for pipe-output-redirection in the Miller put/filter DSL.
//
// Note I am not using os.exec.Cmd which is billed as being simpler than using
// os.StartProcess. It may indeed be simpler when you want to handle the
// subprocess' stdin/stdout/stderr all three within the parent process.  Here I
// found it much easier to use os.StartProcess to let the stdout/stderr run
// free.

func OpenOutboundHalfPipe(commandString string) (*os.File, error) {
	readPipe, writePipe, err := os.Pipe()

	var procAttr os.ProcAttr
	procAttr.Files = []*os.File{
		readPipe,
		os.Stdout,
		os.Stderr,
	}

	// /bin/sh -c "..." or cmd /c "..."
	shellRunArray := platform.GetShellRunArray(commandString)

	process, err := os.StartProcess(shellRunArray[0], shellRunArray, &procAttr)
	if err != nil {
		return nil, err
	}

	go process.Wait()

	return writePipe, nil
}

// OpenInboundHalfPipe returns a handle to a process. Reading from that handle
// reads from the process' stdout. The process' stdin and stderr are the
// current process' stdin and stderr.
//
// This is for the Miller prepipe feature.
//
// Note I am not using os.exec.Cmd which is billed as being simpler than using
// os.StartProcess. It may indeed be simpler when you want to handle the
// subprocess' stdin/stdout/stderr all three within the parent process.  Here I
// found it much easier to use os.StartProcess to let the stdin/stderr run
// free.

func OpenInboundHalfPipe(commandString string) (*os.File, error) {
	readPipe, writePipe, err := os.Pipe()

	var procAttr os.ProcAttr
	procAttr.Files = []*os.File{
		os.Stdin,
		writePipe,
		os.Stderr,
	}

	// /bin/sh -c "..." or cmd /c "..."
	shellRunArray := platform.GetShellRunArray(commandString)

	process, err := os.StartProcess(shellRunArray[0], shellRunArray, &procAttr)
	if err != nil {
		return nil, err
	}

	// TODO comment somewhere
	// https://stackoverflow.com/questions/47486128/why-does-io-pipe-continue-to-block-even-when-eof-is-reached

	// TODO comment
	go func(process *os.Process, readPipe *os.File) {
		_, err := process.Wait()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", "mlr", err)
		}
		readPipe.Close()
	}(process, readPipe)

	return readPipe, nil
}

func waitAndClose(process *os.Process, readPipe *os.File) {
	_, err := process.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", "mlr", err)
	}
	readPipe.Close()
}
