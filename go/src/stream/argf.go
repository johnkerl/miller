package stream

import (
	"io"
	"os"
)

func Argf(filenames []string) (io.Reader, error) {
	if len(filenames) == 0 {
		return os.Stdin, nil
	} else {
		readers := make([]io.Reader, len(filenames))
		for i, filename := range filenames {
			handle, err := os.Open(filename)
			if err == nil {
				readers[i] = handle
			} else {
				return nil, err
			}
		}
		return io.MultiReader(readers...), nil
	}
}
