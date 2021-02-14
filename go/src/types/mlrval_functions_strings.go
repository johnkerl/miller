package types

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"miller/src/lib"
)

// ================================================================
func MlrvalStrlen(ma *Mlrval) Mlrval {
	if !ma.IsStringOrVoid() {
		return MlrvalFromError()
	}
	return MlrvalFromInt(int(utf8.RuneCountInString(ma.printrep)))
}

// ================================================================
func MlrvalToString(ma *Mlrval) Mlrval {
	return MlrvalFromString(ma.String())
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

func dot_s_xx(ma, mb *Mlrval) Mlrval {
	return MlrvalFromString(ma.String() + mb.String())
}

var dot_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT VOID   STRING    INT       FLOAT     BOOL      ARRAY     MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _absn, _absn},
	/*ABSENT */ {_erro, _absn, _void, _2___, _s2__, _s2__, _s2__, _absn, _absn},
	/*VOID   */ {_erro, _void, _void, _2___, _s2__, _s2__, _s2__, _absn, _absn},
	/*STRING */ {_erro, _1___, _1___, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx},
	/*INT    */ {_erro, _s1__, _s1__, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx},
	/*FLOAT  */ {_erro, _s1__, _s1__, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx},
	/*BOOL   */ {_erro, _s1__, _s1__, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx},
	/*ARRAY  */ {_absn, _absn, _absn, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx},
	/*MAP    */ {_absn, _absn, _absn, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx, dot_s_xx},
}

func MlrvalDot(ma, mb *Mlrval) Mlrval {
	return dot_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}

// ================================================================
// substr(s,m,n) gives substring of s from 1-up position m to n inclusive.
// Negative indices -len .. -1 alias to 0 .. len-1.

func MlrvalSubstr(ma, mb, mc *Mlrval) Mlrval {
	if !ma.IsStringOrVoid() {
		if ma.IsNumeric() {
			// JIT-stringify, if not already (e.g. intval scanned from string
			// in input-file data)
			ma.setPrintRep()
		} else {
			return MlrvalFromError()
		}
	}
	// TODO: fix this with regard to UTF-8 and runes.
	strlen := int(len(ma.printrep))

	// For array slices like s[1:2], s[:2], s[1:], when the lower index is
	// empty in the DSL expression it comes in here as a 1. But when the upper
	// index is empty in the DSL expression it comes in here as "".
	if !mb.IsInt() {
		return MlrvalFromError()
	}
	lowerMindex := mb.intval

	upperMindex := strlen
	if mc.IsEmpty() {
		// Keep strlen
	} else if !mc.IsInt() {
		return MlrvalFromError()
	} else {
		upperMindex = mc.intval
	}

	// Convert from negative-aliased 1-up to positive-only 0-up
	m, mok := UnaliasArrayLengthIndex(strlen, lowerMindex)
	n, nok := UnaliasArrayLengthIndex(strlen, upperMindex)

	if !mok || !nok {
		return MlrvalFromString("")
	} else if m > n {
		return MlrvalFromError()
	} else {
		// Note Golang slice indices are 0-up, and the 1st index is inclusive
		// while the 2nd is exclusive. For Miller, indices are 1-up and both
		// are inclusive.
		return MlrvalFromString(ma.printrep[m : n+1])
	}
}

// ================================================================
func MlrvalTruncate(ma, mb *Mlrval) Mlrval {
	if ma.IsErrorOrAbsent() {
		return *ma
	}
	if mb.IsErrorOrAbsent() {
		return *mb
	}
	if !ma.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !mb.IsInt() {
		return MlrvalFromError()
	}
	if mb.intval < 0 {
		return MlrvalFromError()
	}

	oldLength := int(len(ma.printrep))
	maxLength := mb.intval
	if oldLength <= maxLength {
		return *ma
	} else {
		return MlrvalFromString(ma.printrep[0:maxLength])
	}
}

// ================================================================
func MlrvalLStrip(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_STRING {
		return MlrvalFromString(strings.TrimLeft(ma.printrep, " \t"))
	} else {
		return *ma
	}
}

