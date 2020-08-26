package stream

import (
	// System:
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	// Miller:
	"input"
)

func Stream(sourceName string) error {
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
