package types

import (
	"strings"

	"mlr/internal/pkg/lib"
)

// BIF_ssub implements the ssub function -- no-frills string-replace, no
// regexes, no escape sequences.
func BIF_ssub(input1, input2, input3 *Mlrval) *Mlrval {
	if input1.IsErrorOrAbsent() {
		return input1
	}
	if input2.IsErrorOrAbsent() {
		return input2
	}
	if input3.IsErrorOrAbsent() {
		return input3
	}
	if !input1.IsStringOrVoid() {
		return MLRVAL_ERROR
	}
	if !input2.IsStringOrVoid() {
		return MLRVAL_ERROR
	}
	if !input3.IsStringOrVoid() {
		return MLRVAL_ERROR
	}
	return MlrvalFromString(
		strings.Replace(input1.printrep, input2.printrep, input3.printrep, 1),
	)
}

// BIF_sub implements the sub function, with support for regexes and regex captures
// of the form "\1" .. "\9".
//
// TODO: make a variant which allows compiling the regexp once and reusing it
// on each record. Likewise for other regex-using functions in this file.  But
// first, do a profiling run to see how much time would be saved, and if this
// precomputing+caching would be worthwhile.
func BIF_sub(input1, input2, input3 *Mlrval) *Mlrval {
	if input1.IsErrorOrAbsent() {
		return input1
	}
	if input2.IsErrorOrAbsent() {
		return input2
	}
	if input3.IsErrorOrAbsent() {
		return input3
	}
	if !input1.IsStringOrVoid() {
		return MLRVAL_ERROR
	}
	if !input2.IsStringOrVoid() {
		return MLRVAL_ERROR
	}
	if !input3.IsStringOrVoid() {
		return MLRVAL_ERROR
	}

	input := input1.printrep
	sregex := input2.printrep
	replacement := input3.printrep

	stringOutput := lib.RegexSub(input, sregex, replacement)
	return MlrvalFromString(stringOutput)
}

// BIF_gsub implements the gsub function, with support for regexes and regex captures
// of the form "\1" .. "\9".
func BIF_gsub(input1, input2, input3 *Mlrval) *Mlrval {
	if input1.IsErrorOrAbsent() {
		return input1
	}
	if input2.IsErrorOrAbsent() {
		return input2
	}
	if input3.IsErrorOrAbsent() {
		return input3
	}
	if !input1.IsStringOrVoid() {
		return MLRVAL_ERROR
	}
	if !input2.IsStringOrVoid() {
		return MLRVAL_ERROR
	}
	if !input3.IsStringOrVoid() {
		return MLRVAL_ERROR
	}

	input := input1.printrep
	sregex := input2.printrep
	replacement := input3.printrep

	stringOutput := lib.RegexGsub(input, sregex, replacement)
	return MlrvalFromString(stringOutput)
}

// BIF_string_matches_regexp implements the =~ operator, with support for
// setting regex-captures for later expressions to access using "\1" .. "\9".
func BIF_string_matches_regexp(input1, input2 *Mlrval) (retval *Mlrval, captures []string) {
	if !input1.IsLegit() {
		return input1, nil
	}
	if !input2.IsLegit() {
		return input2, nil
	}
	input1string := input1.String()
	if !input2.IsStringOrVoid() {
		return MLRVAL_ERROR, nil
	}

	boolOutput, captures := lib.RegexMatches(input1string, input2.printrep)
	return MlrvalFromBool(boolOutput), captures
}

// BIF_string_matches_regexp implements the !=~ operator.
func BIF_string_does_not_match_regexp(input1, input2 *Mlrval) (retval *Mlrval, captures []string) {
	output, captures := BIF_string_matches_regexp(input1, input2)
	if output.mvtype == MT_BOOL {
		return MlrvalFromBool(!output.boolval), captures
	} else {
		// else leave it as error, absent, etc.
		return output, captures
	}
}

func BIF_regextract(input1, input2 *Mlrval) *Mlrval {
	if !input1.IsString() {
		return MLRVAL_ERROR
	}
	if !input2.IsString() {
		return MLRVAL_ERROR
	}
	regex := lib.CompileMillerRegexOrDie(input2.printrep)
	match := regex.FindStringIndex(input1.printrep)
	if match != nil {
		return MlrvalFromString(input1.printrep[match[0]:match[1]])
	} else {
		return MLRVAL_ABSENT
	}
}

func BIF_regextract_or_else(input1, input2, input3 *Mlrval) *Mlrval {
	if !input1.IsString() {
		return MLRVAL_ERROR
	}
	if !input2.IsString() {
		return MLRVAL_ERROR
	}
	regex := lib.CompileMillerRegexOrDie(input2.printrep)
	match := regex.FindStringIndex(input1.printrep)
	if match != nil {
		return MlrvalFromString(input1.printrep[match[0]:match[1]])
	} else {
		return input3
	}
}
