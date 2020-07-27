package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
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

			// Line to map
			// Note: needs to have insertion-ordering
			mymap := make(map[string]string)
			fields := strings.Split(line, ",");
			for _, field := range(fields) {
				kvps := strings.SplitN(field, "=", 2)
				mymap[kvps[0]] = kvps[1]
			}

			// Map-to-map transform
			newmap := make(map[string]string)
			for _, includeField := range(includeFields) {
				value, present := mymap[includeField]
				if present {
					newmap[includeField] = value
				}
			}

			// Map to string
			outs := make([]string, len(newmap))
			i := 0
			for k, v := range(newmap) {
				outs[i] = k + "=" + v
				i++
			}

			out := strings.Join(outs, ",")

			// Write to output stream
			//fmt.Println("")
			writer.WriteString(out)
			writer.WriteString("\n")
		}
	}
	if fileName != "-" {
		inputStream.Close()
	}
	writer.Flush()

	return true
}
