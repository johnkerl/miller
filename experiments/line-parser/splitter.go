package main

import (
	"bufio"
	"fmt"
	"os"
)

// From https://pkg.go.dev/bufio#SplitFunc:
//
// SplitFunc is the signature of the split function used to tokenize the input.
// The arguments are an initial substring of the remaining unprocessed data and
// a flag, atEOF, that reports whether the Reader has no more data to give. The
// return values are the number of bytes to advance the input and the next
// token to return to the user, if any, plus an error, if any.
//
// Scanning stops if the function returns an error, in which case some of the
// input may be discarded. If that error is ErrFinalToken, scanning stops with
// no error.
//
// Otherwise, the Scanner advances the input. If the token is not nil, the
// Scanner returns it to the user. If the token is nil, the Scanner reads more
// data and continues scanning; if there is no more data--if atEOF was
// true--the Scanner returns. If the data does not yet hold a complete token,
// for instance if it has no newline while scanning lines, a SplitFunc can
// return (0, nil, nil) to signal the Scanner to read more data into the slice
// and try again with a longer slice starting at the same point in the input.
//
// The function is never called with an empty data slice unless atEOF is true.
// If atEOF is true, however, data may be non-empty and, as always, holds
// unprocessed text.

func main() {
	filename := os.Args[1]
	handle, err := os.Open(filename)
	if err != nil {
		fmt.Println("OERR", err)
		os.Exit(1)
	}

	irs := "xy\n"
	irsbytes := []byte(irs)
	irslen := len(irsbytes)

	scanner := bufio.NewScanner(handle)

	// Custom splitter
	splitter := func(
		data []byte,
		atEOF bool,
	) (
		advance int,
		token []byte,
		err error,
	) {
		datalen := len(data)
		// Example:
		// datalen = 10
		// irslen = 3
		// 0123456789 <--- data
		//        xy* <--- IRS
		end := datalen - irslen
		for i := 0; i <= end; i++ {
			if data[i] == irsbytes[0] {
				match := true
				for j := 1; j < irslen; j++ {
					if data[i+j] != irsbytes[j] {
						match = false
						break
					}
				}
				if match {
					return i + irslen, data[:i], nil
				}
			}
		}
		if !atEOF {
			return 0, nil, nil
		}
		// There is one final token to be delivered, which may be the empty string.
		// Returning bufio.ErrFinalToken here tells Scan there are no more tokens after this
		// but does not trigger an error to be returned from Scan itself.
		return 0, data, bufio.ErrFinalToken
	}
	scanner.Split(splitter)

	// Consume input
	atFirst := true
	for scanner.Scan() {
		if atFirst {
			atFirst = false
		} else {
			fmt.Println()
		}
		fmt.Printf("%s\n", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}
}
