// ================================================================
// Wraps 'bash -c foo bar' or 'cmd /c foo bar', nominally for regression-testing.
// ================================================================

// +build windows

package platform

import (
	"os"
)

func GetShellRunArray(command string) (string, []string) {
	if os.Getenv("MSYSTEM") != "" {
		return "bash", []string{"-c", command}
	} else {
		// TODO: type up startup-latency issue, as well as what syntaxes
		// both/one/neither supports.
		//return "powershell", []string{"/c", command}
		return "cmd", []string{"/c", command}
	}
}
