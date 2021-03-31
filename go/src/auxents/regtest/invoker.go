package regtest

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"

	shellquote "github.com/kballard/go-shellquote"
)

// RunMillerCommand runs a string like 'mlr cat foo.dat', with specified mlr
// executable name to be interpolated into the args[0] slot. This allows us to
// compare different versions of Miller using the same test data.
//
// Note the argsString could have left the exe name off entirely, like 'cat
// foo.dat', but it's desirable for debugging to have the command-files be
// directly runnable as-is.
func RunMillerCommand(
	millerExe string,
	argsString string,
) (
	stdout string,
	stderr string,
	exitCode int,
	executionError error, // failure to even start the process
) {
	argsString = strings.TrimRight(argsString, "\n")
	argsString = strings.TrimRight(argsString, "\r")

	argsArray, err := shellquote.Split(argsString)
	if err != nil {
		return "", "", -1, err
	}
	// Given file contents 'mlr cat foo.dat' and millerExe 'mlr-previous',
	// contents split to ['mlr', 'cat', 'foo.dat']. We need a non-empty array
	// in order to overwrite args[0] from 'mlr' to 'mlr-previous'.
	if len(argsArray) < 1 {
		return "", "", -1, errors.New(
			"Empty command in regression-test input",
		)
	}

	argsArray = argsArray[1:] // everything after the Miller-executable name.

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
