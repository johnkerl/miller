package main

import (
	// System:
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	// Miller:
	//"containers"
)

// ----------------------------------------------------------------
// xxx to do: stdout/stderr w/ ternary on exitrc
func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] {filenames ...}\n", os.Args[0])
	//fmt.Fprintf(os.Stderr, "If no file name is given, or if filename is \"-\", stdin is used.\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

// ----------------------------------------------------------------
func main() {
	// pFoo := flag.Bool("f", false, "Foo")

	flag.Usage = usage
	flag.Parse()
	args := flag.Args()

	// foo := *pFoo

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
	os.Exit(0)
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
			// This is how to do a chomp:
			line = strings.TrimRight(line, "\n")
		}
	}

	return nil
}
