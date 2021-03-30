package support

import (
	"bytes"
	"fmt"
	"os/exec"

	shellquote "github.com/kballard/go-shellquote"
)

func RunMillerCommand(
	millerExe string,
	argsString string,
) (
	stdout string,
	stderr string,
	exitCode int,
	executionError error, // TODO: comment: failure to even start the process
) {

	// xxx
	millerExe = "../mlr"

	fmt.Printf("%s %s\n", millerExe, argsString)
	argsArray, err := shellquote.Split(argsString)
	if err != nil {
		return "", "", -1, err
	}

	cmd := exec.Command(millerExe, argsArray...)

	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer

	cmd.Stdout = &stdoutBuffer
	cmd.Stderr = &stderrBuffer

	err = cmd.Run()

	exitCode = 0
	stdout = stdoutBuffer.String()
	stderr = stderrBuffer.String()

	fmt.Println("AAA")
	if err != nil {
	fmt.Println("BBB", exitCode)
		exitError, ok := err.(*exec.ExitError)
		if ok {
			exitCode = exitError.ExitCode()
	fmt.Println("CCC", exitCode)
		}
	}

	return stdout, stderr, exitCode, nil
}
