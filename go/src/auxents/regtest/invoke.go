package regtest

import (
	"bytes"
	"errors"
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

	argsArray, err := shellquote.Split(argsString)
	if err != nil {
		return "", "", -1, err
	}
	//  TODO: comment
	if len(argsArray) < 1 {
		return "", "", -1, errors.New(
			"Empty command in regression-test input",
		)
	}
	argsArray = argsArray[1:] // TODO: comment

	cmd := exec.Command(millerExe, argsArray...)

	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer

	cmd.Stdout = &stdoutBuffer
	cmd.Stderr = &stderrBuffer

	err = cmd.Run()

	exitCode = 0
	stdout = stdoutBuffer.String()
	stderr = stderrBuffer.String()

	if err != nil {
		exitCode = 1
		exitError, ok := err.(*exec.ExitError)
		if ok {
			exitCode = exitError.ExitCode()
		}
	}

	return stdout, stderr, exitCode, nil
}
