/*
Copyright (c) 2013 Jeremy Jay

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

// Package strptime provides a C-style strptime wrappers for time.Parse.
//
// It supports the following subset of format strings (stolen from python docs):
//     %d  Day of the month as a zero-padded decimal number.
//     %b  Month as locale’s abbreviated name.
//     %B  Month as locale’s full name.
//     %m  Month as a zero-padded decimal number.
//     %y  Year without century as a zero-padded decimal number.
//     %Y  Year with century as a decimal number.
//     %H  Hour (24-hour clock) as a zero-padded decimal number.
//     %I  Hour (12-hour clock) as a zero-padded decimal number.
//     %p  Locale’s equivalent of either AM or PM.
//     %M  Minute as a zero-padded decimal number.
//     %S  Second as a zero-padded decimal number.
//     %f  Microsecond as a decimal number, zero-padded on the left.
//     %z  UTC offset in the form +HHMM or -HHMM.
//     %Z  Time zone name. UTC, EST, CST
//     %%  A literal '%' character.
//
// BUG(pbnjay): If an unsupported specifier is used, it may NOT directly precede a
// supported specifier (i.e. there must be intervening text to match first)

// Local mods (johnkerl 2021-10-17): ParseTZ and strptime_tz supporting
// Miller's idiosyncracies.

package strptime

import (
	"errors"
	"os"
	"strings"
	"time"
)

// Parse accepts a percent-encoded strptime format string, converts it for use with
// time.Parse, and returns the resulting time.Time value. If non-date-related format
// text does not match within the string value, then ErrFormatMismatch will be returned.
// Errors from time.Parse are passed through untouched.
//
// If a unsupported format specifier is provided, it will be ignored and matching
// text will be skipped. To receive errors for unsupported formats, use ParseStrict or call Check.
func Parse(value, format string) (time.Time, error) {
	return strptime_tz(value, format, true, false, nil)
}

// ParseLocal is like Parse except it consults the $TZ environment variable.
// This is for Miller.
func ParseLocal(value, format string) (time.Time, error) {
	return strptime_tz(value, format, true, true, nil)
}

// ParseLocation is like Parse except it uses the specified location (timezone).
// This is for Miller.
func ParseLocation(value, format string, location *time.Location) (time.Time, error) {
	return strptime_tz(value, format, true, true, location)
}

// ParseStrict returns ErrFormatUnsupported for unsupported formats strings, but is otherwise
// identical to Parse.
func ParseStrict(value, format string) (time.Time, error) {
	return strptime_tz(value, format, false, false, nil)
}

// MustParse is a wrapper for Parse which panics on any error.
func MustParse(value, format string) time.Time {
	t, err := strptime_tz(value, format, true, false, nil)
	if err != nil {
		panic(err)
	}
	return t
}

// Check verifies that format is a fully-supported strptime format string for this implementation.
func Check(format string) error {
	parts := strings.Split(format, "%")
	for _, ps := range parts {
		// since we split on '%', this is the format code
		c := int(ps[0])
		if c == '%' {
			continue
		}
		if _, found := formatMap[c]; !found {
			return ErrFormatUnsupported
		}
	}

	return nil
}

func strptime_tz(value, format string, ignoreUnsupported bool, useTZ bool, location *time.Location) (time.Time, error) {
	parseStr := ""
	parseFmt := ""
	vi := 0

	parts := strings.Split(format, "%")
	for pi, ps := range parts {
		if pi == 0 {
			// check prefix string
			if value[:len(ps)] != ps {
				return time.Time{}, ErrFormatMismatch
			}
			vi += len(ps)
			continue
		}
		// since we split on '%', this is the format code
		c := int(ps[0])

		if c == '%' { // handle %% quickly
			if ps != value[vi:vi+len(ps)] {
				return time.Time{}, ErrFormatMismatch
			}
			vi += len(ps)
			continue
		}

		// Check if format is supported and get the time.Parse translation
		f, supported := formatMap[c]
		if !supported && !ignoreUnsupported {
			return time.Time{}, ErrFormatUnsupported
		}

		// Check the intervening text between format strings.
		// There may be some edge cases where this isn't quite right
		// but if that's the case you've got other problems...
		vj := len(ps) - 1
		if vj > 0 {
			vj = strings.Index(value[vi:], ps[1:])
		}
		if vj == -1 {
			return time.Time{}, ErrFormatMismatch
		}

		if supported {
			// Build up a new format and date string
			if vj == 0 { // no intervening text
				if c == 'f' {
					vj = len(value) - vi
				} else {
					vj = len(f)
					if vj > len(value)-vi {
						return time.Time{}, ErrFormatMismatch
					}
				}
			}

			if c == 'f' {
				parseFmt += "." + f
				parseStr += "." + value[vi:vi+vj]
			} else if c == 'p' {
				parseFmt += " " + f
				parseStr += " " + strings.ToUpper(value[vi:vi+vj])
			} else {
				parseFmt += " " + f
				parseStr += " " + value[vi:vi+vj]
			}
		}

		if !supported && vj == 0 {
			// ignore to the end of the string
			vi = len(value)
		} else {
			vi += (len(ps) - 1) + vj
		}
	}

	if vi < len(value) {
		// extra text on end of value
		return time.Time{}, ErrFormatMismatch
	}

	if useTZ {
		if location != nil {
			return time.ParseInLocation(parseFmt, parseStr, location)
		} else {
			tz := os.Getenv("TZ")
			if tz == "" {
				return time.Parse(parseFmt, parseStr)
			} else {
				location, err := time.LoadLocation(tz)
				if err != nil {
					return time.Time{}, err
				}
				return time.ParseInLocation(parseFmt, parseStr, location)
			}
		}
	} else {
		return time.Parse(parseFmt, parseStr)
	}
}

var (
	// ErrFormatMismatch means that intervening text in the strptime format string did not
	// match within the parsed string.
	ErrFormatMismatch = errors.New("date format mismatch")
	// ErrFormatUnsupported means that the format string includes unsupport percent-escapes.
	ErrFormatUnsupported = errors.New("date format contains unsupported percent-encodings")

	formatMap = map[int]string{
		'd': "02",
		'b': "Jan",
		'B': "January",
		'm': "01",
		'y': "06",
		'Y': "2006",
		'H': "15",
		'I': "03",
		'p': "PM",
		'M': "04",
		'S': "05",
		'f': "999999",
		'z': "-0700",
		'Z': "MST",
	}
)
