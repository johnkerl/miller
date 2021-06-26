// ================================================================
// Wraps 'sh -c foo bar' or 'cmd /c foo bar', nominally for regression-testing.
// ================================================================

// +build !windows

package platform

func GetShellRunCommandAndArray(command string) (string, []string) {
	return "/bin/sh", []string{"-c", command}
}

func GetShellRunArray(command string) []string {
	return []string{"/bin/sh", "-c", command}
}
