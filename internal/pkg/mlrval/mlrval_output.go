package mlrval

import (
	"fmt"
	"os"
	"strconv"
)

// Must have non-pointer receiver in order to implement the fmt.Stringer
// interface to make this printable via fmt.Println et al.  However, that
// results in a needless copy of the Mlrval. So, we intentionally use pointer
// receiver, and if we need to print to stdout, we can fmt.Printf with "%s" and
// mv.String().
func (mv *Mlrval) String() string {
	// TODO: comment re deferral -- important perf effect!
	// if mv.IsFloat() && floatOutputFormatter != nil
	// if mv.mvtype == MT_FLOAT && floatOutputFormatter != nil {
	//if floatOutputFormatter != nil && (mv.mvtype == MT_FLOAT || mv.mvtype == MT_PENDING) {
	if floatOutputFormatter != nil && mv.Type() == MT_FLOAT {
		// Use the format string from global --ofmt, if supplied
		return floatOutputFormatter.FormatFloat(mv.floatval)
	}

	// TODO: track dirty-flag checking / somesuch.
	// At present it's cumbersome to check if an array or map has been modified
	// and it's safest to always recompute the string-rep.
	if mv.IsArrayOrMap() {
		mv.printrepValid = false
	}

	mv.setPrintRep()
	return mv.printrep
}

// OriginalString gets the field value as a string regardless of --ofmt specification.
// E.g if the ofmt is "%.4f" and input is 3.1415926535, OriginalString() will return
// "3.1415926535" while String() will return "3.1416".
func (mv *Mlrval) OriginalString() string {
	if mv.printrepValid {
		return mv.printrep
	} else {
		return mv.String()
	}
}

// See mlrval.go for more about JIT-formatting of string backings
func (mv *Mlrval) setPrintRep() {
	if !mv.printrepValid {
		switch mv.mvtype {

		case MT_PENDING:
			// Should not have gotten outside of the JSON decoder, so flag this
			// clearly visually if it should (buggily) slip through to
			// user-level visibility.
			mv.printrep = "(bug-if-you-see-this:case=3)" // xxx constdef at top of file

		case MT_ERROR:
			mv.printrep = "(error)" // xxx constdef at top of file

		case MT_ABSENT:
			// Callsites should be using absence to do non-assigns, so flag
			// this clearly visually if it should (buggily) slip through to
			// user-level visibility.
			mv.printrep = "(bug-if-you-see-this:case=4)" // xxx constdef at top of file

		case MT_VOID:
			mv.printrep = "" // xxx constdef at top of file

		case MT_STRING:
			break

		case MT_INT:
			mv.printrep = strconv.FormatInt(mv.intval, 10)

		case MT_FLOAT:
			mv.printrep = strconv.FormatFloat(mv.floatval, 'f', -1, 64)

		case MT_BOOL:
			if mv.boolval == true {
				mv.printrep = "true"
			} else {
				mv.printrep = "false"
			}

		case MT_ARRAY:
			bytes, err := mv.MarshalJSON(JSON_MULTILINE, false)
			// maybe just InternalCodingErrorIf(err != nil)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			mv.printrep = string(bytes)

		case MT_MAP:
			bytes, err := mv.MarshalJSON(JSON_MULTILINE, false)
			// maybe just InternalCodingErrorIf(err != nil)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			mv.printrep = string(bytes)
		}
		mv.printrepValid = true
	}
}

// StringifyValuesRecursively is nominally for the `--jvquoteall` flag.
func (mv *Mlrval) StringifyValuesRecursively() {
	switch mv.mvtype {

	case MT_ARRAY:
		for i, _ := range mv.arrayval {
			mv.arrayval[i].StringifyValuesRecursively()
		}

	case MT_MAP:
		for pe := mv.mapval.Head; pe != nil; pe = pe.Next {
			pe.Value.StringifyValuesRecursively()
		}

	default:
		mv.SetFromString(mv.String())
	}
}
