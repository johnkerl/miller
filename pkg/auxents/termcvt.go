package auxents

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func termcvtUsage(verbName string, o *os.File) {
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
}

func termcvtMain(args []string) int {
	inputTerminator := "\n"
	outputTerminator := "\n"
	doInPlace := false

	// 'mlr' and 'termcvt' are already argv[0] and argv[1].
	verb := args[1]
	args = args[2:]
	if len(args) < 1 {
		termcvtUsage(verb, os.Stderr)
		return 1
	}

	for len(args) >= 1 {
		opt := args[0]
		if opt[0] != '-' {
			break
		}
		args = args[1:]

		switch opt {
		case "-h", "--help":
			termcvtUsage(verb, os.Stdout)
			return 0
		case "-I":
			doInPlace = true
		case "--cr2crlf":
			inputTerminator = "\r"
			outputTerminator = "\r\n"
		case "--lf2crlf":
			inputTerminator = "\n"
			outputTerminator = "\r\n"
		case "--crlf2cr":
			inputTerminator = "\r\n"
			outputTerminator = "\r"
		case "--lf2cr":
			inputTerminator = "\n"
			outputTerminator = "\r"
		case "--crlf2lf":
			inputTerminator = "\r\n"
			outputTerminator = "\n"
		case "--cr2lf":
			inputTerminator = "\r"
			outputTerminator = "\n"
		default:
			termcvtUsage(verb, os.Stderr)
			return 1
		}
	}

	if len(args) == 0 {
		if err := termcvtFile(os.Stdin, os.Stdout, inputTerminator, outputTerminator); err != nil {
			fmt.Fprintf(os.Stderr, "mlr termcvt: %v\n", err)
			return 1
		}

	} else if doInPlace {
		for _, filename := range args {
			// TODO: make re-entrant via long-random suffix
			suffix := "-termcvt-temp"
			tempname := filename + suffix

			istream, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "mlr termcvt: %v\n", err)
				return 1
			}

			ostream, err := os.Open(tempname)
			if err != nil {
				fmt.Fprintf(os.Stderr, "mlr termcvt: %v\n", err)
				return 1
			}

			if err := termcvtFile(istream, ostream, inputTerminator, outputTerminator); err != nil {
				fmt.Fprintf(os.Stderr, "mlr termcvt: %v\n", err)
				return 1
			}

			_ = istream.Close()
			if err := ostream.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "mlr termcvt: %v\n", err)
				return 1
			}

			err = os.Rename(tempname, filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "mlr termcvt: %v\n", err)
				return 1
			}
		}

	} else {
		for _, filename := range args {

			istream, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "mlr termcvt: %v\n", err)
				return 1
			}

			err = termcvtFile(istream, os.Stdout, inputTerminator, outputTerminator)
			_ = istream.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "mlr termcvt: %v\n", err)
				return 1
			}
		}
	}
	return 0
}

func termcvtFile(istream *os.File, ostream *os.File, inputTerminator string, outputTerminator string) error {
	lineReader := bufio.NewReader(istream)
	inputTerminatorBytes := []byte(inputTerminator[len(inputTerminator)-1:])[0] // bufio.Reader.ReadString takes char not string delimiter :(

	for {
		line, err := lineReader.ReadString(inputTerminatorBytes)
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		// This is how to do a chomp:
		line = strings.TrimRight(line, inputTerminator)
		if _, err := ostream.Write([]byte(line + outputTerminator)); err != nil {
			return err
		}
	}
	return nil
}
