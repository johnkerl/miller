// ================================================================
// Wraps 'sh -c foo bar' or 'cmd /c foo bar', nominally for regression-testing.
// ================================================================

// +build !windows

package platform

func GetShellRunArray(command string) (string, []string) {
	return "sh", []string{"-c", command}
}
