package types

import (
	"strings"

	"miller/src/lib"
)

// ================================================================
func MlrvalSsub(input1, input2, input3 *Mlrval) Mlrval {
	if input1.IsErrorOrAbsent() {
		return *input1
	}
	if input2.IsErrorOrAbsent() {
		return *input2
	}
	if input3.IsErrorOrAbsent() {
		return *input3
	}
	if !input1.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !input2.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !input3.IsStringOrVoid() {
		return MlrvalFromError()
	}
	return MlrvalFromString(
		strings.Replace(input1.printrep, input2.printrep, input3.printrep, 1),
	)
}

// ================================================================
// TODO: make a variant which allows compiling the regexp once and reusing it
// on each record
func MlrvalSub(input1, input2, input3 *Mlrval) Mlrval {
	if input1.IsErrorOrAbsent() {
		return *input1
	}
	if input2.IsErrorOrAbsent() {
		return *input2
	}
	if input3.IsErrorOrAbsent() {
		return *input3
	}
	if !input1.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !input2.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !input3.IsStringOrVoid() {
		return MlrvalFromError()
	}
	// TODO: better exception-handling
	re := lib.CompileMillerRegexOrDie(input2.printrep)

	output := lib.RegexReplaceOnce(re, input1.printrep, input3.printrep)

	return MlrvalFromString(output)
}

// ================================================================
// TODO: make a variant which allows compiling the regexp once and reusing it
// on each record
func MlrvalGsub(input1, input2, input3 *Mlrval) Mlrval {
	if input1.IsErrorOrAbsent() {
		return *input1
	}
	if input2.IsErrorOrAbsent() {
		return *input2
	}
	if input3.IsErrorOrAbsent() {
		return *input3
	}
	if !input1.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !input2.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !input3.IsStringOrVoid() {
		return MlrvalFromError()
	}
	// TODO: better exception-handling
	re := lib.CompileMillerRegexOrDie(input2.printrep)
	return MlrvalFromString(
		re.ReplaceAllString(input1.printrep, input3.printrep),
	)
}

// ================================================================
// TODO: make a variant which allows compiling the regexp once and reusing it
// on each record
func MlrvalStringMatchesRegexp(input1, input2 *Mlrval) Mlrval {
	if !input1.IsLegit() {
		return *input1
	}
	if !input2.IsLegit() {
		return *input2
	}
	if !input1.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !input2.IsStringOrVoid() {
		return MlrvalFromError()
	}
	// TODO: better exception-handling
	re := lib.CompileMillerRegexOrDie(input2.printrep)
	return MlrvalFromBool(
		re.MatchString(input1.printrep),
	)
}

func MlrvalStringDoesNotMatchRegexp(input1, input2 *Mlrval) Mlrval {
	matches := MlrvalStringMatchesRegexp(input1, input2)
	if matches.mvtype == MT_BOOL {
		return MlrvalFromBool(!matches.boolval)
	} else {
		return matches // error, absent, etc.
	}
}

// TODO: find a way to keep and stash a precompiled regex, somewhere in the CST ...
func MlrvalRegextract(input1, input2 *Mlrval) Mlrval {
	if !input1.IsString() {
		return MlrvalFromError()
	}
	if !input2.IsString() {
		return MlrvalFromError()
	}
	regex := lib.CompileMillerRegexOrDie(input2.printrep)
	// TODO: See if we need FindStringIndex or FindStringSubmatch to distinguish from matching "".
	output := regex.FindString(input1.printrep)
	if output != "" {
		return MlrvalFromString(output)
	} else {
		return MlrvalFromAbsent()
	}
}

func MlrvalRegextractOrElse(input1, input2, input3 *Mlrval) Mlrval {
	if !input1.IsString() {
		return MlrvalFromError()
	}
	if !input2.IsString() {
		return MlrvalFromError()
	}
	regex := lib.CompileMillerRegexOrDie(input2.printrep)
	// TODO: See if we need FindStringIndex or FindStringSubmatch to distinguish from matching "".
	output := regex.FindString(input1.printrep)
	if output != "" {
		return MlrvalFromString(output)
	} else {
		return *input3
	}
}
