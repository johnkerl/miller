// ================================================================
// Little side-programs for hex, unhex, LF <-> CR/LF, etc which are delivered
// within the mlr exectuable.
// ================================================================

package auxents

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// ----------------------------------------------------------------
type tAuxMain func(args []string) int
type tAuxUsage func(verbName string, ostream *os.File, exitCode int)

type tAuxLookupEntry struct {
	name  string
	main  tAuxMain
	usage tAuxUsage
}

// We get a Golang "initialization loop" if this is defined statically. So, we
// use a "package init" function.
var _AUX_LOOKUP_TABLE = []tAuxLookupEntry{}

func init() {
	_AUX_LOOKUP_TABLE = []tAuxLookupEntry{
		{"aux-list", auxListMain, auxListUsage},
		{"lecat", lecatMain, lecatUsage},
		{"termcvt", termcvtMain, termcvtUsage},
		{"hex", hexMain, hexUsage},
		{"unhex", unhexMain, unhexUsage},
	}
}

// ================================================================
func Dispatch(args []string) {
	if len(args) < 2 {
		return
	}
	verb := args[1]

	for _, entry := range _AUX_LOOKUP_TABLE {
		if verb == entry.name {
			os.Exit(entry.main(args))
		}
	}

	// Else, return control to mlr.go for the rest of Miller.
}

// ================================================================
func auxListUsage(verbName string, ostream *os.File, exitCode int) {
	fmt.Fprintf(ostream, "Usage: %s %s [options]\n", os.Args[0], verbName)
	fmt.Fprintf(ostream, "Options:\n")
	fmt.Fprintf(ostream, "-h or --help: print this message\n")
	os.Exit(exitCode)
}

func auxListMain(args []string) int {
	ShowAuxEntries(os.Stdout)
	return 0
}

// This symbol is exported for 'mlr --help'.
func ShowAuxEntries(ostream *os.File) {
	fmt.Fprintf(ostream, "Available subcommands:\n")
	for _, entry := range _AUX_LOOKUP_TABLE {
		fmt.Fprintf(ostream, "  %s\n", entry.name)
	}

	fmt.Fprintf(
		ostream,
		"For more information, please invoke %s {subcommand} --help.\n",
		os.Args[0],
	)
}

// ================================================================
func lecatUsage(verbName string, ostream *os.File, exitCode int) {
	fmt.Fprintf(ostream, "Usage: %s %s [options] {zero or more file names}\n", os.Args[0], verbName)
	fmt.Fprintf(ostream, "Simply echoes input, but flags CR characters in red and LF characters in green.\n")
	fmt.Fprintf(ostream, "If zero file names are supplied, standard input is read.\n")
	fmt.Fprintf(ostream, "Options:\n")
	fmt.Fprintf(ostream, "--mono: don't try to colorize the output\n")
	fmt.Fprintf(ostream, "-h or --help: print this message\n")
	os.Exit(exitCode)
}

