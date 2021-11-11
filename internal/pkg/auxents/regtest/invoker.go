package regtest

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/platform"
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

// RunDiffCommandOnStrings runs either diff or fc (not-Windows / Windows
// respectively) to show differences between actual and expected
// regression-test output.
func RunDiffCommandOnStrings(
	actualOutput string,
	expectedOutput string,
) (
	diffOutput string,
) {
	actualOutputFileName := lib.WriteTempFileOrDie(actualOutput)
	expectedOutputFileName := lib.WriteTempFileOrDie(expectedOutput)
	defer os.Remove(actualOutputFileName)
	defer os.Remove(expectedOutputFileName)

	// This is diff or fc
	diffRunArray := platform.GetDiffRunArray(actualOutputFileName, expectedOutputFileName)

	cmd := exec.Command(diffRunArray[0], diffRunArray[1:]...)

	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer
	cmd.Stdout = &stdoutBuffer
	cmd.Stderr = &stderrBuffer

	// Ignore the error-return since it's likely the fact that diff exits
	// non-zero when files differ at all. Otherwise it's a failure to invoke
	// diff itself, about which we can do little within the regtest.  A diff
	// output is simply something (in addition to printing the actual &
	// expected outputs) to help people debug, and hey, we tried.

	// err := cmd.Run()
	_ = cmd.Run()
	//if err != nil {
	//	fmt.Printf("Error executing %s:\n", strings.Join(diffRunArray, " "))
	//	fmt.Println(err)
	//	fmt.Println(stderrBuffer.String())
	//	os.Exit(1)
	//}

	return stdoutBuffer.String()
}

// RunDiffCommandOnFilenames runs either diff or fc (not-Windows / Windows
// respectively) to show differences between actual and expected
// regression-test output.
func RunDiffCommandOnFilenames(
	actualOutputFileName string,
	expectedOutputFileName string,
) (
	diffOutput string,
) {
	// This is diff or fc
	diffRunArray := platform.GetDiffRunArray(actualOutputFileName, expectedOutputFileName)

	cmd := exec.Command(diffRunArray[0], diffRunArray[1:]...)

	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer
	cmd.Stdout = &stdoutBuffer
	cmd.Stderr = &stderrBuffer

	// Ignore the error-return since it's likely the fact that diff exits
	// non-zero when files differ at all. Otherwise it's a failure to invoke
	// diff itself, about which we can do little within the regtest.  A diff
	// output is simply something (in addition to printing the actual &
	// expected outputs) to help people debug, and hey, we tried.

	// err := cmd.Run()
	_ = cmd.Run()
	//if err != nil {
	//	fmt.Printf("Error executing %s:\n", strings.Join(diffRunArray, " "))
	//	fmt.Println(err)
	//	fmt.Println(stderrBuffer.String())
	//	os.Exit(1)
	//}

	return stdoutBuffer.String()
}
