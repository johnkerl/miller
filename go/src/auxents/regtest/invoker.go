package regtest

import (
	"bytes"
	"os/exec"
	"strings"

	"mlr/src/platform"
)

// RunMillerCommand runs a string like 'mlr cat foo.dat', with specified mlr
// executable name to be interpolated into the args[0] slot. This allows us to
// compare different versions of Miller using the same test data.
//
// Note the argsString could have left the exe name off entirely, like 'tac
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

	// Insert the desired Miller executable.
	if strings.HasPrefix(argsString, "mlr ") {
		argsString = strings.Replace(argsString, "mlr", millerExe, 1)
	}

	// This is bash -c ... or cmd /c ...
	shellRunArray := platform.GetShellRunArray(argsString)

	cmd := exec.Command(shellRunArray[0], shellRunArray[1:]...)

	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer

	cmd.Stdout = &stdoutBuffer
	cmd.Stderr = &stderrBuffer

	err := cmd.Run()

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
