// Wrapper for os.Open which maps string filename to *os.File, which in turn
// implements io.ReadCloser, and optional in turn wrapping that in a
// gzip/zlib/bunzip2 reader. Shared across record-readers for all the various
// input-file formats Miller supports.

package lib

import (
	"compress/bzip2"
	"compress/gzip"
	"compress/zlib"
	"errors"
	"io"
	"os"
)

type TFileInputEncoding int

const (
	FileInputEncodingCat TFileInputEncoding = iota
	FileInputEncodingBzip2
	FileInputEncodingGzip
	FileInputEncodingZlib
)

func OpenFile(
	filename string,
	encoding TFileInputEncoding,
) (io.ReadCloser, error) {
	handle, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return openHandle(handle, encoding)
}

func OpenStdin(
	encoding TFileInputEncoding,
) (io.ReadCloser, error) {
	return openHandle(os.Stdin, encoding)
}

func openHandle(
	handle *os.File,
	encoding TFileInputEncoding,
) (io.ReadCloser, error) {
	switch encoding {
	case FileInputEncodingCat:
		return handle, nil
		break
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
	InternalCodingErrorIf(true) // should not have been reached
	return nil, errors.New("to make the compiler happy")
}

// ----------------------------------------------------------------
// BZip2ReadCloser remedies the fact that bzip2.NewReader does not implement io.ReadCloser.
type BZip2ReadCloser struct {
	originalHandle *os.File
	bzip2Handle    io.Reader
}

func NewBZip2ReadCloser(handle *os.File) *BZip2ReadCloser {
	return &BZip2ReadCloser{
		originalHandle: handle,
		bzip2Handle:    bzip2.NewReader(handle),
	}
}

func (this *BZip2ReadCloser) Read(p []byte) (n int, err error) {
	return this.bzip2Handle.Read(p)
}

func (this *BZip2ReadCloser) Close() error {
	return this.originalHandle.Close()
}
