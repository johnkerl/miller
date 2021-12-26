// This is the entry point for the mlr executable.
package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"

	"github.com/johnkerl/miller/internal/pkg/entrypoint"
	"github.com/pkg/profile" // for trace.out
)

func main() {
	// For mlr --time
	startTime := time.Now()

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

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// CPU profiling
	//
	// We do this here, not in the command-line parser, since
	// pprof.StopCPUProfile() needs to be called at the very end of everything.
	// Putting this pprof logic into a go func running in parallel with main,
	// and properly stopping the profile only when main ends via chan-sync,
	// results in a zero-length pprof file.
	//
	// Please see README-profiling.md for more information.

	if len(os.Args) >= 3 && os.Args[1] == "--cpuprofile" {
		profFilename := os.Args[2]
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
		defer fmt.Fprintf(os.Stderr, "CPU profile finished.\ngo tool pprof -http=:8080 %s\n", profFilename)
	}

	if len(os.Args) >= 2 && os.Args[1] == "--traceprofile" {
		defer profile.Start(profile.TraceProfile, profile.ProfilePath(".")).Stop()
		defer fmt.Fprintf(os.Stderr, "go tool trace trace.out\n")
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// This will obtain os.Args and go from there.  All the usual contents of
	// main() are put into this package for ease of testing.
	mainReturn := entrypoint.Main()

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Timing
	//
	// The system 'time' command is built-in, of course but it's nice to have
	// simply wall-time without the real/user/sys distinction. Also, making
	// this a Miller built-in is nice for Windows.
	if mainReturn.PrintElapsedTime {
		endTime := time.Now()
		startNanos := startTime.UnixNano()
		endNanos := endTime.UnixNano()
		seconds := float64(endNanos-startNanos) / 1e9
		fmt.Fprintf(os.Stderr, "%.6f", seconds)
		for _, arg := range os.Args {
			if strings.Contains(arg, " ") || strings.Contains(arg, "\t") {
				fmt.Fprintf(os.Stderr, " '%s'", arg)
			} else {
				fmt.Fprintf(os.Stderr, " %s", arg)
			}
		}
		fmt.Fprintf(os.Stderr, "\n")
	}
}
