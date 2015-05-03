package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

// ----------------------------------------------------------------
func main() {
	args := os.Args[1:]
	includeFields := []string {"a", "x"};

	ok := true
	if len(args) == 0 {
		ok = handle("-", includeFields) && ok
	} else {
		for _, arg := range args {
			ok = handle(arg, includeFields) && ok
		}
	}
	if ok {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

// ----------------------------------------------------------------
func handle(fileName string, includeFields []string) (ok bool) {
	inputStream := os.Stdin
	if fileName != "-" {
		var err error
		if inputStream, err = os.Open(fileName); err != nil {
			log.Println(err)
			return false
		}
	}

	reader := bufio.NewReader(inputStream)
	writer := bufio.NewWriter(os.Stdout)
	eof := false

	for !eof {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			err = nil
			eof = true
		} else if err != nil {
			log.Println(err)
			if fileName != "-" {
				inputStream.Close()
			}
			return false
		} else {
			writer.WriteString(line)
		}
	}
	if fileName != "-" {
		inputStream.Close()
	}
	writer.Flush()

	return true
}
