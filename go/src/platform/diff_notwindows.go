// ================================================================
// Handling diff or fc for regression-test.
// ================================================================

//go:build !windows
// +build !windows

package platform

func GetDiffRunArray(filename1, filename2 string) []string {
	return []string{"diff", "-u", filename1, filename2}
}
