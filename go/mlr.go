package main

import (
	"fmt"
	"os"
	"runtime"

	"miller/cli"
	"miller/stream"
)

// ----------------------------------------------------------------
// xxx to do: stdout/stderr w/ ternary on exitrc
func usage() {
	// xxx temp grammar
	fmt.Fprintf(os.Stderr, "Usage: %s [options] {ifmt} {mapper} {ofmt} {filenames ...}\n",
		os.Args[0])
	fmt.Fprintf(os.Stderr, "If no file name is given, or if filename is \"-\", stdin is used.\n")
	// stub
	os.Exit(1)
}

// ----------------------------------------------------------------
func main() {
	runtime.GOMAXPROCS(4) // Seems reasonable these days

	options, recordMappers, filenames, err := cli.ParseCommandLine(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, os.Args[0], ": ", err)
		os.Exit(1)
	}

	err = stream.Stream(options, recordMappers, filenames)
	if err != nil {
		fmt.Fprintln(os.Stderr, os.Args[0], ": ", err)
		os.Exit(1)
	}
}
