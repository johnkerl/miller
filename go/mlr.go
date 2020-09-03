package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"miller/cli"
	"miller/stream"
)

// ----------------------------------------------------------------
func main() {
	runtime.GOMAXPROCS(4) // Seems reasonable these days

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// CPU profiling
	//
	// We do this here not in the command-line parser since
	// pprof.StopCPUProfile() needs to be called at the very end of everything.
	// Putting this pprof logic into a go func running in parallel with main,
	// and properly stopping the profile only when main ends via chan-sync,
	// results in a zero-length pprof file.
	//
	// Usage:
	// * mlr --cpuprofile cpu.pprof put -f example.mlr then nothing ~/tmp/huge > /dev/null
	// * go tool pprof mlr cpu.pprof
	//   top10
	// * go tool pprof --pdf mlr cpu.pprof > mlr-call-graph.pdf

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
		//	c := make(chan os.Signal)
		//	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		//	go func() {
		//		<-c
		//		pprof.StopCPUProfile()
		//		os.Exit(0)
		//	}()
	}
	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

	// Start of Miller main per se
	options, recordMappers, err := cli.ParseCommandLine(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, os.Args[0], ": ", err)
		os.Exit(1)
	}

	if options.FileNames != nil { // nil if mlr -n
		err = stream.Stream(options, recordMappers)
		if err != nil {
			fmt.Fprintln(os.Stderr, os.Args[0], ": ", err)
			os.Exit(1)
		}
	}
}
