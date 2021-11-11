// Package main is the entry point for Miller.
package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strconv"

	"mlr/internal/pkg/entrypoint"
)

// main is the entry point for Miller.
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

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// CPU profiling
	//
	// We do this here not in the command-line parser since
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
		defer fmt.Fprintf(os.Stderr, "CPU profile finished.\n")
	}

	// This will obtain os.Args and go from there.  All the usual contents of
	// main() are put into this package for ease of testing.
	entrypoint.Main()
}
