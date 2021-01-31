package lib

import (
	"os"
	"path"
)

func MlrExeName() string {
	return path.Base(os.Args[0])
}
