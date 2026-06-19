package input

import (
	"bytes"
	"io"
	"unsafe"

	csv "github.com/johnkerl/miller/v6/pkg/go-csv"
)

// ZeroCopyCSVReader is a prototype CSV record reader that avoids the
// per-record string allocation made by the stdlib-derived go-csv parser.
//
// The stdlib parser, for each record, (a) accumulates the record's bytes into a
// reused buffer and (b) does `string(buffer)` -- a fresh heap allocation and
// copy -- so the returned field substrings have stable backing. That is one
// allocation (plus the field-slice) per record.
//
// This reader instead reads the input in large, persistent blocks and returns
// field strings that point *directly into the block* via unsafe.String. The
// block stays alive for exactly as long as some field references it (ordinary
// Go GC reachability), so no per-record copy is needed. Field bytes are never
// mutated after being read, so the unsafe string views are sound.
//
// Scope of the prototype:
//   - The fast path handles unquoted records (no '"' in the record) with ','
//     (or the configured single-byte delimiter), LF or CRLF line endings. This
//     is the overwhelmingly common case and is fully zero-copy for the field
//     backing bytes.
//   - Records containing '"' (quoted fields, embedded newlines, escaped quotes)
//     are delegated to the proven go-csv parser over the isolated record bytes,
//     for guaranteed correctness. These allocate, but are comparatively rare.
//   - It is used only when LazyQuotes and TrimLeadingSpace are off (see the
//     gating in processHandle); otherwise the original reader is used.
//
// It exposes Read() []string to be a drop-in for csv.Reader in the scanner.
type ZeroCopyCSVReader struct {
	handle  io.Reader
	comma   byte
	comment byte // 0 if comment handling is disabled

	lazyQuotes       bool
	trimLeadingSpace bool

	buf           []byte // current persistent block
	parseOff      int    // parse cursor within buf
	filled        int    // valid bytes in buf
	eofSeen       bool
	pendingErr    error // non-EOF read error, surfaced after buffered data is drained
	linesConsumed int   // physical input lines consumed so far (for error line numbers)

	blockSize int
}

const zeroCopyCSVBlockSize = 64 * 1024

func NewZeroCopyCSVReader(handle io.Reader, comma byte, comment byte, lazyQuotes, trimLeadingSpace bool) *ZeroCopyCSVReader {
	return &ZeroCopyCSVReader{
		handle:           handle,
		comma:            comma,
		comment:          comment,
		lazyQuotes:       lazyQuotes,
		trimLeadingSpace: trimLeadingSpace,
		blockSize:        zeroCopyCSVBlockSize,
	}
}

// bytesToString returns a string sharing b's backing array (no copy). Safe only
// while b's bytes are immutable for the string's lifetime, which holds here:
// block bytes are written once (by Read from the input) and never modified.
func bytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// fill moves the unparsed tail to the front of a fresh block and reads more
// input after it. A fresh block is allocated (rather than reusing buf) so that
// field strings already handed out -- which point into the old block -- remain
// valid; the old block is freed by GC once its last field is unreferenced.
func (r *ZeroCopyCSVReader) fill() {
	tailLen := r.filled - r.parseOff
	newSize := r.blockSize
	if tailLen*2 > newSize {
		newSize = tailLen * 2
	}
	newbuf := make([]byte, newSize)
	copy(newbuf, r.buf[r.parseOff:r.filled])
	r.buf = newbuf
	r.parseOff = 0
	// Read until we make progress (n > 0) or hit EOF. An io.Reader is permitted
	// to return (0, nil); looping here guarantees the caller's fill loop always
	// makes progress, the way bufio.Reader did for the original parser.
	for tailLen < len(r.buf) {
		n, err := r.handle.Read(r.buf[tailLen:])
		tailLen += n
		if err == io.EOF {
			r.eofSeen = true
			break
		}
		if err != nil {
			r.pendingErr = err
			r.eofSeen = true
			break
		}
		if n > 0 {
			break
		}
	}
	r.filled = tailLen
}

// ensureFirstNewline ensures buf holds a '\n' at or after parseOff (filling as
// needed), returning its index, or filled if EOF is reached without one. fill()
// preserves parseOff as the current line start, so the returned index is valid
// in the current buf.
func (r *ZeroCopyCSVReader) ensureFirstNewline() int {
	for {
		if nl := bytes.IndexByte(r.buf[r.parseOff:r.filled], '\n'); nl >= 0 {
			return r.parseOff + nl
		}
		if r.eofSeen {
			return r.filled
		}
		r.fill()
	}
}

