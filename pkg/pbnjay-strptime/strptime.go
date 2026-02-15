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
// Miller's idiosyncrasies.

package strptime

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

const _ignoreUnsupported = false

var _debug = os.Getenv("MLR_DEBUG_STRPTIME") != ""

// Parse accepts a percent-encoded strptime format string, converts it for use with
// time.Parse, and returns the resulting time.Time value. If non-date-related format
// text does not match within the string value, then ErrFormatMismatch will be returned.
// Errors from time.Parse are passed through untouched.
//
// If a unsupported format specifier is provided, it will be ignored and matching
// text will be skipped. To receive errors for unsupported formats, use ParseStrict or call Check.
func Parse(value, format string) (time.Time, error) {
	return strptime_tz(value, format, _ignoreUnsupported, false, nil)
}

// ParseLocal is like Parse except it consults the $TZ environment variable.
// This is for Miller.
func ParseLocal(value, format string) (time.Time, error) {
	return strptime_tz(value, format, _ignoreUnsupported, true, nil)
}

// ParseLocation is like Parse except it uses the specified location (timezone).
// This is for Miller.
func ParseLocation(value, format string, location *time.Location) (time.Time, error) {
	return strptime_tz(value, format, _ignoreUnsupported, true, location)
}

// Check verifies that format is a fully-supported strptime format string for this implementation.
// Not used by Miller.
func Check(format string) error {
	format = expandShorthands(format)

	parts := strings.Split(format, "%")
	for _, ps := range parts {
		// Since we split on '%', this is the format code

		// This is for "%%"
		if ps == "" {
			continue
		}

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

func strptime_tz(
	strptime_input, strptime_format string, ignoreUnsupported bool, useTZ bool, location *time.Location,
) (time.Time, error) {
	if _debug {
		fmt.Printf("================================================================ STRPTIME ENTER\n")
		fmt.Printf("strptime_input    \"%s\"\n", strptime_input)
		fmt.Printf("strptime_format   \"%s\"\n", strptime_format)
		defer fmt.Printf("================================================================ STRPTIME EXIT\n")
	}

	// E.g. re-write "%F" to "%Y-%m-%d".
	strptime_format = expandShorthands(strptime_format)
	if _debug {
		fmt.Printf("strptime_input    \"%s\"\n", strptime_input)
	}

	// The job of strptime is to map "format strings" like "%Y-%m-%d %H:%M:%S" to
	// Go-library "templates" like "2006 01 02 15 04 05".
	//
	// The way this works within pbnjay/strptime is to split the format string on "%", then walk
	// through and modify the input string as well.
	//
	// Example:
	// * strptime("2015-08-28T13:33:21Z", "%Y-%m-%dT%H:%M:%SZ")
	// * strptime input  "2015-08-28T13:33:21Z"
	// * strptime format "%Y-%m-%dT%H:%M:%SZ"
	// * go-lib input    "2015 08 28 13 33 21"
	// * go-lib template "2006 01 02 15 04 05"
	//
	// Note that since we split the strptime-style format string on "%", the first character in each
	// part is a format character like 'Y', 'm', etc -- except for the very start of the format
	// string which may have some prefix text before its very first percent sign.

	goLibInput := ""
	goLibTemplate := ""
	// inputIdx: index into strptime_input (the date string). Format is walked via partsIndex.
	inputIdx := 0
	firstComponent := true // no leading space before first component (avoids Go parse quirks)

	partsBetweenPercentSigns := strings.Split(strptime_format, "%")
	nparts := len(partsBetweenPercentSigns)
	for partsIndex := 0; partsIndex < nparts; /* increment in loop */ {
		partBetweenPercentSigns := partsBetweenPercentSigns[partsIndex]

		if _debug {
			fmt.Printf("\n")
			fmt.Printf("partsIndex %d:     \"%s\"\n", partsIndex, partBetweenPercentSigns)
		}
		if partsIndex == 0 {
			// Check for prefix text. It must be an exact match, e.g. with input "foo 2021" and
			// format "foo %Y", "foo " == "foo ". Or, if the format starts with a "%", we're
			// checking "" == "".
			if strptime_input[inputIdx:inputIdx+len(partBetweenPercentSigns)] != partBetweenPercentSigns {
				if _debug {
					fmt.Printf("\"%s\" != \"%s\"\n",
						strptime_input[inputIdx:inputIdx+len(partBetweenPercentSigns)], partBetweenPercentSigns,
					)
				}
				return time.Time{}, ErrFormatMismatch
			}
			inputIdx += len(partBetweenPercentSigns)
			partsIndex++
			continue
		}

		// Handle %% straight off, as this is a special case.
		if partBetweenPercentSigns == "" {
			if _debug {
				fmt.Printf("formatCode        '%c'\n", '%')
			}
			if strptime_input[inputIdx:inputIdx+1] != "%" {
				if _debug {
					fmt.Println("did not match %%")
				}
				return time.Time{}, ErrFormatMismatch
			}

			if _debug {
				fmt.Printf("templateComponent \"%s\"\n", "%")
				fmt.Printf("inputComponent    \"%s\"\n", "%")
			}

			inputIdx += 1
			partsIndex += 2 // TODO: TYPE ME UP
			continue
		}

		// Since we split on '%', this is the format code
		formatCode := int(partBetweenPercentSigns[0])

		// Check if the format code is supported, and map the strptime-style format code to the
		// Go-library (time.Parse) template component, e.g. 'Y' -> "2006".
		templateComponent, supported := formatMap[formatCode]
		if !supported && !ignoreUnsupported {
			if _debug {
				fmt.Printf("formatCode '%c' is unsupported\n", formatCode)
			}
			return time.Time{}, ErrFormatUnsupported
		}
		if _debug {
			fmt.Printf("formatCode        '%c'\n", formatCode)
			fmt.Printf("templateComponent \"%s\"\n", templateComponent)
		}

		// Check the intervening text between format strings, e.g. the ":" in "%Y:%m".  There may be
		// some edge cases where this isn't quite right but if that's the case you've got other
		// problems ...

		// Subtract 1 for the format code itself. E.g. with "%Y:%m", splitting on "%", one piece
		// is "Y:". sil is the length of the ":" part.
		sil := len(partBetweenPercentSigns) - 1
		// Now sil becomes the offset of this part within the strptime-style input.
		if sil > 0 {
			sil = strings.Index(strptime_input[inputIdx:], partBetweenPercentSigns[1:])
		}
		if sil == -1 {
			if _debug {
				fmt.Printf("format/template mismatch 1\n")
			}
			return time.Time{}, ErrFormatMismatch
		}
		if _debug {
			fmt.Printf("inputComponent    \"%s\"\n", strptime_input[inputIdx:inputIdx+sil])
		}

		if supported {
			// Accumulate the go-lib style template and input strings.
			if sil == 0 { // No intervening text, e.g. "%Y%m%d"
				if formatCode == 'f' {
					// %f is optional decimal point + 1-6 digit runes (microseconds).
					// Do not consume the rest of the string so that %f%z works:
					// e.g. ".160001+0100" -> %f takes ".160001", %z takes "+0100".
					sil = parseFracLen(strptime_input[inputIdx:])
					if sil == 0 {
						if _debug {
							fmt.Printf("format/template mismatch: no fractional digits for %%f\n")
						}
						return time.Time{}, ErrFormatMismatch
					}
				} else {
					want := len(templateComponent)
					remaining := len(strptime_input) - inputIdx
					if remaining == 0 {
						if _debug {
							fmt.Printf("format/template mismatch 2\n")
						}
						return time.Time{}, ErrFormatMismatch
					}
					// Allow single-digit at end of string (e.g. "1989-1-2" for %Y-%m-%d); we zero-pad when building.
					if want > remaining {
						sil = remaining
					} else {
						sil = want
					}
				}
			}

			// Use the format's literal as separator after each value (e.g. "/" for %m/%d/%Y) so Go parses unambiguously.
			sep := partBetweenPercentSigns[1:]
			if firstComponent {
				sep = ""
				firstComponent = false
			}
			if formatCode == 'f' {
				goLibTemplate += "." + templateComponent
				goLibInput += "." + strptime_input[inputIdx:inputIdx+sil]
			} else if formatCode == 'p' {
				goLibTemplate += templateComponent + sep
				goLibInput += strings.ToUpper(strptime_input[inputIdx:inputIdx+sil]) + sep
			} else {
				comp := strptime_input[inputIdx : inputIdx+sil]
				// Zero-pad numeric fields so single-digit input works (e.g. "1/07/2022" for %d/%m/%Y).
				comp = zeroPadLeft(comp, len(templateComponent))
				goLibTemplate += templateComponent + sep
				goLibInput += comp + sep
			}
		}

		if !supported && sil == 0 {
			// Ignore to the end of the string
			inputIdx = len(strptime_input)
		} else {
			inputIdx += (len(partBetweenPercentSigns) - 1) + sil
		}
		partsIndex++
	}

	if inputIdx < len(strptime_input) {
		if _debug {
			fmt.Printf("Extra text on end of strptime_input\n")
		}
		return time.Time{}, ErrFormatMismatch
	}

	if _debug {
		fmt.Printf("goLibInput        \"%s\"\n", goLibInput)
		fmt.Printf("goLibTemplate     \"%s\"\n", goLibTemplate)
	}

	// Now call the Go time library with template and input formatted the way it wants.
	if useTZ {
		if location != nil {
			return time.ParseInLocation(goLibTemplate, goLibInput, location)
		} else {
			tz := os.Getenv("TZ")
			if tz == "" {
				return time.Parse(goLibTemplate, goLibInput)
			} else {
				location, err := time.LoadLocation(tz)
				if err != nil {
					return time.Time{}, err
				}
				return time.ParseInLocation(goLibTemplate, goLibInput, location)
			}
		}
	} else {
		// Parse in UTC so strptime (without _local) is deterministic and matches docs.
		return time.ParseInLocation(goLibTemplate, goLibInput, time.UTC)
	}
}

// zeroPadLeft pads s with leading zeros to length n. If s is already >= n chars or
// contains non-digits, s is returned unchanged. Used so %d/%m etc. accept both
// "1" and "01" (Go's time.Parse with "02"/"01" requires zero-padded).
func zeroPadLeft(s string, n int) string {
	if n <= 0 || len(s) >= n {
		return s
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return s
		}
	}
	return strings.Repeat("0", n-len(s)) + s
}

// parseFracLen returns the byte length of a strptime %f field in s: optional '.'
// followed by 1-6 digit runes (microseconds). Returns 0 if no valid fraction.
func parseFracLen(s string) int {
	if s == "" {
		return 0
	}
	n := 0
	if s[0] == '.' {
		n = 1
	}
	digits := 0
	for n < len(s) && digits < 6 && s[n] >= '0' && s[n] <= '9' {
		n++
		digits++
	}
	if digits == 0 {
		return 0
	}
	return n
}

// expandShorthands handles some shorthands that the C library uses, which we can easily
// replicate -- e.g. "%F" is "%Y-%m-%d".
func expandShorthands(format string) string {
	// TODO: mem cache
	format = strings.ReplaceAll(format, "%T", "%H:%M:%S")
	format = strings.ReplaceAll(format, "%D", "%m/%d/%y")
	format = strings.ReplaceAll(format, "%F", "%Y-%m-%d")
	format = strings.ReplaceAll(format, "%R", "%H:%M")
	format = strings.ReplaceAll(format, "%r", "%I:%M:%S %p")
	format = strings.ReplaceAll(format, "%T", "%H:%M:%S")
	// We've no %e in this package
	// format = strings.ReplaceAll(format, "%v", "%e-%b-%Y")
	return format
}

var (
	// ErrFormatMismatch means that intervening text in the strptime format string did not
	// match within the parsed string.
	ErrFormatMismatch = errors.New("date format mismatch")
	// ErrFormatUnsupported means that the format string includes unsupported percent-escapes.
	ErrFormatUnsupported = errors.New("date format contains unsupported percent-encodings")

	formatMap = map[int]string{
		'b': "Jan",
		'B': "January",
		'd': "02",
		'f': "999999",
		'H': "15",
		'I': "03",
		'j': "__2",
		'm': "01",
		'M': "04",
		'p': "PM",
		'S': "05",
		'y': "06",
		'Y': "2006",
		'z': "-0700",
		'Z': "MST",
	}
)
