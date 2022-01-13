package bifs

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
)

// ================================================================
func BIF_strlen(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsStringOrVoid() {
		return mlrval.ERROR
	} else {
		return mlrval.FromInt(int(utf8.RuneCountInString(input1.AcquireStringValue())))
	}
}

// ================================================================
func BIF_string(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromString(input1.String())
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

func dot_s_xx(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return mlrval.FromString(input1.String() + input2.String())
}

var dot_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT       FLOAT     BOOL      VOID   STRING    ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {dot_s_xx, dot_s_xx, dot_s_xx, _s1__, dot_s_xx, _erro, _erro, _erro, _erro, _1___, _s1__},
	/*FLOAT  */ {dot_s_xx, dot_s_xx, dot_s_xx, _s1__, dot_s_xx, _erro, _erro, _erro, _erro, _1___, _s1__},
	/*BOOL   */ {dot_s_xx, dot_s_xx, dot_s_xx, _s1__, dot_s_xx, _erro, _erro, _erro, _erro, _1___, _s1__},
	/*VOID   */ {_s2__, _s2__, _s2__, _void, _2___, _absn, _absn, _erro, _erro, _void, _void},
	/*STRING */ {dot_s_xx, dot_s_xx, dot_s_xx, _1___, dot_s_xx, _erro, _erro, _erro, _erro, _1___, _1___},
	/*ARRAY  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*MAP    */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _absn, _absn, _erro, _erro, _erro, _erro},
	/*NULL   */ {_s2__, _s2__, _s2__, _void, _2___, _absn, _absn, _erro, _erro, _null, _null},
	/*ABSENT */ {_s2__, _s2__, _s2__, _void, _2___, _absn, _absn, _erro, _erro, _null, _absn},
}

func BIF_dot(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return dot_dispositions[input1.Type()][input2.Type()](input1, input2)
}

// ================================================================
// substr1(s,m,n) gives substring of s from 1-up position m to n inclusive.
// Negative indices -len .. -1 alias to 0 .. len-1.

func BIF_substr_1_up(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsAbsent() {
		return mlrval.ABSENT
	}
	if input1.IsError() {
		return mlrval.ERROR
	}
	sinput := input1.String()

	// Handle UTF-8 correctly: len(input1.AcquireStringValue()) will count bytes, not runes.
	runes := []rune(sinput)
	strlen := int(len(runes))

	// For array slices like s[1:2], s[:2], s[1:], when the lower index is
	// empty in the DSL expression it comes in here as a 1. But when the upper
	// index is empty in the DSL expression it comes in here as "".
	if !input2.IsInt() {
		return mlrval.ERROR
	}
	lowerMindex := input2.AcquireIntValue()

	upperMindex := strlen
	if input3.IsVoid() {
		// Keep strlen
	} else if !input3.IsInt() {
		return mlrval.ERROR
	} else {
		upperMindex = input3.AcquireIntValue()
	}

	// Convert from negative-aliased 1-up to positive-only 0-up
	m, mok := unaliasArrayLengthIndex(strlen, lowerMindex)
	n, nok := unaliasArrayLengthIndex(strlen, upperMindex)

	if !mok || !nok {
		return mlrval.VOID
	} else if m > n {
		return mlrval.VOID
	} else {
		// Note Golang slice indices are 0-up, and the 1st index is inclusive
		// while the 2nd is exclusive. For Miller, indices are 1-up and both
		// are inclusive.
		return mlrval.FromString(string(runes[m : n+1]))
	}
}

// ================================================================
// substr0(s,m,n) gives substring of s from 0-up position m to n inclusive.
// Negative indices -len .. -1 alias to 0 .. len-1.