func MlrvalRStrip(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_STRING {
		return MlrvalFromString(strings.TrimRight(ma.printrep, " \t"))
	} else {
		return *ma
	}
}

func MlrvalStrip(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_STRING {
		return MlrvalFromString(strings.Trim(ma.printrep, " \t"))
	} else {
		return *ma
	}
}

// ----------------------------------------------------------------
func MlrvalCollapseWhitespace(ma *Mlrval) Mlrval {
	return MlrvalCollapseWhitespaceRegexp(ma, WhitespaceRegexp())
}

func MlrvalCollapseWhitespaceRegexp(ma *Mlrval, whitespaceRegexp *regexp.Regexp) Mlrval {
	if ma.mvtype == MT_STRING {
		return MlrvalFromString(whitespaceRegexp.ReplaceAllString(ma.printrep, " "))
	} else {
		return *ma
	}
}

func WhitespaceRegexp() *regexp.Regexp {
	return regexp.MustCompile("\\s+")
}

// ================================================================
func MlrvalToUpper(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_STRING {
		return MlrvalFromString(strings.ToUpper(ma.printrep))
	} else if ma.mvtype == MT_VOID {
		return *ma
	} else {
		return *ma
	}
}

func MlrvalToLower(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_STRING {
		return MlrvalFromString(strings.ToLower(ma.printrep))
	} else if ma.mvtype == MT_VOID {
		return *ma
	} else {
		return *ma
	}
}

func MlrvalCapitalize(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_STRING {
		if ma.printrep == "" {
			return *ma
		} else {
			runes := []rune(ma.printrep)
			rfirst := runes[0]
			rrest := runes[1:]
			sfirst := strings.ToUpper(string(rfirst))
			srest := string(rrest)
			return MlrvalFromString(sfirst + srest)
		}
	} else {
		return *ma
	}
}

// ----------------------------------------------------------------
func MlrvalCleanWhitespace(ma *Mlrval) Mlrval {
	temp := MlrvalCollapseWhitespaceRegexp(ma, WhitespaceRegexp())
	return MlrvalStrip(&temp)
}

// ================================================================
func MlrvalHexfmt(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_INT {
		return MlrvalFromString("0x" + strconv.FormatUint(uint64(ma.intval), 16))
	} else {
		return *ma
	}
}

// ----------------------------------------------------------------
func fmtnum_is(ma, mb *Mlrval) Mlrval {
	return MlrvalFromString(
		fmt.Sprintf(
			mb.printrep,
			ma.intval,
		),
	)
}

func fmtnum_fs(ma, mb *Mlrval) Mlrval {
	return MlrvalFromString(
		fmt.Sprintf(
			mb.printrep,
			ma.floatval,
		),
	)
}

func fmtnum_bs(ma, mb *Mlrval) Mlrval {
	return MlrvalFromString(
		fmt.Sprintf(
			mb.printrep,
			lib.BoolToInt(ma.boolval),
		),
	)
}

var fmtnum_dispositions = [MT_DIM][MT_DIM]BinaryFunc{
	//       .  ERROR   ABSENT VOID   STRING INT    FLOAT  BOOL   ARRAY  MAP
	/*ERROR  */ {_erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*ABSENT */ {_erro, _absn, _erro, _absn, _erro, _erro, _erro, _erro, _erro},
	/*VOID   */ {_erro, _absn, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*STRING */ {_erro, _absn, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*INT    */ {_erro, _absn, _erro, fmtnum_is, _erro, _erro, _erro, _erro, _erro},
	/*FLOAT  */ {_erro, _absn, _erro, fmtnum_fs, _erro, _erro, _erro, _erro, _erro},
	/*BOOL   */ {_erro, _absn, _erro, fmtnum_bs, _erro, _erro, _erro, _erro, _erro},
	/*ARRAY  */ {_erro, _absn, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
	/*MAP    */ {_erro, _absn, _erro, _erro, _erro, _erro, _erro, _erro, _erro},
}

func MlrvalFmtNum(ma, mb *Mlrval) Mlrval {
	return fmtnum_dispositions[ma.mvtype][mb.mvtype](ma, mb)
}
