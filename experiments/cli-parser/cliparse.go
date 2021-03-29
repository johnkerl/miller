// go mod init cliparse.go
// go get golang.org/x/sys
// go get github.com/kballard/go-shellquote

package main

import (
	"fmt"
	"os"

	"cliparse/arch"
)

func main() {
	fmt.Println("-- os.Args:")
	for i, arg := range os.Args {
		fmt.Printf("args[%d] <<%s>>\n", i, arg)
	}
	fmt.Println()

	fmt.Println("-- canonical args:")
	args := arch.GetMainArgs()

	for i, arg := range args {
		fmt.Printf("args[%d] <<%s>>\n", i, arg)
	}
}