func BIF_substr_0_up(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsAbsent() {
		return mlrval.ABSENT
	}
	if input1.IsError() {
		return mlrval.ERROR
	}
	sinput := input1.String()

	// Handle UTF-8 correctly: len(input1.AcquireStringValue()) will count bytes, not runes.
	runes := []rune(sinput)
	strlen := int(len(runes))

	// For array slices like s[1:2], s[:2], s[1:], when the lower index is
	// empty in the DSL expression it comes in here as a 1. But when the upper
	// index is empty in the DSL expression it comes in here as "".
	if !input2.IsInt() {
		return mlrval.ERROR
	}
	lowerMindex := input2.AcquireIntValue()
	if lowerMindex >= 0 {
		// Make 1-up
		lowerMindex += 1
	}

	upperMindex := strlen
	if input3.IsVoid() {
		// Keep strlen
	} else if !input3.IsInt() {
		return mlrval.ERROR
	} else {
		upperMindex = input3.AcquireIntValue()
		if upperMindex >= 0 {
			// Make 1-up
			upperMindex += 1
		}
	}

	// Convert from negative-aliased 1-up to positive-only 0-up
	m, mok := unaliasArrayLengthIndex(strlen, lowerMindex)
	n, nok := unaliasArrayLengthIndex(strlen, upperMindex)

	if !mok || !nok {
		return mlrval.VOID
	} else if m > n {
		return mlrval.VOID
	} else {
		// Note Golang slice indices are 0-up, and the 1st index is inclusive
		// while the 2nd is exclusive. For Miller, indices are 1-up and both
		// are inclusive.
		return mlrval.FromString(string(runes[m : n+1]))
	}
}

// ================================================================
func BIF_truncate(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsErrorOrAbsent() {
		return input1
	}
	if input2.IsErrorOrAbsent() {
		return input2
	}
	if !input1.IsStringOrVoid() {
		return mlrval.ERROR
	}
	if !input2.IsInt() {
		return mlrval.ERROR
	}
	if input2.AcquireIntValue() < 0 {
		return mlrval.ERROR
	}

	// Handle UTF-8 correctly: len(input1.AcquireStringValue()) will count bytes, not runes.
	runes := []rune(input1.AcquireStringValue())
	oldLength := int(len(runes))
	maxLength := input2.AcquireIntValue()
	if oldLength <= maxLength {
		return input1
	} else {
		return mlrval.FromString(string(runes[0:maxLength]))
	}
}

// ================================================================
func BIF_lstrip(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsString() {
		return mlrval.FromString(strings.TrimLeft(input1.AcquireStringValue(), " \t"))
	} else {
		return input1
	}
}

func BIF_rstrip(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsString() {
		return mlrval.FromString(strings.TrimRight(input1.AcquireStringValue(), " \t"))
	} else {
		return input1
	}
}

func BIF_strip(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsString() {
		return mlrval.FromString(strings.Trim(input1.AcquireStringValue(), " \t"))
	} else {
		return input1
	}
}

// ----------------------------------------------------------------
func BIF_collapse_whitespace(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return BIF_collapse_whitespace_regexp(input1, WhitespaceRegexp())
}

func BIF_collapse_whitespace_regexp(input1 *mlrval.Mlrval, whitespaceRegexp *regexp.Regexp) *mlrval.Mlrval {
	if input1.IsString() {
		return mlrval.FromString(whitespaceRegexp.ReplaceAllString(input1.AcquireStringValue(), " "))
	} else {
		return input1
	}
}

func WhitespaceRegexp() *regexp.Regexp {
	return regexp.MustCompile("\\s+")
}

// ================================================================
func BIF_toupper(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsString() {
		return mlrval.FromString(strings.ToUpper(input1.AcquireStringValue()))
	} else if input1.IsVoid() {
		return input1
	} else {
		return input1
	}
}

func BIF_tolower(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsString() {
		return mlrval.FromString(strings.ToLower(input1.AcquireStringValue()))
	} else if input1.IsVoid() {
		return input1
	} else {
		return input1
	}
}

func BIF_capitalize(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsString() {
		if input1.AcquireStringValue() == "" {
			return input1
		} else {
			runes := []rune(input1.AcquireStringValue())
			rfirst := runes[0]
			rrest := runes[1:]
			sfirst := strings.ToUpper(string(rfirst))
			srest := string(rrest)
			return mlrval.FromString(sfirst + srest)
		}
	} else {
		return input1
	}
}

