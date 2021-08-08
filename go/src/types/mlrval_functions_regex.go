package types

import (
	"strings"

	"mlr/src/lib"
)

// MlrvalSsub implements the ssub function -- no-frills string-replace, no
// regexes, no escape sequences.
func MlrvalSsub(input1, input2, input3 *Mlrval) *Mlrval {
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
	return MlrvalPointerFromString(
		strings.Replace(input1.printrep, input2.printrep, input3.printrep, 1),
	)
}

// MlrvalSub implements the sub function, with support for regexes and regex captures
// of the form "\1" .. "\9".
//
// TODO: make a variant which allows compiling the regexp once and reusing it
// on each record
func MlrvalSub(input1, input2, input3 *Mlrval) *Mlrval {
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

	stringOutput := lib.RegexSubWithCaptures(input1.printrep, input2.printrep, input3.printrep)
	return MlrvalPointerFromString(stringOutput)
}

// MlrvalGsub implements the gsub function, with support for regexes and regex captures
// of the form "\1" .. "\9".
// TODO: make a variant which allows compiling the regexp once and reusing it
// on each record
func MlrvalGsub(input1, input2, input3 *Mlrval) *Mlrval {
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

	stringOutput := lib.RegexGsubWithCaptures(input1.printrep, input2.printrep, input3.printrep)
	return MlrvalPointerFromString(stringOutput)
}

// MlrvalStringMatchesRegexp implements the =~ operator, with support for
// setting regex-captures for later expressions to access using "\1" .. "\9".
//
// TODO: make a variant which allows compiling the regexp once and reusing it
// on each record
func MlrvalStringMatchesRegexp(input1, input2 *Mlrval) (retval *Mlrval, captures []string) {
	if !input1.IsLegit() {
		return input1, nil
	}
	if !input2.IsLegit() {
		return input2, nil
	}
	if !input1.IsStringOrVoid() {
		return MLRVAL_ERROR, nil
	}
	if !input2.IsStringOrVoid() {
		return MLRVAL_ERROR, nil
	}

	// TODO
	boolOutput, captures := lib.RegexMatches(input1.printrep, input2.printrep)
	return MlrvalPointerFromBool(boolOutput), captures
}

// MlrvalStringMatchesRegexp implements the !=~ operator.
func MlrvalStringDoesNotMatchRegexp(input1, input2 *Mlrval) (retval *Mlrval, captures []string) {
	output, captures := MlrvalStringMatchesRegexp(input1, input2)
	if output.mvtype == MT_BOOL {
		return MlrvalPointerFromBool(!output.boolval), captures
	} else {
		// else leave it as error, absent, etc.
		return output, captures
	}
}

// TODO: find a way to keep and stash a precompiled regex, somewhere in the CST ...
func MlrvalRegextract(input1, input2 *Mlrval) *Mlrval {
	if !input1.IsString() {
		return MLRVAL_ERROR
	}
	if !input2.IsString() {
		return MLRVAL_ERROR
	}
	regex := lib.CompileMillerRegexOrDie(input2.printrep)
	// TODO: See if we need FindStringIndex or FindStringSubmatch to distinguish from matching "".
	match := regex.FindString(input1.printrep)
	if match != "" {
		return MlrvalPointerFromString(match)
	} else {
		return MLRVAL_ABSENT
	}
}

func MlrvalRegextractOrElse(input1, input2, input3 *Mlrval) *Mlrval {
	if !input1.IsString() {
		return MLRVAL_ERROR
	}
	if !input2.IsString() {
		return MLRVAL_ERROR
	}
	regex := lib.CompileMillerRegexOrDie(input2.printrep)
	// TODO: See if we need FindStringIndex or FindStringSubmatch to distinguish from matching "".
	found := regex.FindString(input1.printrep)
	if found != "" {
		return MlrvalPointerFromString(found)
	} else {
		return input3
	}
}