// scanRecord returns the record region [parseOff, lineEnd) (trailing CR
// stripped), the offset of the next record, whether the record contained a
// quote, and whether a complete record terminator was found within buf.
func (r *ZeroCopyCSVReader) scanRecord() (lineEnd, nextOff int, hasQuote, complete bool) {
	// Fast path: locate the next newline, then check the segment for a quote.
	if nl := bytes.IndexByte(r.buf[r.parseOff:r.filled], '\n'); nl >= 0 {
		segEnd := r.parseOff + nl
		if bytes.IndexByte(r.buf[r.parseOff:segEnd], '"') < 0 {
			end := segEnd
			if end > r.parseOff && r.buf[end-1] == '\r' {
				end--
			}
			return end, segEnd + 1, false, true
		}
	}
	// Quote-aware scan: a newline only ends the record when not inside quotes.
	inQuote := false
	for i := r.parseOff; i < r.filled; i++ {
		c := r.buf[i]
		if c == '"' {
			hasQuote = true
			inQuote = !inQuote
		} else if c == '\n' && !inQuote {
			end := i
			if end > r.parseOff && r.buf[end-1] == '\r' {
				end--
			}
			return end, i + 1, hasQuote, true
		}
	}
	return r.filled, r.filled, hasQuote, false
}

func (r *ZeroCopyCSVReader) Read() ([]string, error) {
	// Ensure there is some unparsed data.
	for r.parseOff >= r.filled && !r.eofSeen {
		r.fill()
	}
	if r.parseOff >= r.filled {
		if r.pendingErr != nil {
			err := r.pendingErr
			r.pendingErr = nil
			return nil, err
		}
		return nil, io.EOF
	}

	// Comment lines are terminated by a single newline regardless of any quote
	// characters they contain, and MUST be detected before quote parsing --
	// otherwise a '"' inside a comment would make the quote-aware scan consume
	// across the newline and swallow the following record. The go-csv fork's
	// readLine includes the trailing terminator in the returned line, and
	// --pass-comments echoes it verbatim, so we hand back the raw line with its
	// \n / \r\n included.
	if r.comment != 0 && r.buf[r.parseOff] == r.comment {
		nl := r.ensureFirstNewline()
		start := r.parseOff
		if nl < r.filled {
			r.parseOff = nl + 1
		} else {
			r.parseOff = r.filled
		}
		r.linesConsumed += bytes.Count(r.buf[start:r.parseOff], nlByte)
		return []string{bytesToString(r.buf[start:r.parseOff])}, nil
	}

	// Ensure the whole record is buffered (it may span blocks, esp. quoted
	// records with embedded newlines).
	lineEnd, nextOff, hasQuote, complete := r.scanRecord()
	for !complete && !r.eofSeen {
		r.fill()
		lineEnd, nextOff, hasQuote, complete = r.scanRecord()
	}
	if !complete {
		// EOF with a final record lacking a trailing newline.
		lineEnd = r.filled
		if lineEnd > r.parseOff && r.buf[lineEnd-1] == '\r' {
			lineEnd--
		}
		nextOff = r.filled
	}

	recordStartLine := r.linesConsumed + 1
	record := r.buf[r.parseOff:lineEnd]
	consumedStart := r.parseOff
	r.parseOff = nextOff
	r.linesConsumed += bytes.Count(r.buf[consumedStart:nextOff], nlByte)

	if !hasQuote {
		return r.splitNoQuote(record), nil
	}
	return r.parseQuoted(record, recordStartLine)
}

// splitNoQuote splits an unquoted record into zero-copy field views.
func (r *ZeroCopyCSVReader) splitNoQuote(record []byte) []string {
	n := 1
	for i := 0; i < len(record); i++ {
		if record[i] == r.comma {
			n++
		}
	}
	out := make([]string, n)
	fi := 0
	start := 0
	for i := 0; i < len(record); i++ {
		if record[i] == r.comma {
			out[fi] = bytesToString(record[start:i])
			fi++
			start = i + 1
		}
	}
	out[fi] = bytesToString(record[start:])
	return out
}

// nlByte is the newline separator for bytes.Count (line accounting).
var nlByte = []byte{'\n'}

// parseQuoted delegates a quote-containing record to the proven go-csv parser
// over the isolated record bytes. Correct (including escaped quotes and
// embedded newlines) at the cost of allocation; rare for typical data. Because
// the sub-parser numbers lines relative to the record, any ParseError's line
// numbers are shifted to absolute file lines via startLine.
func (r *ZeroCopyCSVReader) parseQuoted(record []byte, startLine int) ([]string, error) {
	sub := csv.NewReader(bytes.NewReader(record))
	sub.Comma = rune(r.comma)
	sub.LazyQuotes = r.lazyQuotes
	sub.TrimLeadingSpace = r.trimLeadingSpace
	// Comment handling is applied at the line level above, not here.
	fields, err := sub.Read()
	if err != nil {
		if pe, ok := err.(*csv.ParseError); ok {
			offset := startLine - 1
			pe.StartLine += offset
			pe.Line += offset
		}
	}
	return fields, err
}
