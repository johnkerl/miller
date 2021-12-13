// Experiments in performance/profiling.
package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strconv"

	"github.com/pkg/profile" // for trace.out

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/stream"
	"github.com/johnkerl/miller/internal/pkg/transformers"
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

	cat, err := transformers.NewTransformerCat(false, "", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mprof5: %v\n", err)
		os.Exit(1)
	}
	xforms := []transformers.IRecordTransformer{cat}

	filenames := os.Args[1:]

	err = stream.Stream(filenames, options, xforms, os.Stdout, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mlr: %v.\n", err)
		os.Exit(1)
	}
}
