package main

// https://golang.org/pkg/encoding/json

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func main() {
	decoder := json.NewDecoder(os.Stdin)

	var token json.Token
	var err error
	for decoder.More() {
		token, err = decoder.Token()
		if err == io.EOF {
			fmt.Println("EOF")
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(token)
	}
}
