package auxents

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// ================================================================
func termcvtUsage(verbName string, o *os.File, exitCode int) {
	fmt.Fprintf(o, "Usage: mlr %s [option] {zero or more file names}\n", verbName)
	fmt.Fprintf(o, "Option (exactly one is required):\n")
	fmt.Fprintf(o, "--cr2crlf\n")
	fmt.Fprintf(o, "--lf2crlf\n")
	fmt.Fprintf(o, "--crlf2cr\n")
	fmt.Fprintf(o, "--crlf2lf\n")
	fmt.Fprintf(o, "--cr2lf\n")
	fmt.Fprintf(o, "--lf2cr\n")
	fmt.Fprintf(o, "-I in-place processing (default is to write to stdout)\n")
	fmt.Fprintf(o, "-h or --help: print this message\n")
	fmt.Fprintf(o, "Zero file names means read from standard input.\n")
	fmt.Fprintf(o, "Output is always to standard output; files are not written in-place.\n")
	os.Exit(exitCode)
}

func termcvtMain(args []string) int {
	inputTerminator := "\n"
	outputTerminator := "\n"
	doInPlace := false

	// 'mlr' and 'termcvt' are already argv[0] and argv[1].
	verb := args[1]
	args = args[2:]
	if len(args) < 1 {
		termcvtUsage(verb, os.Stderr, 1)
	}

	for len(args) >= 1 {
		opt := args[0]
		if opt[0] != '-' {
			break
		}
		args = args[1:]

		if opt == "-h" || opt == "--help" {
			termcvtUsage(verb, os.Stdout, 0)
		} else if opt == "-I" {
			doInPlace = true
		} else if opt == "--cr2crlf" {
			inputTerminator = "\r"
			outputTerminator = "\r\n"
		} else if opt == "--lf2crlf" {
			inputTerminator = "\n"
			outputTerminator = "\r\n"
		} else if opt == "--crlf2cr" {
			inputTerminator = "\r\n"
			outputTerminator = "\r"
		} else if opt == "--lf2cr" {
			inputTerminator = "\n"
			outputTerminator = "\r"
		} else if opt == "--crlf2lf" {
			inputTerminator = "\r\n"
			outputTerminator = "\n"
		} else if opt == "--cr2lf" {
			inputTerminator = "\r"
			outputTerminator = "\n"
		} else {
			termcvtUsage(verb, os.Stderr, 1)
		}
	}

	if len(args) == 0 {
		termcvtFile(os.Stdin, os.Stdout, inputTerminator, outputTerminator)

	} else if doInPlace {
		for _, filename := range args {
			// TODO: make re-entrant via long-random suffix
			suffix := "-termcvt-temp"
			tempname := filename + suffix

			istream, err := os.Open(filename)
			if err != nil {
				// TODO: "mlr"
				fmt.Fprintln(os.Stderr, "mlr termcvt:", err)
				os.Exit(1)
			}

			ostream, err := os.Open(tempname)
			if err != nil {
				// TODO: "mlr"
				fmt.Fprintln(os.Stderr, "mlr termcvt:", err)
				os.Exit(1)
			}

			termcvtFile(istream, ostream, inputTerminator, outputTerminator)

			istream.Close()
			// TODO: check return status
			ostream.Close()

			err = os.Rename(tempname, filename)
			if err != nil {
				// TODO: "mlr"
				fmt.Fprintln(os.Stderr, "mlr termcvt:", err)
				os.Exit(1)
			}
		}

	} else {
		for _, filename := range args {

			istream, err := os.Open(filename)
			if err != nil {
				// TODO: "mlr"
				fmt.Fprintln(os.Stderr, "mlr termcvt:", err)
				os.Exit(1)
			}

			termcvtFile(istream, os.Stdout, inputTerminator, outputTerminator)

			istream.Close()
		}
	}
	return 0
}

func termcvtFile(istream *os.File, ostream *os.File, inputTerminator string, outputTerminator string) {
	lineReader := bufio.NewReader(istream)
	inputTerminatorBytes := []byte(inputTerminator[len(inputTerminator)-1:])[0] // bufio.Reader.ReadString takes char not string delimiter :(

	for {
		line, err := lineReader.ReadString(inputTerminatorBytes)
		if err == io.EOF {
			break
		}

		if err != nil {
			// TODO: "mlr"
			fmt.Fprintln(os.Stderr, "mlr termcvt:", err)
			os.Exit(1)
		}

		// This is how to do a chomp:
		line = strings.TrimRight(line, inputTerminator)
		ostream.Write([]byte(line + outputTerminator))
	}
}
