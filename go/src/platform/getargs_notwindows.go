// ================================================================
// Handling single quotes and double quotes is different on Windows unless
// particular care is taken, which is what this file does.
// ================================================================

//go:build !windows
// +build !windows

package platform

import (
	"os"
)

func GetArgs() []string {
	return os.Args
}
