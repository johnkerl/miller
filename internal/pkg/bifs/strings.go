package bifs

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
)

// ================================================================
func BIF_strlen(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsStringOrVoid() {
		return mlrval.ERROR
	} else {
		return mlrval.FromInt(lib.UTF8Strlen(input1.AcquireStringValue()))
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

	sliceIsEmpty, absentOrError, lowerZindex, upperZindex := MillerSliceAccess(input2, input3, strlen, false)

	if sliceIsEmpty {
		return mlrval.VOID
	}
	if absentOrError != nil {
		return absentOrError
	}

	// Note Golang slice indices are 0-up, and the 1st index is inclusive
	// while the 2nd is exclusive. For Miller, indices are 1-up and both
	// are inclusive.
	return mlrval.FromString(string(runes[lowerZindex : upperZindex+1]))
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

	sliceIsEmpty, absentOrError, lowerZindex, upperZindex := MillerSliceAccess(input2, input3, strlen, true)

	if sliceIsEmpty {
		return mlrval.VOID
	}
	if absentOrError != nil {
		return absentOrError
	}

	// Note Golang slice indices are 0-up, and the 1st index is inclusive
	// while the 2nd is exclusive. For Miller, indices are 1-up and both
	// are inclusive.
	return mlrval.FromString(string(runes[lowerZindex : upperZindex+1]))
}

// ================================================================
// index(string, substring) returns the index of substring within string (if found), or -1 if not
// found.

func BIF_index(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsAbsent() {
		return mlrval.ABSENT
	}
	if input1.IsError() {
		return mlrval.ERROR
	}
	sinput1 := input1.String()
	sinput2 := input2.String()

	// Handle UTF-8 correctly, since Go's strings.Index counts bytes
	iindex := strings.Index(sinput1, sinput2)
	if iindex < 0 {
		return mlrval.FromInt(int64(iindex))
	}

	// Go indices are 0-up; Miller indices are 1-up.
	return mlrval.FromInt(lib.UTF8Strlen(sinput1[:iindex]) + 1)
}

// ================================================================
// contains(string, substring) returns true if string contains substring, else false.

func BIF_contains(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsAbsent() {
		return mlrval.ABSENT
	}
	if input1.IsError() {
		return mlrval.ERROR
	}

	return mlrval.FromBool(strings.Contains(input1.String(), input2.String()))
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
	maxLength := int(input2.AcquireIntValue())
	if oldLength <= maxLength {
		return input1
	} else {
		return mlrval.FromString(string(runes[0:maxLength]))
	}
}

// ================================================================
func BIF_leftpad(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsErrorOrAbsent() {
		return input1
	}
	if input2.IsErrorOrAbsent() {
		return input2
	}
	if input3.IsErrorOrAbsent() {
		return input3
	}

	if !input2.IsInt() {
		return mlrval.ERROR
	}

	inputString := input1.String()
	padString := input3.String()

	inputLength := lib.UTF8Strlen(inputString)
	padLength := lib.UTF8Strlen(padString)
	targetLength := input2.AcquireIntValue()
	outputLength := inputLength

	var buffer bytes.Buffer
	for outputLength+padLength <= targetLength {
		buffer.WriteString(padString)
		outputLength += padLength
	}
	buffer.WriteString(inputString)

	return mlrval.FromString(buffer.String())
}

func BIF_rightpad(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsErrorOrAbsent() {
		return input1
	}
	if input2.IsErrorOrAbsent() {
		return input2
	}
	if input3.IsErrorOrAbsent() {
		return input3
	}

	if !input2.IsInt() {
		return mlrval.ERROR
	}

	inputString := input1.String()
	padString := input3.String()

	inputLength := lib.UTF8Strlen(inputString)
	padLength := lib.UTF8Strlen(padString)
	targetLength := input2.AcquireIntValue()
	outputLength := inputLength

	var buffer bytes.Buffer
	buffer.WriteString(inputString)
	for outputLength+padLength <= targetLength {
		buffer.WriteString(padString)
		outputLength += padLength
	}

	return mlrval.FromString(buffer.String())
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
	return BIF_collapse_whitespace_regexp(input1, _whitespace_regexp)
}

func BIF_collapse_whitespace_regexp(input1 *mlrval.Mlrval, whitespaceRegexp *regexp.Regexp) *mlrval.Mlrval {
	if input1.IsString() {
		return mlrval.FromString(whitespaceRegexp.ReplaceAllString(input1.AcquireStringValue(), " "))
	} else {
		return input1
	}
}

var _whitespace_regexp = regexp.MustCompile(`\s+`)

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
			input1, _whitespace_regexp,
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

// unformat("{}:{}:{}",  "1:2:3")    gives [1, 2]
// unformat("{}h{}m{}s", "3h47m22s") gives [3, 47, 22]
func BIF_unformat(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return bif_unformat_aux(input1, input2, true)
}
func BIF_unformatx(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	return bif_unformat_aux(input1, input2, false)
}

func bif_unformat_aux(input1, input2 *mlrval.Mlrval, inferTypes bool) *mlrval.Mlrval {
	template, ok1 := input1.GetStringValue()
	if !ok1 {
		return mlrval.ERROR
	}
	input, ok2 := input2.GetStringValue()
	if !ok2 {
		return mlrval.ERROR
	}

	templatePieces := strings.Split(template, "{}")
	output := mlrval.FromEmptyArray()

	// template "{}h{}m{}s"
	// input    "12h34m56s"
	// templatePieces   ["", "h", "m", "s"]

	remaining := input

	if !strings.HasPrefix(remaining, templatePieces[0]) {
		return mlrval.ERROR
	}
	remaining = remaining[len(templatePieces[0]):]
	templatePieces = templatePieces[1:]

	n := len(templatePieces)
	for i, templatePiece := range templatePieces {

		var index int
		if i == n-1 && templatePiece == "" {
			// strings.Index("", ...) will match the *start* of what's
			// remaining, whereas we want it to match the end.
			index = len(remaining)
		} else {
			index = strings.Index(remaining, templatePiece)
			if index < 0 {
				return mlrval.ERROR
			}
		}

		inputPiece := remaining[:index]
		remaining = remaining[index+len(templatePiece):]
		if inferTypes {
			output.ArrayAppend(mlrval.FromInferredType(inputPiece))
		} else {
			output.ArrayAppend(mlrval.FromString(inputPiece))
		}
	}

	return output
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
	if input1.IsArray() || input1.IsMap() {
		return recurseBinaryFuncOnInput1(BIF_fmtnum, input1, input2)
	} else {
		return fmtnum_dispositions[input1.Type()][input2.Type()](input1, input2)
	}
}

func BIF_fmtifnum(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsArray() || input1.IsMap() {
		return recurseBinaryFuncOnInput1(BIF_fmtifnum, input1, input2)
	} else {
		output := fmtnum_dispositions[input1.Type()][input2.Type()](input1, input2)
		if output.IsError() {
			return input1
		} else {
			return output
		}
	}
}

func BIF_latin1_to_utf8(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsArray() || input1.IsMap() {
		return recurseUnaryFuncOnInput1(BIF_latin1_to_utf8, input1)
	} else if input1.IsString() {
		output, err := lib.TryLatin1ToUTF8(input1.String())
		if err != nil {
			// Somewhat arbitrary design decision
			// return input1
			return mlrval.ERROR
		} else {
			return mlrval.FromString(output)
		}
	} else {
		return input1
	}
}

func BIF_utf8_to_latin1(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsArray() || input1.IsMap() {
		return recurseUnaryFuncOnInput1(BIF_utf8_to_latin1, input1)
	} else if input1.IsString() {
		output, err := lib.TryUTF8ToLatin1(input1.String())
		if err != nil {
			// Somewhat arbitrary design decision
			// return input1
			return mlrval.ERROR
		} else {
			return mlrval.FromString(output)
		}
	} else {
		return input1
	}
}
