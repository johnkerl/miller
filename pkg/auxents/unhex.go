package auxents

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func unhexUsage(verbName string, o *os.File) {
	fmt.Fprintf(o, "Usage: mlr %s [option] {zero or more file names}\n", verbName)
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h or --help: print this message\n")
	fmt.Fprintf(o, "Zero file names means read from standard input.\n")
	fmt.Fprintf(o, "Output is always to standard output; files are not written in-place.\n")
}

func unhexMain(args []string) int {
	// 'mlr' and 'hex' are already argv[0] and argv[1].
	verb := args[1]
	args = args[2:]

	if len(args) >= 1 {
		if args[0] == "-h" || args[0] == "--help" {
			unhexUsage(verb, os.Stdout)
			return 0
		}
	}

	if len(args) == 0 {
		if err := unhexFile(os.Stdin, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "mlr unhex: %v\n", err)
			return 1
		}
	} else {
		for _, filename := range args {
			istream, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "mlr unhex: %v\n", err)
				return 1
			}
			err = unhexFile(istream, os.Stdout)
			_ = istream.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "mlr unhex: %v\n", err)
				return 1
			}
		}
	}

	return 0
}

func unhexFile(istream *os.File, ostream *os.File) error {
	// Key insight is os.File implements io.Reader
	lineReader := bufio.NewReader(istream)

	var scanValue int
	byteArray := make([]byte, 1)

	re := regexp.MustCompile(`\s+`)

	eof := false
	for !eof {
		line, err := lineReader.ReadString('\n') // TODO: auto-detect
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// This is how to do a chomp:
		line = strings.TrimRight(line, "\n")

		// Ignore "" which can happen on empty lines
		fields := re.Split(line, -1)
		for _, field := range fields {
			if field != "" {
				n, err := fmt.Sscanf(field, "%x", &scanValue)
				if err != nil {
					return err
				}
				if n != 1 {
					return errors.New("internal coding error")
				}
				byteArray[0] = byte(scanValue)
				if _, err := ostream.Write(byteArray); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
