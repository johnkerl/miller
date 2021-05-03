package types

import (
	"os"
	"runtime"

	"miller/src/version"
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
