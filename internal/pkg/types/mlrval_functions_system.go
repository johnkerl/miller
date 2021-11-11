package types

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"mlr/internal/pkg/platform"
	"mlr/internal/pkg/version"
)

func BIF_version() *Mlrval {
	return MlrvalFromString(version.STRING)
}

func BIF_os() *Mlrval {
	return MlrvalFromString(runtime.GOOS)
}

func BIF_hostname() *Mlrval {
	hostname, err := os.Hostname()
	if err != nil {
		return MLRVAL_ERROR
	} else {
		return MlrvalFromString(hostname)
	}
}

func BIF_system(input1 *Mlrval) *Mlrval {
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

	return MlrvalFromString(outputString)
}
