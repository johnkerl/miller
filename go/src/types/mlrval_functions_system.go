package types

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"mlr/src/platform"
	"mlr/src/version"
)

func MlrvalVersion() *Mlrval {
	return MlrvalPointerFromString(version.STRING)
}

func MlrvalOS() *Mlrval {
	return MlrvalPointerFromString(runtime.GOOS)
}

func MlrvalHostname() *Mlrval {
	hostname, err := os.Hostname()
	if err != nil {
		return MLRVAL_ERROR
	} else {
		return MlrvalPointerFromString(hostname)
	}
}

func MlrvalSystem(input1 *Mlrval) *Mlrval {
	if !input1.IsStringOrVoid() {
		return MLRVAL_ERROR
	}
	commandString := input1.printrep

	shellRunArray := platform.GetShellRunArray(commandString)

	outputBytes, err := exec.Command(shellRunArray[0], shellRunArray[1:]...).Output()
	if err != nil {
		return MLRVAL_ERROR
	}
	outputString := strings.TrimRight(string(outputBytes), "\n")

	return MlrvalPointerFromString(outputString)
}
