package types

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"mlr/internal/pkg/lib"
)

// ================================================================
func BIF_strlen(input1 *Mlrval) *Mlrval {
	if !input1.IsStringOrVoid() {
		return MLRVAL_ERROR
	} else {
		return MlrvalFromInt(int(utf8.RuneCountInString(input1.printrep)))
	}
}

// ================================================================
func BIF_string(input1 *Mlrval) *Mlrval {
	return MlrvalFromString(input1.String())
}

// ================================================================
// Dot operator, with loose typecasting.
//
// For most operations, I don't like loose typecasting -- for example, in PHP
// "10" + 2 is the number 12 and in JavaScript it's the string "102", and I
// find both of those horrid and error-prone. In Miller, "10"+2 is MT_ERROR, by
// design, unless intentional casting is done like '$x=int("10")+2'.
//
// However, for dotting, in practice I tipped over and allowed dotting of
// strings and ints: so while "10" + 2 is an error in Miller, '"10". 2' is
// "102". Unlike with "+", with "." there is no ambiguity about what the output
// should be: always the string concatenation of the string representations of
// the two arguments. So, we do the string-cast for the user.

func dot_s_xx(input1, input2 *Mlrval) *Mlrval {
	return MlrvalFromString(input1.String() + input2.String())
}

var dot_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT NULL   VOID   STRING    INT       FLOAT     BOOL      ARRAY  MAP     FUNC
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro},
	/*ABSENT */ {_erro, _absn, _null, _void, _2___, _s2__, _s2__, _s2__, _absn, _absn, _erro},
	/*NULL   */ {_erro, _null, _null, _void, _2___, _s2__, _s2__, _s2__, _absn, _absn, _erro},
	/*VOID   */ {_erro, _void, _void, _void, _2___, _s2__, _s2__, _s2__, _absn, _absn, _erro},
	/*STRING */ {_erro, _1___, _1___, _1___, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, _erro, _erro, _erro},
	/*INT    */ {_erro, _s1__, _1___, _s1__, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, _erro, _erro, _erro},
	/*FLOAT  */ {_erro, _s1__, _1___, _s1__, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, _erro, _erro, _erro},
	/*BOOL   */ {_erro, _s1__, _1___, _s1__, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, _erro, _erro, _erro},
	/*ARRAY  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*MAP    */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*FUNC    */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func BIF_dot(input1, input2 *Mlrval) *Mlrval {
	return dot_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}

// ================================================================
// substr1(s,m,n) gives substring of s from 1-up position m to n inclusive.
// Negative indices -len .. -1 alias to 0 .. len-1.

func BIF_substr_1_up(input1, input2, input3 *Mlrval) *Mlrval {
	if !input1.IsStringOrVoid() {
		return MLRVAL_ERROR
	}

	// Handle UTF-8 correctly: len(input1.printrep) will count bytes, not runes.
	runes := []rune(input1.printrep)
	strlen := int(len(runes))

	// For array slices like s[1:2], s[:2], s[1:], when the lower index is
	// empty in the DSL expression it comes in here as a 1. But when the upper
	// index is empty in the DSL expression it comes in here as "".
	if !input2.IsInt() {
		return MLRVAL_ERROR
	}
	lowerMindex := input2.intval

	upperMindex := strlen
	if input3.IsEmpty() {
		// Keep strlen
	} else if !input3.IsInt() {
		return MLRVAL_ERROR
	} else {
		upperMindex = input3.intval
	}

	// Convert from negative-aliased 1-up to positive-only 0-up
	m, mok := UnaliasArrayLengthIndex(strlen, lowerMindex)
	n, nok := UnaliasArrayLengthIndex(strlen, upperMindex)

	if !mok || !nok {
		return MLRVAL_VOID
	} else if m > n {
		return MLRVAL_VOID
	} else {
		// Note Golang slice indices are 0-up, and the 1st index is inclusive
		// while the 2nd is exclusive. For Miller, indices are 1-up and both
		// are inclusive.
		return MlrvalFromString(string(runes[m : n+1]))
	}
}

