// ================================================================
// Wrapper for os.Open which maps string filename to *os.File, which in turn
// implements io.ReadCloser, and optional in turn wrapping that in a
// gzip/zlib/bunzip2 reader. Shared across record-readers for all the various
// input-file formats (CSV, JSON, XTAB, DKVP, NIDX, PPRINT) which Miller
// supports.
//
// There are two ways of handling compressed data in the Miller Go port:
//
// * A user-specified 'prepipe' command such as 'gunzip', where we popen a
//   process, hand it the filename via '< filename', and read from that pipe;
//
// * An indication to use an in-process encoding reader (gzip or bzip2, etc).
//
// If a prepipe is specified, it is used; else if an encoding is specified, it
// is used; otherwise the file suffix (.bz2, .gz, .z) is consulted; otherwise
// the file is treated as text.
// ================================================================

package lib

import (
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"compress/zlib"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

type TFileInputEncoding int

const (
	FileInputEncodingDefault TFileInputEncoding = iota
	FileInputEncodingBzip2
	FileInputEncodingGzip
	FileInputEncodingZlib
)

// OpenFileForRead: If prepipe is non-empty, popens "{prepipe} < {filename}"
// and returns a handle to that where prepipe is nominally things like
// "gunzip", "cat", etc.  Otherwise, delegates to an in-process reader which
// can natively handle gzip/bzip2/zlib depending on the specified encoding.  If
// the encoding isn't a compression encoding, this ends up being simply
// os.Open.
func OpenFileForRead(
	filename string,
	prepipe string,
	prepipeIsRaw bool,
	encoding TFileInputEncoding, // ignored if prepipe is non-empty
) (io.ReadCloser, error) {
	if prepipe != "" {
		return openPrepipedHandleForRead(filename, prepipe, prepipeIsRaw)
	} else {
		handle, err := PathToHandle(filename)
		if err != nil {
			return nil, err
		}
		return openEncodedHandleForRead(handle, encoding, filename)
	}
}

// PathToHandle maps various back-ends to a stream. As of 2021-07-07, the
// following URI schemes are supported:
// * https://... and http://...
// * file://...
// * plain disk files
func PathToHandle(
	path string,
) (io.ReadCloser, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		resp, err := http.Get(path)
		if err != nil {
			return nil, err
		}
		handle := resp.Body
		return handle, err
	} else if strings.HasPrefix(path, "file://") {
		return os.Open(strings.Replace(path, "file://", "", 1))
	} else {
		return os.Open(path)
	}
}

// OpenStdin: if prepipe is non-empty, popens "{prepipe}" and returns a handle
// to that where prepipe is nominally things like "gunzip", "cat", etc.
// Otherwise, delegates to an in-process reader which can natively handle
// gzip/bzip2/zlib depending on the specified encoding.  If the encoding isn't
// a compression encoding, this ends up being simply os.Stdin.
func OpenStdin(
	prepipe string,
	prepipeIsRaw bool,
	encoding TFileInputEncoding, // ignored if prepipe is non-empty
) (io.ReadCloser, error) {
	if prepipe != "" {
		return openPrepipedHandleForRead("", prepipe, prepipeIsRaw)
	} else {
		return openEncodedHandleForRead(os.Stdin, encoding, "")
	}
}

func openPrepipedHandleForRead(
	filename string,
	prepipe string,
	prepipeIsRaw bool,
) (io.ReadCloser, error) {
	escapedFilename := escapeFileNameForPopen(filename)

	var command string
	if filename == "" { // stdin
		command = prepipe
	} else {
		if prepipeIsRaw {
			command = prepipe + " " + escapedFilename
		} else {
			command = prepipe + " < " + escapedFilename
		}
	}

	return OpenInboundHalfPipe(command)
}

// Avoids shell-injection cases by replacing single-quote with backslash
// single-quote and double-quote with backslack double-quote, then wrapping the
// entire result in initial and final single-quote.
//
// TODO: test on Windows. Maybe needs move to internal/pkg/platform.
func escapeFileNameForPopen(filename string) string {
	var buffer bytes.Buffer
	foundQuote := false
	for _, c := range filename {
		if c == '\'' || c == '"' {
			buffer.WriteRune('\'')
			buffer.WriteRune(c)
			buffer.WriteRune('\'')
		} else {
			buffer.WriteRune(c)
		}
	}
	if foundQuote {
		return "'" + buffer.String() + "'"
	} else {
		return buffer.String()
	}
}

