// Package dkvpx reads DKVPX records: comma-delimited key=value pairs with
// CSV-style quoting. Each record maps to lib.OrderedMap[string].
//
// Input format examples:
//
//	x=1,y=2,z=3              -> OrderedMap: x->1, y->2, z->3
//	x=1,2,z=3                 -> OrderedMap: x->1, "2"->2, z->3  (implicit keys from 1-up index)
//	"x,y"="a,b,c",z=3         -> OrderedMap: "x,y"->"a,b,c", z->3  (quoting for keys/values with commas)
//
// Quoting uses " with "" as escape. Keys and values may be quoted independently.
// Inside quotes, commas, equals, and newlines are literal.
package dkvpx

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"unicode/utf8"

	"github.com/johnkerl/miller/v6/pkg/lib"
)

// A Reader reads DKVPX records from an io.Reader.
type Reader struct {
	// Comma is the pair delimiter (default ',').
	Comma rune

	// Comment, if not 0, causes lines beginning with this character to be skipped.
	Comment rune

	// TrimLeadingSpace, if true, trims leading space from each key and value.
	TrimLeadingSpace bool

	r *bufio.Reader

	numLine   int
	offset    int64
	rawBuffer []byte
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		Comma: ',',
		r:     bufio.NewReader(r),
	}
}

// Read reads one DKVPX record. It returns nil, io.EOF when there is no more input.
func (r *Reader) Read() (*lib.OrderedMap[string], error) {
	return r.readRecord()
}

// readLine reads the next line (with trailing newline). Normalizes \r\n to \n.
// The result is only valid until the next call to readLine.
func (r *Reader) readLine() ([]byte, error) {
	line, err := r.r.ReadSlice('\n')
	if err == bufio.ErrBufferFull {
		r.rawBuffer = append(r.rawBuffer[:0], line...)
		for err == bufio.ErrBufferFull {
			line, err = r.r.ReadSlice('\n')
			r.rawBuffer = append(r.rawBuffer, line...)
		}
		line = r.rawBuffer
	}
	readSize := len(line)
	if readSize > 0 && err == io.EOF {
		err = nil
		if line[readSize-1] == '\r' {
			line = line[:readSize-1]
		}
	}
	r.numLine++
	r.offset += int64(readSize)
	if n := len(line); n >= 2 && line[n-2] == '\r' && line[n-1] == '\n' {
		line[n-2] = '\n'
		line = line[:n-1]
	}
	return line, err
}

func lengthNL(b []byte) int {
	if len(b) > 0 && b[len(b)-1] == '\n' {
		return 1
	}
	return 0
}

func nextRune(b []byte) rune {
	r, _ := utf8.DecodeRune(b)
	return r
}

func (r *Reader) readRecord() (*lib.OrderedMap[string], error) {
	const quoteLen = len(`"`)

	// Read lines until we have a non-comment line or EOF.
	var line []byte
	var errRead error
	for errRead == nil {
		line, errRead = r.readLine()
		if errRead == io.EOF {
			return nil, io.EOF
		}
		if r.Comment != 0 && len(line) > lengthNL(line) && nextRune(line) == r.Comment {
			continue // skip comment lines
		}
		break
	}

	result := lib.NewOrderedMap[string]()
	var keyBuf, valBuf []byte
	inQuotes := false
	haveKey := false
	pairIndex := 0

	finishRecord := func() {
		if len(keyBuf) > 0 || len(valBuf) > 0 || haveKey {
			key, val := r.finalizePair(keyBuf, valBuf, haveKey, pairIndex)
			result.Put(key, val)
		}
	}

	for {
		if r.TrimLeadingSpace && !inQuotes {
			for len(line) > 0 && line[0] != '\n' && (line[0] == ' ' || line[0] == '\t') {
				line = line[1:]
			}
		}

		if len(line) == 0 {
			if errRead != nil {
				finishRecord()
				return result, errRead
			}
			line, errRead = r.readLine()
			if errRead == io.EOF {
				errRead = nil
				finishRecord()
				return result, nil
			}
			continue
		}

		if inQuotes {
			i := bytes.IndexByte(line, '"')
			if i >= 0 {
				if haveKey {
					valBuf = append(valBuf, line[:i]...)
				} else {
					keyBuf = append(keyBuf, line[:i]...)
				}
				line = line[i+quoteLen:]
				if len(line) > 0 && line[0] == '"' {
					// Escaped quote ""
					if haveKey {
						valBuf = append(valBuf, '"')
					} else {
						keyBuf = append(keyBuf, '"')
					}
					line = line[quoteLen:]
				} else {
					inQuotes = false
				}
			} else {
				contentLen := len(line) - lengthNL(line)
				if haveKey {
					valBuf = append(valBuf, line[:contentLen]...)
				} else {
					keyBuf = append(keyBuf, line[:contentLen]...)
				}
				if contentLen > 0 && len(line) > contentLen {
					if haveKey {
						valBuf = append(valBuf, '\n')
					} else {
						keyBuf = append(keyBuf, '\n')
					}
				}
				line = nil
			}
			continue
		}

		// Not in quotes: process one rune at a time
		rn, rnLen := utf8.DecodeRune(line)
		if rn == utf8.RuneError {
			rnLen = 1
		}

		if line[0] == '"' {
			inQuotes = true
			line = line[quoteLen:]
			continue
		}

		if line[0] == '\n' {
			finishRecord()
			return result, errRead
		}

		if rn == r.Comma {
			key, val := r.finalizePair(keyBuf, valBuf, haveKey, pairIndex)
			result.Put(key, val)
			keyBuf = keyBuf[:0]
			valBuf = valBuf[:0]
			haveKey = false
			pairIndex++
			line = line[rnLen:]
			continue
		}

		if rn == '=' && !haveKey {
			haveKey = true
			line = line[rnLen:]
			continue
		}

		// Regular character
		if haveKey {
			valBuf = append(valBuf, line[:rnLen]...)
		} else {
			keyBuf = append(keyBuf, line[:rnLen]...)
		}
		line = line[rnLen:]
	}
}

func (r *Reader) finalizePair(keyBuf, valBuf []byte, haveKey bool, pairIndex int) (key, val string) {
	if haveKey {
		key = string(keyBuf)
		if key == "" {
			key = strconv.Itoa(pairIndex + 1)
		}
		val = string(valBuf)
	} else {
		// No = seen: entire pair is value, key is 1-up positional
		key = strconv.Itoa(pairIndex + 1)
		val = string(keyBuf)
	}
	return key, val
}