// ================================================================
// substr0(s,m,n) gives substring of s from 0-up position m to n inclusive.
// Negative indices -len .. -1 alias to 0 .. len-1.

func BIF_substr_0_up(input1, input2, input3 *Mlrval) *Mlrval {
	if !input1.IsStringOrVoid() {
		return MLRVAL_ERROR
	}

	// Handle UTF-8 correctly: len(input1.printrep) will count bytes, not runes.
	runes := []rune(input1.printrep)
	strlen := int(len(runes))

	// For array slices like s[1:2], s[:2], s[1:], when the lower index is
	// empty in the DSL expression it comes in here as a 1. But when the upper
	// index is empty in the DSL expression it comes in here as "".
	if !input2.IsInt() {
		return MLRVAL_ERROR
	}
	lowerMindex := input2.intval
	if lowerMindex >= 0 {
		// Make 1-up
		lowerMindex += 1
	}

	upperMindex := strlen
	if input3.IsEmpty() {
		// Keep strlen
	} else if !input3.IsInt() {
		return MLRVAL_ERROR
	} else {
		upperMindex = input3.intval
		if upperMindex >= 0 {
			// Make 1-up
			upperMindex += 1
		}
	}

	// Convert from negative-aliased 1-up to positive-only 0-up
	m, mok := UnaliasArrayLengthIndex(strlen, lowerMindex)
	n, nok := UnaliasArrayLengthIndex(strlen, upperMindex)

	if !mok || !nok {
		return MLRVAL_VOID
	} else if m > n {
		return MLRVAL_VOID
	} else {
		// Note Golang slice indices are 0-up, and the 1st index is inclusive
		// while the 2nd is exclusive. For Miller, indices are 1-up and both
		// are inclusive.
		return MlrvalFromString(string(runes[m : n+1]))
	}
}

// ================================================================
func BIF_truncate(input1, input2 *Mlrval) *Mlrval {
	if input1.IsErrorOrAbsent() {
		return input1
	}
	if input2.IsErrorOrAbsent() {
		return input2
	}
	if !input1.IsStringOrVoid() {
		return MLRVAL_ERROR
	}
	if !input2.IsInt() {
		return MLRVAL_ERROR
	}
	if input2.intval < 0 {
		return MLRVAL_ERROR
	}

	// Handle UTF-8 correctly: len(input1.printrep) will count bytes, not runes.
	runes := []rune(input1.printrep)
	oldLength := int(len(runes))
	maxLength := input2.intval
	if oldLength <= maxLength {
		return input1
	} else {
		return MlrvalFromString(string(runes[0:maxLength]))
	}
}

// ================================================================
func BIF_lstrip(input1 *Mlrval) *Mlrval {
	if input1.mvtype == MT_STRING {
		return MlrvalFromString(strings.TrimLeft(input1.printrep, " \t"))
	} else {
		return input1
	}
}

func BIF_rstrip(input1 *Mlrval) *Mlrval {
	if input1.mvtype == MT_STRING {
		return MlrvalFromString(strings.TrimRight(input1.printrep, " \t"))
	} else {
		return input1
	}
}

func BIF_strip(input1 *Mlrval) *Mlrval {
	if input1.mvtype == MT_STRING {
		return MlrvalFromString(strings.Trim(input1.printrep, " \t"))
	} else {
		return input1
	}
}

// ----------------------------------------------------------------
func BIF_collapse_whitespace(input1 *Mlrval) *Mlrval {
	return MlrvalCollapseWhitespaceRegexp(input1, WhitespaceRegexp())
}

func MlrvalCollapseWhitespaceRegexp(input1 *Mlrval, whitespaceRegexp *regexp.Regexp) *Mlrval {
	if input1.mvtype == MT_STRING {
		return MlrvalFromString(whitespaceRegexp.ReplaceAllString(input1.printrep, " "))
	} else {
		return input1
	}
}

