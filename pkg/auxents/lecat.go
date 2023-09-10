package auxents

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// ================================================================
func lecatUsage(verbName string, o *os.File, exitCode int) {
	fmt.Fprintf(o, "Usage: mlr %s [options] {zero or more file names}\n", verbName)
	fmt.Fprintf(o, "Simply echoes input, but flags CR characters in red and LF characters in green.\n")
	fmt.Fprintf(o, "If zero file names are supplied, standard input is read.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "--mono: don't try to colorize the output\n")
	fmt.Fprintf(o, "-h or --help: print this message\n")
	os.Exit(exitCode)
}

func lecatMain(args []string) int {
	doColor := true

	// 'mlr' and 'hex' are already argv[0] and argv[1].
	verb := args[1]
	args = args[2:]
	if len(args) >= 1 {
		if args[0] == "-h" || args[0] == "--help" {
			lecatUsage(verb, os.Stdout, 0)
		}

		if args[0][0] == '-' {
			if args[0] == "--mono" {
				doColor = false
				args = args[1:]
			} else {
				fmt.Fprintf(os.Stderr, "mlr %s: unrecognized option \"%s\".\n",
					verb, args[0],
				)
				os.Exit(1)
			}
		}
	}

	if len(args) == 0 {
		lecatFile(os.Stdin, doColor)
	} else {
		for _, filename := range args {

			istream, err := os.Open(filename)
			if err != nil {
				fmt.Fprintln(os.Stderr, "mlr lecat:", err)
				os.Exit(1)
			}

			lecatFile(istream, doColor)

			istream.Close()
		}
	}
	return 0
}

func lecatFile(istream *os.File, doColor bool) {
	reader := bufio.NewReader(istream)
	for {
		c, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		if c == '\r' {
			if doColor {
				fmt.Printf("\033[31;01m") // xterm red
			}
			fmt.Printf("[CR]")
			if doColor {
				fmt.Printf("\033[0m")
			}
		} else if c == '\n' {
			if doColor {
				fmt.Printf("\033[32;01m") // xterm green
			}
			fmt.Printf("[LF]\n")
			if doColor {
				fmt.Printf("\033[0m")
			}
		} else {
			fmt.Printf("%c", c)
		}
	}
}
