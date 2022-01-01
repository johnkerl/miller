// +build !windows

package arch

import (
	"os"
)

func GetMainArgs() []string {
	return os.Args
}