func WhitespaceRegexp() *regexp.Regexp {
	return regexp.MustCompile("\\s+")
}

// ================================================================
func BIF_toupper(input1 *Mlrval) *Mlrval {
	if input1.mvtype == MT_STRING {
		return MlrvalFromString(strings.ToUpper(input1.printrep))
	} else if input1.mvtype == MT_VOID {
		return input1
	} else {
		return input1
	}
}

func BIF_tolower(input1 *Mlrval) *Mlrval {
	if input1.mvtype == MT_STRING {
		return MlrvalFromString(strings.ToLower(input1.printrep))
	} else if input1.mvtype == MT_VOID {
		return input1
	} else {
		return input1
	}
}

func BIF_capitalize(input1 *Mlrval) *Mlrval {
	if input1.mvtype == MT_STRING {
		if input1.printrep == "" {
			return input1
		} else {
			runes := []rune(input1.printrep)
			rfirst := runes[0]
			rrest := runes[1:]
			sfirst := strings.ToUpper(string(rfirst))
			srest := string(rrest)
			return MlrvalFromString(sfirst + srest)
		}
	} else {
		return input1
	}
}

// ----------------------------------------------------------------
func BIF_clean_whitespace(input1 *Mlrval) *Mlrval {
	return BIF_strip(
		MlrvalCollapseWhitespaceRegexp(
			input1, WhitespaceRegexp(),
		),
	)
}

// ================================================================
func BIF_hexfmt(input1 *Mlrval) *Mlrval {
	if input1.mvtype == MT_INT {
		return MlrvalFromString("0x" + strconv.FormatUint(uint64(input1.intval), 16))
	} else {
		return input1
	}
}

// ----------------------------------------------------------------
func fmtnum_is(input1, input2 *Mlrval) *Mlrval {
	if !input2.IsString() {
		return MLRVAL_ERROR
	}
	formatString := input2.printrep
	formatter, err := GetMlrvalFormatter(formatString)
	if err != nil {
		return MLRVAL_ERROR
	}

	return formatter.Format(input1)
}

func fmtnum_fs(input1, input2 *Mlrval) *Mlrval {
	if !input2.IsString() {
		return MLRVAL_ERROR
	}
	formatString := input2.printrep
	formatter, err := GetMlrvalFormatter(formatString)
	if err != nil {
		return MLRVAL_ERROR
	}

	return formatter.Format(input1)
}

func fmtnum_bs(input1, input2 *Mlrval) *Mlrval {
	if !input2.IsString() {
		return MLRVAL_ERROR
	}
	formatString := input2.printrep
	formatter, err := GetMlrvalFormatter(formatString)
	if err != nil {
		return MLRVAL_ERROR
	}

	intMv := MlrvalFromInt(lib.BoolToInt(input1.boolval))

	return formatter.Format(intMv)
}

var fmtnum_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT NULL   VOID   STRING     INT    FLOAT  BOOL   ARRAY  MAP    FUNC
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _absn, _absn, _absn, _absn, _absn, _erro, _erro, _erro, _erro},
	/*NULL   */ {_erro, _absn, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*VOID   */ {_erro, _absn, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*STRING */ {_erro, _absn, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*INT    */ {_erro, _absn, _erro, _erro, fmtnum_is, _erro, _erro, _erro, _erro, _erro, _erro},
	/*FLOAT  */ {_erro, _absn, _erro, _erro, fmtnum_fs, _erro, _erro, _erro, _erro, _erro, _erro},
	/*BOOL   */ {_erro, _absn, _erro, _erro, fmtnum_bs, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ARRAY  */ {_erro, _absn, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*MAP    */ {_erro, _absn, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func BIF_fmtnum(input1, input2 *Mlrval) *Mlrval {
	return fmtnum_dispositions[input1.mvtype][input2.mvtype](input1, input2)
}
