package types

import (
	"fmt"
	"os"
	"strconv"
)

// Empty string means use default format.
// Set from the CLI parser using mlr --ofmt.
var mlrvalFloatOutputFormatter IMlrvalFormatter = nil

func SetMlrvalFloatOutputFormat(formatString string) error {
	formatter, err := GetMlrvalFormatter(formatString)
	if err != nil {
		return err
	}
	mlrvalFloatOutputFormatter = formatter
	return nil
}

// See mlrval.go for more about JIT-formatting of string backings
func (mv *Mlrval) setPrintRep() {
	if !mv.printrepValid {
		// xxx do it -- disposition vector
		// xxx temp temp temp temp temp
		switch mv.mvtype {
		case MT_PENDING:
			// Should not have gotten outside of the JSON decoder, so flag this
			// clearly visually if it should (buggily) slip through to
			// user-level visibility.
			mv.printrep = "(bug-if-you-see-this)" // xxx constdef at top of file
			break
		case MT_ERROR:
			mv.printrep = "(error)" // xxx constdef at top of file
			break
		case MT_ABSENT:
			// Callsites should be using absence to do non-assigns, so flag
			// this clearly visually if it should (buggily) slip through to
			// user-level visibility.
			mv.printrep = "(bug-if-you-see-this)" // xxx constdef at top of file
			break
		case MT_VOID:
			mv.printrep = "" // xxx constdef at top of file
			break
		case MT_STRING:
			// panic i suppose
			break
		case MT_INT:
			mv.printrep = strconv.Itoa(mv.intval)
			break
		case MT_FLOAT:
			mv.printrep = strconv.FormatFloat(mv.floatval, 'f', -1, 64)
			break
		case MT_BOOL:
			if mv.boolval == true {
				mv.printrep = "true"
			} else {
				mv.printrep = "false"
			}
			break
		// TODO: handling indentation
		case MT_ARRAY:

			bytes, err := mv.MarshalJSON(JSON_MULTILINE, false)
			// maybe just InternalCodingErrorIf(err != nil)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			mv.printrep = string(bytes)

			break
		case MT_MAP:

			bytes, err := mv.MarshalJSON(JSON_MULTILINE, false)
			// maybe just InternalCodingErrorIf(err != nil)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			mv.printrep = string(bytes)

			break
		}
		mv.printrepValid = true
	}
}

// Must have non-pointer receiver in order to implement the fmt.Stringer
// interface to make this printable via fmt.Println et al.
func (mv Mlrval) String() string {
	if mv.mvtype == MT_FLOAT && mlrvalFloatOutputFormatter != nil {
		// Use the format string from global --ofmt, if supplied
		return mlrvalFloatOutputFormatter.FormatFloat(mv.floatval)
	} else {
		mv.setPrintRep()
		return mv.printrep
	}
}
