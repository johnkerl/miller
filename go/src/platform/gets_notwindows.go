// +build !windows

package platform

import (
	"os"
)

func GetArgs() []string {
	return os.Args
}