func lecatMain(args []string) int {
	doColor := true

	// 'mlr' and 'hex' are already argv[0] and argv[1].
	verb := args[1]
	args = args[2:]
	if len(args) >= 1 {
		if args[0] == "-h" || args[0] == "--help" {
			hexUsage(verb, os.Stdout, 0)
		}

		if args[0][0] == '-' {
			if args[0] == "--mono" {
				doColor = false
				args = args[1:]
			} else {
				fmt.Fprintf(os.Stderr, "%s %s: unrecognized option \"%s\".\n",
					os.Args[0], verb, args[0],
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
				// TODO: os.Args[0]
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

// ================================================================
func termcvtUsage(verbName string, ostream *os.File, exitCode int) {
	fmt.Fprintf(ostream, "Usage: %s %s [option] {zero or more file names}\n", os.Args[0], verbName)
	fmt.Fprintf(ostream, "Option (exactly one is required):\n")
	fmt.Fprintf(ostream, "--cr2crlf\n")
	fmt.Fprintf(ostream, "--lf2crlf\n")
	fmt.Fprintf(ostream, "--crlf2cr\n")
	fmt.Fprintf(ostream, "--crlf2lf\n")
	fmt.Fprintf(ostream, "--cr2lf\n")
	fmt.Fprintf(ostream, "--lf2cr\n")
	fmt.Fprintf(ostream, "-I in-place processing (default is to write to stdout)\n")
	fmt.Fprintf(ostream, "-h or --help: print this message\n")
	fmt.Fprintf(ostream, "Zero file names means read from standard input.\n")
	fmt.Fprintf(ostream, "Output is always to standard output; files are not written in-place.\n")
	os.Exit(exitCode)
}

func termcvtMain(args []string) int {
	inTerm := "\n"
	outTerm := "\n"
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
			inTerm = "\r"
			outTerm = "\r\n"
		} else if opt == "--lf2crlf" {
			inTerm = "\n"
			outTerm = "\r\n"
		} else if opt == "--crlf2cr" {
			inTerm = "\r\n"
			outTerm = "\r"
		} else if opt == "--lf2cr" {
			inTerm = "\n"
			outTerm = "\r"
		} else if opt == "--crlf2lf" {
			inTerm = "\r\n"
			outTerm = "\n"
		} else if opt == "--cr2lf" {
			inTerm = "\r"
			outTerm = "\n"
		} else {
			termcvtUsage(verb, os.Stderr, 1)
		}
	}

	if len(args) == 0 {
		termcvtFile(os.Stdin, os.Stdout, inTerm, outTerm)

	} else if doInPlace {
		for _, filename := range args {
			// TODO: make re-entrant via long-random suffix
			suffix := "-termcvt-temp"
			tempname := filename + suffix

			istream, err := os.Open(filename)
			if err != nil {
				// TODO: os.Args[0]
				fmt.Fprintln(os.Stderr, "mlr termcvt:", err)
				os.Exit(1)
			}

			ostream, err := os.Open(tempname)
			if err != nil {
				// TODO: os.Args[0]
				fmt.Fprintln(os.Stderr, "mlr termcvt:", err)
				os.Exit(1)
			}

			termcvtFile(istream, ostream, inTerm, outTerm)

			istream.Close()
			// TODO: check return status
			ostream.Close()

			err = os.Rename(tempname, filename)
			if err != nil {
				// TODO: os.Args[0]
				fmt.Fprintln(os.Stderr, "mlr termcvt:", err)
				os.Exit(1)
			}
		}

	} else {
		for _, filename := range args {

			istream, err := os.Open(filename)
			if err != nil {
				// TODO: os.Args[0]
				fmt.Fprintln(os.Stderr, "mlr termcvt:", err)
				os.Exit(1)
			}

			termcvtFile(istream, os.Stdout, inTerm, outTerm)

			istream.Close()
		}
	}
	return 0
}

func termcvtFile(istream *os.File, ostream *os.File, inTerm string, outTerm string) {
	lineReader := bufio.NewReader(istream)
	inTermFinal := []byte(inTerm[len(inTerm)-1:])[0] // bufio.Reader.ReadString takes char not string delimiter :(

	for {
		line, err := lineReader.ReadString(inTermFinal)
		if err == io.EOF {
			break
		}

		if err != nil {
			// TODO: os.Args[0]
			fmt.Fprintln(os.Stderr, "mlr termcvt:", err)
			os.Exit(1)
		}

		// This is how to do a chomp:
		line = strings.TrimRight(line, inTerm)
		ostream.Write([]byte(line + outTerm))
	}
}

// ================================================================
// Copyright (c) 1998 John Kerl.
// ================================================================
// This is a simple hex dump with hex offsets to the left, hex data in the
// middle, and ASCII at the right.  This is a subset of the functionality of
// Unix od; I wrote it in my NT days.
//
// Example:
//
// $ d2h $(jot 0 128) | unhex | hex
// 00000000: 00 01 02 03  04 05 06 07  08 09 0a 0b  0c 0d 0e 0f |................|
// 00000010: 10 11 12 13  14 15 16 17  18 19 1a 1b  1c 1d 1e 1f |................|
// 00000020: 20 21 22 23  24 25 26 27  28 29 2a 2b  2c 2d 2e 2f | !"#$%&'()*+,-./|
// 00000030: 30 31 32 33  34 35 36 37  38 39 3a 3b  3c 3d 3e 3f |0123456789:;<=>?|
// 00000040: 40 41 42 43  44 45 46 47  48 49 4a 4b  4c 4d 4e 4f |@ABCDEFGHIJKLMNO|
// 00000050: 50 51 52 53  54 55 56 57  58 59 5a 5b  5c 5d 5e 5f |PQRSTUVWXYZ[\]^_|
// 00000060: 60 61 62 63  64 65 66 67  68 69 6a 6b  6c 6d 6e 6f |`abcdefghijklmno|
// 00000070: 70 71 72 73  74 75 76 77  78 79 7a 7b  7c 7d 7e 7f |pqrstuvwxyz{|}~.|
// ================================================================

// ----------------------------------------------------------------
func hexUsage(verbName string, ostream *os.File, exitCode int) {
	fmt.Fprintf(ostream, "Usage: %s %s [options] {zero or more file names}\n", os.Args[0], verbName)
	fmt.Fprintf(ostream, "Simple hex-dump.\n")
	fmt.Fprintf(ostream, "If zero file names are supplied, standard input is read.\n")
	fmt.Fprintf(ostream, "Options:\n")
	fmt.Fprintf(ostream, "-r: print only raw hex without leading offset indicators or trailing ASCII dump.\n")
	fmt.Fprintf(ostream, "-h or --help: print this message\n")
	os.Exit(exitCode)
}

func hexMain(args []string) int {
	doRaw := false

	// 'mlr' and 'hex' are already argv[0] and argv[1].
	verb := args[1]
	args = args[2:]
	if len(args) >= 1 {
		if args[0] == "-r" {
			doRaw = true
			args = args[1:]
		} else if args[0] == "-h" || args[0] == "--help" {
			hexUsage(verb, os.Stdout, 0)
		}
	}

	if len(args) == 0 {
		hexDumpFile(os.Stdin, os.Stdout, doRaw)
	} else {
		for _, filename := range args {
			// Print filename if there is more than one file, unless raw output
			// was requested.
			if !doRaw && len(args) > 1 {
				fmt.Printf("%s:\n", filename)
			}

			istream, err := os.Open(filename)
			if err != nil {
				// TODO: os.Args[0]
				fmt.Fprintln(os.Stderr, "mlr hex:", err)
				os.Exit(1)
			}

			hexDumpFile(istream, os.Stdout, doRaw)

			istream.Close()
			if !doRaw && len(args) > 1 {
				fmt.Println()
			}
		}
	}

	return 0
}

func hexDumpFile(istream *os.File, ostream *os.File, doRaw bool) {
	const bytesPerClump = 4
	const clumpsPerLine = 4
	const bufferSize = bytesPerClump * clumpsPerLine

	buffer := make([]byte, bufferSize)
	eof := false
	offset := 0

	for !eof {
		numBytesRead, err := io.ReadFull(istream, buffer)
		if err == io.EOF {
			err = nil
			eof = true
			break
		}

		// io.ErrUnexpectedEOF is the normal case when the file size isn't an
		// exact multiple of our buffer size.
		if err != nil && err != io.ErrUnexpectedEOF {
			fmt.Fprintln(os.Stderr, "mlr hex:", err)
			os.Exit(1)
		}

		// Print offset "pre" part
		if !doRaw {
			fmt.Printf("%08x: ", offset)
		}

		// Print hex payload
		for i := 0; i < bufferSize; i++ {
			if i < numBytesRead {
				fmt.Printf("%02x ", buffer[i])
			} else {
				fmt.Printf("   ")
			}
			if (i % bytesPerClump) == (bytesPerClump - 1) {
				if (i > 0) && (i < bufferSize-1) {
					fmt.Printf(" ")
				}
			}
		}

		// Print ASCII-dump "post" part
		if !doRaw {
			fmt.Printf("|")

			for i := 0; i < numBytesRead; i++ {
				if buffer[i] >= 0x20 && buffer[i] <= 0x7e {
					fmt.Printf("%c", buffer[i])
				} else {
					fmt.Printf(".")
				}
			}
			fmt.Printf("|")
		}

		// Print line end
		fmt.Printf("\n")

		offset += numBytesRead

	}
}

// ================================================================
func unhexUsage(verbName string, ostream *os.File, exitCode int) {
	fmt.Fprintf(ostream, "Usage: %s %s [option] {zero or more file names}\n", os.Args[0], verbName)
	fmt.Fprintf(ostream, "Options:\n")
	fmt.Fprintf(ostream, "-h or --help: print this message\n")
	fmt.Fprintf(ostream, "Zero file names means read from standard input.\n")
	fmt.Fprintf(ostream, "Output is always to standard output; files are not written in-place.\n")
	os.Exit(exitCode)
}

func unhexMain(args []string) int {
	// 'mlr' and 'hex' are already argv[0] and argv[1].
	verb := args[1]
	args = args[2:]

	if len(args) >= 1 {
		if args[0] == "-h" || args[0] == "--help" {
			hexUsage(verb, os.Stdout, 0)
		}
	}

	if len(args) == 0 {
		unhexFile(os.Stdin, os.Stdout)
	} else {
		for _, filename := range args {
			istream, err := os.Open(filename)
			if err != nil {
				fmt.Fprintln(os.Stderr, "mlr unhex:", err)
				os.Exit(1)
			}
			unhexFile(istream, os.Stdout)
			istream.Close()
		}
	}

	return 0
}

func unhexFile(istream *os.File, ostream *os.File) {
	// Key insight is os.File implements io.Reader
	lineReader := bufio.NewReader(istream)

	var scanValue int
	byteArray := make([]byte, 1)

	re, err := regexp.Compile("\\s+")
	if err != nil {
		fmt.Fprintln(os.Stderr, "mlr unhex: internal coding error detected.")
		os.Exit(1)
	}

	eof := false
	for !eof {
		line, err := lineReader.ReadString('\n') // TODO: auto-detect
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "mlr unhex:", err)
			os.Exit(1)
		}

		// This is how to do a chomp:
		line = strings.TrimRight(line, "\n")

		// Ignore "" which can happen on empty lines
		fields := re.Split(line, -1)
		for _, field := range fields {
			if field != "" {
				n, err := fmt.Sscanf(field, "%x", &scanValue)
				if err != nil {
					fmt.Fprintln(os.Stderr, "mlr unhex:", err)
					os.Exit(1)
				}
				if n != 1 {
					fmt.Fprintln(os.Stderr, "mlr unhex: internal coding error")
					os.Exit(1)
				}
				byteArray[0] = byte(scanValue)
				ostream.Write(byteArray)
			}
		}
	}
}
