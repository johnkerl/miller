package main

import (
	"fmt"
	"os"
)

func main() {
	base := "/bin/sh"
	rest := []string{"/bin/sh", "-c", "echo hello"}
	//base := "cmd"
	//rest := []string{"cmd", "/c", "echo hello"}

	var procAttr os.ProcAttr
	procAttr.Files = []*os.File{
		os.Stdin,
		os.Stdout,
		os.Stderr,
	}
	process, err := os.StartProcess(base, rest, &procAttr)
	if err != nil {
		fmt.Println(err)
	} else {
		go process.Wait()
	}
}
