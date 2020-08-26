package main

import (
	// System:
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	// Miller:
	"stream"
)

// ----------------------------------------------------------------
// xxx to do: stdout/stderr w/ ternary on exitrc
func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] {filenames ...}\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "If no file name is given, or if filename is \"-\", stdin is used.\n")
	flag.PrintDefaults()
	os.Exit(1)
}

// ----------------------------------------------------------------
func main() {
	// pFoo := flag.Bool("f", false, "Foo")
	cpuprofile := flag.String("cpuprofile", "", "Write CPU profile to `file`")

	flag.Usage = usage
	flag.Parse()
	args := flag.Args()

	// foo := *pFoo

	// go tool pprof mlr foo.prof
	//   top10
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("Could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("Could not start CPU profile: ", err)
		}
	    defer pprof.StopCPUProfile()
	}

	err := stream.Stream(args)
	if err != nil {
		log.Fatal(err)
	}
}
