// go mod init cliparse.go
// go get golang.org/x/sys

package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows"
)

func main() {
	fmt.Printf("<<%s>>\n", windows.UTF16PtrToString(windows.GetCommandLine()))
	fmt.Println()
	for i, arg := range os.Args {
		fmt.Printf("args[%d] <<%s>>\n", i, arg)
	}
}
