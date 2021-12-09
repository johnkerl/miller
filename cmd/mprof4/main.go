// Experiments in performance/profiling.
package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strconv"

	"github.com/pkg/profile" // for trace.out

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/input"
	"github.com/johnkerl/miller/internal/pkg/output"
	"github.com/johnkerl/miller/internal/pkg/types"
)

func main() {

	// Respect env $GOMAXPROCS, if provided, else set default.
	haveSetGoMaxProcs := false
	goMaxProcsString := os.Getenv("GOMAXPROCS")
	if goMaxProcsString != "" {
		goMaxProcs, err := strconv.Atoi(goMaxProcsString)
		if err != nil {
			runtime.GOMAXPROCS(goMaxProcs)
			haveSetGoMaxProcs = true
		}
	}
	if !haveSetGoMaxProcs {
		// As of Go 1.16 this is the default anyway. For 1.15 and below we need
		// to explicitly set this.
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	debug.SetGCPercent(500) // Empirical: See README-profiling.md

	if os.Getenv("MPROF_PPROF") != "" {
		// profiling with cpu.pprof and go tool pprof -http=:8080 cpu.pprof
		profFilename := "cpu.pprof"
		handle, err := os.Create(profFilename)
		if err != nil {
			fmt.Fprintln(os.Stderr, os.Args[0], ": ", "Could not start CPU profile: ", err)
			return
		}
		defer handle.Close()

		if err := pprof.StartCPUProfile(handle); err != nil {
			fmt.Fprintln(os.Stderr, os.Args[0], ": ", "Could not start CPU profile: ", err)
			return
		}
		defer pprof.StopCPUProfile()

		fmt.Fprintf(os.Stderr, "CPU profile started.\n")
		fmt.Fprintf(os.Stderr, "go tool pprof -http=:8080 cpu.pprof\n")
		defer fmt.Fprintf(os.Stderr, "CPU profile finished.\n")
	}

	if os.Getenv("MPROF_TRACE") != "" {
		// tracing with trace.out and go tool trace trace.out
		fmt.Fprintf(os.Stderr, "go tool trace trace.out\n")
		defer profile.Start(profile.TraceProfile, profile.ProfilePath(".")).Stop()
	}

	options := cli.DefaultOptions()
	types.SetInferrerStringOnly()

	filenames := os.Args[1:]

	err := Stream(filenames, options, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mlr: %v.\n", err)
		os.Exit(1)
	}
}

func getBatchSize() int {
	return 1000
}

func Stream(
	filenames []string,
	options *cli.TOptions,
	outputStream io.WriteCloser,
) error {
	initialContext := types.NewContext()

	// Instantiate the record-reader
	recordReader, err := input.NewRecordReaderDKVP(&options.ReaderOptions, getBatchSize())
	if err != nil {
		return err
	}

	// Instantiate the record-writer
	recordWriter, err := output.NewRecordWriterDKVP(&options.WriterOptions)
	if err != nil {
		return err
	}

	bufferedOutputStream := bufio.NewWriter(os.Stdout)
	defer bufferedOutputStream.Flush()

	ioChannel := make(chan *list.List, 1)
	downstreamDoneChannel := make(chan bool, 0)
	errorChannel := make(chan error, 1)
	doneWritingChannel := make(chan bool, 1)

	go recordReader.Read(filenames, *initialContext, ioChannel, errorChannel, downstreamDoneChannel)
	go output.ChannelWriter(ioChannel, recordWriter, &options.WriterOptions, doneWritingChannel,
		bufferedOutputStream, true)

	done := false
	for !done {
		select {
		case err := <-errorChannel:
			fmt.Fprintln(os.Stderr, "mlr", ": ", err)
			os.Exit(1)
		case _ = <-doneWritingChannel:
			done = true
			break
		}
	}

	return nil
}
