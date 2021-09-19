// ================================================================
// Handling diff or fc for regression-test.
// ================================================================

//go:build windows
// +build windows

package platform

// GetDiffRunArray gets the command for diffing actual/expected output in a
// regression test.
func GetDiffRunArray(filename1, filename2 string) []string {
	return []string{"fc", filename1, filename2}
}
