// ================================================================
// Wraps 'sh -c foo bar' or 'cmd /c foo bar', nominally for regression-testing.
// ================================================================

// +build !windows

package platform

func GetShellRunArray(command string) []string {
	return []string{"/bin/sh", "-c", command}
}