// TODO: comment
func openEncodedHandleForRead(
	handle io.ReadCloser,
	encoding TFileInputEncoding,
	filename string,
) (io.ReadCloser, error) {
	switch encoding {
	case FileInputEncodingBzip2:
		return NewBZip2ReadCloser(handle), nil
		break
	case FileInputEncodingGzip:
		return gzip.NewReader(handle)
		break
	case FileInputEncodingZlib:
		return zlib.NewReader(handle)
		break
	}

	InternalCodingErrorIf(encoding != FileInputEncodingDefault)

	if strings.HasSuffix(filename, ".bz2") {
		return NewBZip2ReadCloser(handle), nil
	}
	if strings.HasSuffix(filename, ".gz") {
		return gzip.NewReader(handle)
	}
	if strings.HasSuffix(filename, ".z") {
		return zlib.NewReader(handle)
	}

	// Pass along os.Stdin or os.Open(filename)
	return handle, nil
}

// ----------------------------------------------------------------
// BZip2ReadCloser remedies the fact that bzip2.NewReader does not implement io.ReadCloser.
type BZip2ReadCloser struct {
	originalHandle io.ReadCloser
	bzip2Handle    io.Reader
}

func NewBZip2ReadCloser(handle io.ReadCloser) *BZip2ReadCloser {
	return &BZip2ReadCloser{
		originalHandle: handle,
		bzip2Handle:    bzip2.NewReader(handle),
	}
}

func (rc *BZip2ReadCloser) Read(p []byte) (n int, err error) {
	return rc.bzip2Handle.Read(p)
}

func (rc *BZip2ReadCloser) Close() error {
	return rc.originalHandle.Close()
}

// ----------------------------------------------------------------

// IsEOF handles the following problem: reading past end of files opened with
// os.Open returns the error which is io.EOF. Reading past close of pipes
// opened with popen (e.g.  Miller's prepipe, where the file isn't 'foo.dat'
// but rather the process 'gunzip < foo.dat |') returns not io.EOF but an error
// with 'file already closed' within it. See also
// https://stackoverflow.com/questions/47486128/why-does-io-pipe-continue-to-block-even-when-eof-is-reached
func IsEOF(err error) bool {
	if err == nil {
		return false
	} else if err == io.EOF {
		return true
	} else if strings.Contains(err.Error(), "file already closed") {
		return true
	} else {
		return false
	}
}

// ----------------------------------------------------------------
// Functions for in-place mode

// IsUpdateableInPlace tells if we can use the input with mlr -I: not for URLs,
// and not for prepipe commands (which we don't presume to know how to invert
// for output).
func IsUpdateableInPlace(
	filename string,
	prepipe string,
) error {
	if strings.HasPrefix(filename, "http://") ||
		strings.HasPrefix(filename, "https://") ||
		strings.HasPrefix(filename, "file://") {
		return errors.New("http://, https://, and file:// URLs are not updateable in place.")
	}
	if prepipe != "" {
		return errors.New("input with --prepipe or --prepipex is not updateable in place.")
	}
	return nil
}

// FindInputEncoding determines the input encoding (compression), whether from
// a flag like --gzin, or from filename suffix like ".gz".  If the user did
// --gzin on the command line, TFileInputEncoding will be
// FileInputEncodingGzip.  If they didn't, but the filename ends in ".gz", then
// we auto-infer FileInputEncodingGzip.  Either way, this function tells if we
// will be using in-process decompression within the file-format-specific
// record reader.
func FindInputEncoding(
	filename string,
	inputFileInputEncoding TFileInputEncoding,
) TFileInputEncoding {
	if inputFileInputEncoding != FileInputEncodingDefault {
		return inputFileInputEncoding
	}
	if strings.HasSuffix(filename, ".bz2") {
		return FileInputEncodingBzip2
	}
	if strings.HasSuffix(filename, ".gz") {
		return FileInputEncodingGzip
	}
	if strings.HasSuffix(filename, ".z") {
		return FileInputEncodingZlib
	}
	return FileInputEncodingDefault
}

// WrapOutputHandle wraps a file-write handle with a decompressor.  The first
// return value is the wrapped handle. The second is true if the returned
// handle needs to be closed separately from the original.  The third is for
// in-process compression we can't undo: namely, as of September 2021 the gzip
// and zlib libraries support write-closers, but the bzip2 library does not.
func WrapOutputHandle(
	fileWriteHandle io.WriteCloser,
	inputFileEncoding TFileInputEncoding,
) (io.WriteCloser, bool, error) {
	switch inputFileEncoding {
	case FileInputEncodingBzip2:
		return fileWriteHandle, false, errors.New("bzip2 is not currently supported for in-place mode.")
	case FileInputEncodingGzip:
		return gzip.NewWriter(fileWriteHandle), true, nil
	case FileInputEncodingZlib:
		return zlib.NewWriter(fileWriteHandle), true, nil
	default:
		return fileWriteHandle, false, nil
	}
}
