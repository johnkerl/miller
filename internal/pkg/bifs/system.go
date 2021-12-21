package bifs

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/platform"
	"github.com/johnkerl/miller/internal/pkg/version"
)

func BIF_version() *mlrval.Mlrval {
	return mlrval.FromString(version.STRING)
}

func BIF_os() *mlrval.Mlrval {
	return mlrval.FromString(runtime.GOOS)
}

func BIF_hostname() *mlrval.Mlrval {
	hostname, err := os.Hostname()
	if err != nil {
		return mlrval.ERROR
	} else {
		return mlrval.FromString(hostname)
	}
}

func BIF_system(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsStringOrVoid() {
		return mlrval.ERROR
	}
	commandString := input1.AcquireStringValue()

	shellRunArray := platform.GetShellRunArray(commandString)

	outputBytes, err := exec.Command(shellRunArray[0], shellRunArray[1:]...).Output()
	if err != nil {
		return mlrval.ERROR
	}
	outputString := strings.TrimRight(string(outputBytes), "\n")

	return mlrval.FromString(outputString)
}
