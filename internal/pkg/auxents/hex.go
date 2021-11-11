package auxents

import (
	"fmt"
	"io"
	"os"
)

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
func hexUsage(verbName string, o *os.File, exitCode int) {
	fmt.Fprintf(o, "Usage: mlr %s [options] {zero or more file names}\n", verbName)
	fmt.Fprintf(o, "Simple hex-dump.\n")
	fmt.Fprintf(o, "If zero file names are supplied, standard input is read.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-r: print only raw hex without leading offset indicators or trailing ASCII dump.\n")
	fmt.Fprintf(o, "-h or --help: print this message\n")
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
		hexDumpFile(os.Stdin, doRaw)
	} else {
		for _, filename := range args {
			// Print filename if there is more than one file, unless raw output
			// was requested.
			if !doRaw && len(args) > 1 {
				fmt.Printf("%s:\n", filename)
			}

			istream, err := os.Open(filename)
			if err != nil {
				// TODO: "mlr"
				fmt.Fprintln(os.Stderr, "mlr hex:", err)
				os.Exit(1)
			}

			hexDumpFile(istream, doRaw)

			istream.Close()
			if !doRaw && len(args) > 1 {
				fmt.Println()
			}
		}
	}

	return 0
}

func hexDumpFile(istream *os.File, doRaw bool) {
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
		// exact multiple of our buffer size. We'll break the loop after
		// hex-dumping this last, partial fragment.
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

		// Break the loop if this was the last, partial fragment.  If we don't
		// break here we'll go back to the top of the loop and try the next
		// read and get EOF -- which works fine if the input is a file. But if
		// the input is at the terminal, the user will have to control-D twice
		// which will be unsettling.
		if numBytesRead < bufferSize {
			eof = true
			break
		}

		offset += numBytesRead
	}
}
