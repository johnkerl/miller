// ================================================================
// This file is adapted from
// https://cs.opensource.google/go/go/+/refs/tags/go1.17.6:src/encoding/csv/writer.go
// and used in accordance with its open-source license.
//
// The reason for the fork is https://github.com/johnkerl/miller/issues/853.
// Namely, for colorized output which uses ANSI escape sequences, which have
// things like ';' in them, csv.Writer.Write wraps fields in double quotes if
// the CSV delimiter is ';', and likewise for other special characters in the
// ANSI escape sequences.
//
// ================================================================
// Copyright (c) 2009 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
// ================================================================

package output

import (
	"bufio"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/johnkerl/miller/pkg/colorizer"
)

var errInvalidDelim = errors.New("csv: invalid field or comment delimiter")

func (writer *RecordWriterCSV) WriteCSVRecordMaybeColorized(
	record []string,
	bufferedOutputStream *bufio.Writer,
	outputIsStdout bool,
	isKey bool,
	quoteAll bool,
) error {
	comma := writer.csvWriter.Comma

	// Output prefix + field + suffix to the screen, but do needs-quoting
	// checks on the original fields.  Also note colorization is conditional on
	// command-line flags, whether the output is to stdout, etc. -- the prefix
	// and suffix may be ANSI escape sequences, or may be the empty string.
	prefix, suffix := colorizer.GetColorization(outputIsStdout, isKey)

	if !validDelim(comma) {
		return errInvalidDelim
	}

	for i, field := range record {
		if i > 0 {
			if _, err := bufferedOutputStream.WriteRune(comma); err != nil {
				return err
			}
		}

		// TODO: refacctor the maybe-colorized to return prefix, middle, suffix. Much easier that way.

		// If we don't have to have a quoted field then just
		// write out the field and continue to the next field.
		needsQuotes := quoteAll || fieldNeedsQuotes(field, comma)
		if !needsQuotes {
			if _, err := bufferedOutputStream.WriteString(prefix); err != nil {
				return err
			}
			if _, err := bufferedOutputStream.WriteString(field); err != nil {
				return err
			}
			if _, err := bufferedOutputStream.WriteString(suffix); err != nil {
				return err
			}
			continue
		}

		if _, err := bufferedOutputStream.WriteString(prefix); err != nil {
			return err
		}
		if err := bufferedOutputStream.WriteByte('"'); err != nil {
			return err
		}
		for len(field) > 0 {
			// Search for special characters.
			j := strings.IndexAny(field, "\"\r\n")
			if j < 0 {
				j = len(field)
			}

			// Copy verbatim everything before the special character.
			if _, err := bufferedOutputStream.WriteString(field[:j]); err != nil {
				return err
			}
			field = field[j:]

			// Encode the special character.
			if len(field) > 0 {
				var err error
				switch field[0] {
				case '"':
					if _, err := bufferedOutputStream.WriteString(`""`); err != nil {
						return err
					}
				case '\r':
					if !writer.csvWriter.UseCRLF {
						if err := bufferedOutputStream.WriteByte('\r'); err != nil {
							return err
						}
					}
				case '\n':
					if writer.csvWriter.UseCRLF {
						if _, err := bufferedOutputStream.WriteString("\r\n"); err != nil {
							return err
						}
					} else {
						if err := bufferedOutputStream.WriteByte('\n'); err != nil {
							return err
						}
					}
				}
				field = field[1:]
				if err != nil {
					return err
				}
			}
		}
		if err := bufferedOutputStream.WriteByte('"'); err != nil {
			return err
		}
		if _, err := bufferedOutputStream.WriteString(suffix); err != nil {
			return err
		}
	}
	if writer.csvWriter.UseCRLF {
		if _, err := bufferedOutputStream.WriteString("\r\n"); err != nil {
			return err
		}
	} else {
		if err := bufferedOutputStream.WriteByte('\n'); err != nil {
			return err
		}
	}
	return nil
}

func validDelim(r rune) bool {
	return r != 0 && r != '"' && r != '\r' && r != '\n' && utf8.ValidRune(r) && r != utf8.RuneError
}

// fieldNeedsQuotes reports whether our field must be enclosed in quotes.
// Fields with a Comma, fields with a quote or newline, and
// fields which start with a space must be enclosed in quotes.
// [NOTE: https://www.rfc-editor.org/rfc/rfc4180 doesn't specify this so Miller
// does not use this.]
// We used to quote empty strings, but we do not anymore (as of Go 1.4).
// The two representations should be equivalent, but Postgres distinguishes
// quoted vs non-quoted empty string during database imports, and it has
// an option to force the quoted behavior for non-quoted CSV but it has
// no option to force the non-quoted behavior for quoted CSV, making
// CSV with quoted empty strings strictly less useful.
// Not quoting the empty string also makes this package match the behavior
// of Microsoft Excel and Google Drive.
// For Postgres, quote the data terminating string `\.`.
func fieldNeedsQuotes(field string, comma rune) bool {

	if field == "" {
		return false
	}

	if field == `\.` {
		return true
	}

	if comma < utf8.RuneSelf {
		for i := 0; i < len(field); i++ {
			c := field[i]
			if c == '\n' || c == '\r' || c == '"' || c == byte(comma) {
				return true
			}
		}
	} else {
		if strings.ContainsRune(field, comma) || strings.ContainsAny(field, "\"\r\n") {
			return true
		}
	}

	// Not used by Miller as noted above
	// r1, _ := utf8.DecodeRuneInString(field)
	// return unicode.IsSpace(r1)
	return false
}
