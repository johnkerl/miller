package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/pprof"

	"miller/src/auxents"
	"miller/src/cli"
	"miller/src/stream"
)

// ----------------------------------------------------------------
func main() {
	runtime.GOMAXPROCS(4) // Seems reasonable these days
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

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// 'mlr repl' or 'mlr lecat' or any other non-miller-per-se toolery which
	// is delivered (for convenience) within the mlr executable. If argv[1] is
	// found then this function will not return.
	auxents.Dispatch(os.Args)

	options, recordTransformers, err := cli.ParseCommandLine(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, os.Args[0], ": ", err)
		os.Exit(1)
	}

	err = stream.Stream(options, recordTransformers)
	if err != nil {
		fmt.Fprintln(os.Stderr, os.Args[0], ": ", err)
		os.Exit(1)
	}
}
