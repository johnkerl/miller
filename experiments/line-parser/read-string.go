package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	filename := os.Args[1]
	handle, err := os.Open(filename)
	if err != nil {
		fmt.Println("OERR", err)
		os.Exit(1)
	}
	defer handle.Close()

	lineReader := bufio.NewReader(handle)

	eof := false
	for !eof {

		line, err := lineReader.ReadString('\n') // TODO: auto-detect
		if err != nil {
			if line != "" {
				fmt.Println(line)
			}
			break
		}
		fmt.Print(line)
	}
}
