package main

import (
	// System:
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"strings"
	// Miller:
	//"containers"
	"input"
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

	if len(args) == 0 {
		err := stream("-")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		for _, arg := range args {
			err := stream(arg)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// ----------------------------------------------------------------
func stream(sourceName string) error {
	inputStream := os.Stdin
	if sourceName != "-" {
		var err error
		if inputStream, err = os.Open(sourceName); err != nil {
			return err
		}
	}

	reader := bufio.NewReader(inputStream)

	eof := false

	for !eof {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			err = nil
			eof = true
		} else if err != nil {
			return err
		} else {
			if false {
				fmt.Print(line)
			} else {
				// This is how to do a chomp:
				line = strings.TrimRight(line, "\n")

				// xxx temp
				ifs := ","
				ips := "="
				lrec := input.LrecFromDKVPLine(&line, &ifs, &ips)

				lrec.Print(os.Stdout)
			}
		}
	}

	return nil
}
