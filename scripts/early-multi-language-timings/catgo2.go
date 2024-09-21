package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// ----------------------------------------------------------------
func main2() {
	args := os.Args[1:]

	ok := true
	if len(args) == 0 {
		ok = handle2("-") && ok
	} else {
		for _, arg := range args {
			ok = handle2(arg) && ok
		}
	}
	if ok {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

// ----------------------------------------------------------------
func handle2(fileName string) (ok bool) {
	inputStream := os.Stdin
	if fileName != "-" {
		var err error
		if inputStream, err = os.Open(fileName); err != nil {
			log.Println(err)
			return false
		}
	}

	scanner := bufio.NewScanner(inputStream)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	if fileName != "-" {
		inputStream.Close()
	}

	return true
}
