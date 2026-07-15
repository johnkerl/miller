package auxents

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func lecatUsage(verbName string, o *os.File) {
	fmt.Fprintf(o, "Usage: mlr %s [options] {zero or more file names}\n", verbName)
	fmt.Fprintf(o, "Simply echoes input, but flags CR characters in red and LF characters in green.\n")
	fmt.Fprintf(o, "If zero file names are supplied, standard input is read.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "--mono: don't try to colorize the output\n")
	fmt.Fprintf(o, "-h or --help: print this message\n")
}

func lecatMain(args []string) int {
	doColor := true

	// 'mlr' and 'hex' are already argv[0] and argv[1].
	verb := args[1]
	args = args[2:]
	if len(args) >= 1 {
		if args[0] == "-h" || args[0] == "--help" {
			lecatUsage(verb, os.Stdout)
			return 0
		}

		if args[0][0] == '-' {
			if args[0] == "--mono" {
				doColor = false
				args = args[1:]
			} else {
				fmt.Fprintf(os.Stderr, "mlr %s: unrecognized option \"%s\".\n",
					verb, args[0],
				)
				return 1
			}
		}
	}

	if len(args) == 0 {
		if err := lecatFile(os.Stdin, doColor); err != nil {
			fmt.Fprintf(os.Stderr, "mlr lecat: %v\n", err)
			return 1
		}
	} else {
		for _, filename := range args {

			istream, err := os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "mlr lecat: %v\n", err)
				return 1
			}

			err = lecatFile(istream, doColor)
			_ = istream.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "mlr lecat: %v\n", err)
				return 1
			}
		}
	}
	return 0
}

func lecatFile(istream *os.File, doColor bool) error {
	reader := bufio.NewReader(istream)
	for {
		c, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch c {
		case '\r':
			if doColor {
				fmt.Printf("\033[31;01m") // xterm red
			}
			fmt.Printf("[CR]")
			if doColor {
				fmt.Printf("\033[0m")
			}
		case '\n':
			if doColor {
				fmt.Printf("\033[32;01m") // xterm green
			}
			fmt.Printf("[LF]\n")
			if doColor {
				fmt.Printf("\033[0m")
			}
		default:
			fmt.Printf("%c", c)
		}
	}
	return nil
}
