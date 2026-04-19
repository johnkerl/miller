// Package dkvpx (see also dkvpx_reader.go) writes DKVPX records: comma-delimited
// key=value pairs with selective quoting. Keys and values are double-quoted only
// when they contain comma, equals, newline, or double-quote.
package dkvpx

import (
	"bufio"
	"io"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/lib"
)

// A Writer writes DKVPX records. It mirrors the structure of csv.Writer.
type Writer struct {
	Comma   rune // Pair delimiter (set to ',' by NewWriter)
	UseCRLF bool // True to use \r\n as the line terminator
	w       *bufio.Writer
}

// NewWriter returns a new Writer that writes to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		Comma: ',',
		w:     bufio.NewWriter(w),
	}
}

// Write writes a single DKVPX record. Keys and values are quoted only when
// they contain comma, equals, newline, or double-quote.
func (w *Writer) Write(record *lib.OrderedMap[string]) error {
	if record == nil || record.IsEmpty() {
		return w.writeLineEnd()
	}

	first := true
	for pe := record.Head; pe != nil; pe = pe.Next {
		if !first {
			if _, err := w.w.WriteRune(w.Comma); err != nil {
				return err
			}
		}
		first = false

		if err := w.writeQuotedIfNeeded(pe.Key); err != nil {
			return err
		}
		if err := w.w.WriteByte('='); err != nil {
			return err
		}
		if err := w.writeQuotedIfNeeded(pe.Value); err != nil {
			return err
		}
	}

	return w.writeLineEnd()
}

// writeQuotedIfNeeded writes the string, adding surrounding quotes and escaping
// internal quotes only when the string contains comma, equals, newline, or quote.
func (w *Writer) writeQuotedIfNeeded(s string) error {
	if !needsQuoting(s) {
		_, err := w.w.WriteString(s)
		return err
	}

	if err := w.w.WriteByte('"'); err != nil {
		return err
	}
	// Escape " as ""
	for len(s) > 0 {
		i := strings.IndexAny(s, "\"\r\n")
		if i < 0 {
			i = len(s)
		}
		if _, err := w.w.WriteString(s[:i]); err != nil {
			return err
		}
		s = s[i:]
		if len(s) > 0 {
			switch s[0] {
			case '"':
				if _, err := w.w.WriteString(`""`); err != nil {
					return err
				}
			case '\r', '\n':
				if err := w.w.WriteByte(s[0]); err != nil {
					return err
				}
			}
			s = s[1:]
		}
	}
	return w.w.WriteByte('"')
}

// needsQuoting reports whether the string must be quoted (contains comma,
// equals, newline, or double-quote).
func needsQuoting(s string) bool {
	return strings.ContainsAny(s, ",\n\r=\"")
}

// FormatField returns the string formatted for DKVPX output: quoted and escaped
// if it contains comma, equals, newline, or quote; otherwise unchanged.
// Useful for callers that need to apply colorization or other wrapping.
func FormatField(s string) string {
	if !needsQuoting(s) {
		return s
	}
	var b strings.Builder
	b.WriteByte('"')
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"':
			b.WriteString(`""`)
		default:
			b.WriteByte(s[i])
		}
	}
	b.WriteByte('"')
	return b.String()
}

func (w *Writer) writeLineEnd() error {
	if w.UseCRLF {
		_, err := w.w.WriteString("\r\n")
		return err
	}
	return w.w.WriteByte('\n')
}

// Flush writes any buffered data to the underlying io.Writer.
func (w *Writer) Flush() {
	w.w.Flush()
}

// Error reports any error that has occurred during a previous Write or Flush.
func (w *Writer) Error() error {
	_, err := w.w.Write(nil)
	return err
}

// WriteAll writes multiple DKVPX records and then calls Flush.
func (w *Writer) WriteAll(records []*lib.OrderedMap[string]) error {
	for _, record := range records {
		if err := w.Write(record); err != nil {
			return err
		}
	}
	return w.w.Flush()
}
