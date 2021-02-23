package types

import (
	"strings"

	"miller/src/lib"
)

// ================================================================
func MlrvalSsub(input1, input2, input3 *Mlrval) *Mlrval {
	if input1.IsErrorOrAbsent() {
		return input1
	} else if input2.IsErrorOrAbsent() {
		return input2
	} else if input3.IsErrorOrAbsent() {
		return input3
	} else if !input1.IsStringOrVoid() {
		return MLRVAL_ERROR
	} else if !input2.IsStringOrVoid() {
		return MLRVAL_ERROR
	} else if !input3.IsStringOrVoid() {
		return MLRVAL_ERROR
	} else {
		return MlrvalPointerFromString(
			strings.Replace(input1.printrep, input2.printrep, input3.printrep, 1),
		)
	}
}

// ================================================================
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
	// TODO: better exception-handling
	re := lib.CompileMillerRegexOrDie(input2.printrep)

	replacement := lib.RegexReplaceOnce(re, input1.printrep, input3.printrep)

	return MlrvalPointerFromString(replacement)
}

// ================================================================
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
	// TODO: better exception-handling
	re := lib.CompileMillerRegexOrDie(input2.printrep)
	return MlrvalPointerFromString(
		re.ReplaceAllString(input1.printrep, input3.printrep),
	)
}

// ================================================================
// TODO: make a variant which allows compiling the regexp once and reusing it
// on each record
func MlrvalStringMatchesRegexp(input1, input2 *Mlrval) *Mlrval {
	if !input1.IsLegit() {
		return input1
	}
	if !input2.IsLegit() {
		return input2
	}
	if !input1.IsStringOrVoid() {
		return MLRVAL_ERROR
	}
	if !input2.IsStringOrVoid() {
		return MLRVAL_ERROR
	}
	// TODO: better exception-handling
	re := lib.CompileMillerRegexOrDie(input2.printrep)
	return MlrvalPointerFromBool(
		re.MatchString(input1.printrep),
	)
}

func MlrvalStringDoesNotMatchRegexp(input1, input2 *Mlrval) *Mlrval {
	output := MlrvalStringMatchesRegexp(input1, input2)
	if output.mvtype == MT_BOOL {
		return MlrvalPointerFromBool(!output.boolval)
	} else {
		// else leave it as error, absent, etc.
		return output
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