// ----------------------------------------------------------------
func BIF_clean_whitespace(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	return BIF_strip(
		BIF_collapse_whitespace_regexp(
			input1, WhitespaceRegexp(),
		),
	)
}

// ================================================================
func BIF_format(mlrvals []*mlrval.Mlrval) *mlrval.Mlrval {
	if len(mlrvals) == 0 {
		return mlrval.VOID
	}
	formatString, ok := mlrvals[0].GetStringValue()
	if !ok { // not a string
		return mlrval.ERROR
	}

	pieces := lib.SplitString(formatString, "{}")

	var buffer bytes.Buffer

	// Example: format("{}:{}", 8, 9)
	//
	// * piece[0] ""
	// * piece[1] ":"
	// * piece[2] ""
	// * mlrval[1] 8
	// * mlrval[2] 9
	//
	// So:
	// * Write piece[0]
	// * Write mlrvals[1]
	// * Write piece[1]
	// * Write mlrvals[2]
	// * Write piece[2]

	// Q: What if too few arguments for format?
	// A: Leave them off
	// Q: What if too many arguments for format?
	// A: Leave them off

	n := len(mlrvals)
	for i, piece := range pieces {
		if i > 0 {
			if i < n {
				buffer.WriteString(mlrvals[i].String())
			}
		}
		buffer.WriteString(piece)
	}

	return mlrval.FromString(buffer.String())
}

// ================================================================
func BIF_hexfmt(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsInt() {
		return mlrval.FromString("0x" + strconv.FormatUint(uint64(input1.AcquireIntValue()), 16))
	} else {
		return input1
	}
}

// ----------------------------------------------------------------
func fmtnum_is(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input2.IsString() {
		return mlrval.ERROR
	}
	formatString := input2.AcquireStringValue()
	formatter, err := mlrval.GetFormatter(formatString)
	if err != nil {
		return mlrval.ERROR
	}

	return formatter.Format(input1)
}

func fmtnum_fs(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input2.IsString() {
		return mlrval.ERROR
	}
	formatString := input2.AcquireStringValue()
	formatter, err := mlrval.GetFormatter(formatString)
	if err != nil {
		return mlrval.ERROR
	}

	return formatter.Format(input1)
}

func fmtnum_bs(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input2.IsString() {
		return mlrval.ERROR
	}
	formatString := input2.AcquireStringValue()
	formatter, err := mlrval.GetFormatter(formatString)
	if err != nil {
		return mlrval.ERROR
	}

	intMv := mlrval.FromInt(lib.BoolToInt(input1.AcquireBoolValue()))

	return formatter.Format(intMv)
}

var fmtnum_dispositions = [mlrval.MT_DIM][mlrval.MT_DIM]BinaryFunc{
	//       .  INT    FLOAT  BOOL   VOID   STRING     ARRAY  MAP    FUNC    ERROR   NULL   ABSENT
	/*INT    */ {_erro, _erro, _erro, _erro, fmtnum_is, _erro, _erro, _erro, _erro, _erro, _absn},
	/*FLOAT  */ {_erro, _erro, _erro, _erro, fmtnum_fs, _erro, _erro, _erro, _erro, _erro, _absn},
	/*BOOL   */ {_erro, _erro, _erro, _erro, fmtnum_bs, _erro, _erro, _erro, _erro, _erro, _absn},
	/*VOID   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn},
	/*STRING */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn},
	/*ARRAY  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn},
	/*MAP    */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn},
	/*FUNC   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*NULL   */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn},
	/*ABSENT */ {_absn, _absn, _erro, _absn, _absn, _erro, _erro, _erro, _erro, _absn, _absn},
}

func BIF_fmtnum(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return fmtnum_dispositions[input1.Type()][input2.Type()](input1, input2)
}
