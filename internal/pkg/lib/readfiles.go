// ================================================================
// Routines for loading strings from files. Nominally for the put/filter verbs
// to load DSL strings from .mlr files.
// ================================================================

package lib

import (
	"encoding/csv"
	"io/ioutil"
	"os"
	"strings"
)

// LoadStringsFromFileOrDir calls LoadStringFromFile if path exists and is a
// file, or LoadStringsFromDir if path exists and is a directory.  In the
// former case the extension is ignored; in the latter case it's used as a
// filter on the directory entries.
func LoadStringsFromFileOrDir(path string, extension string) ([]string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		return LoadStringsFromDir(path, extension)
	} else {
		dslString, err := LoadStringFromFile(path)
		if err != nil {
			return nil, err
		} else {
			return []string{dslString}, nil
		}
	}
}

// LoadStringFromFile is just a wrapper around ioutil.ReadFile,
// with a cast from []byte to string.
func LoadStringFromFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// LoadStringsFromDir loads all file contents for files in the given directory
// having the given extension. E.g. LoadStringsFromDir("/u/myfiles", ".mlr")
// will load /u/myfiles/foo.mlr and /u/myfiles/bar.mlr but will skip over
// /u/myfiles/data.csv and /u/myfiles/todo.txt.
func LoadStringsFromDir(dirname string, extension string) ([]string, error) {
	dslStrings := make([]string, 0)

	entries, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	for i := range entries {
		entry := &entries[i]
		name := (*entry).Name()
		if !strings.HasSuffix(name, extension) {
			continue
		}

		path := dirname + "/" + name
		dslString, err := LoadStringFromFile(path)
		if err != nil {
			return nil, err
		}

		dslStrings = append(dslStrings, dslString)
	}

	return dslStrings, nil
}

func ReadCSVHeader(filename string) ([]string, error) {
	handle, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer handle.Close()
	csvReader := csv.NewReader(handle)
	header, err := csvReader.Read()
	if err != nil {
		return nil, err
	}
	return header, nil
}
