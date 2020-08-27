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
	// xxx temp grammar
	fmt.Fprintf(os.Stderr, "Usage: %s [options] {ifmt} {mapper} {ofmt} {filenames ...}\n",
		os.Args[0])
	fmt.Fprintf(os.Stderr, "If no file name is given, or if filename is \"-\", stdin is used.\n")
	flag.PrintDefaults()
	os.Exit(1)
}

// ----------------------------------------------------------------
func main() {
	cpuprofile := flag.String("cpuprofile", "", "Write CPU profile to `file`")

	flag.Usage = usage
	flag.Parse()
	maybeProfile(cpuprofile)

	args := flag.Args()

	if len(args) < 3 {
		usage()
	}
	inputFormatName := args[0]
	mapperName := args[1]
	outputFormatName := args[2]
	args = args[3:]

	err := stream.Stream(inputFormatName, mapperName, outputFormatName, args)
	if err != nil {
		log.Fatal(err)
	}
}

// ----------------------------------------------------------------
func maybeProfile(cpuprofile *string) {
	// to do: move to method
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

}
