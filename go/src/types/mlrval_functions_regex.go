package types

import (
	"strings"

	"miller/src/lib"
)

// ================================================================
func MlrvalSsub(output, input1, input2, input3 *Mlrval) {
	if input1.IsErrorOrAbsent() {
		output.CopyFrom(input1)
	} else if input2.IsErrorOrAbsent() {
		output.CopyFrom(input2)
	} else if input3.IsErrorOrAbsent() {
		output.CopyFrom(input3)
	} else if !input1.IsStringOrVoid() {
		output.SetFromError()
	} else if !input2.IsStringOrVoid() {
		output.SetFromError()
	} else if !input3.IsStringOrVoid() {
		output.SetFromError()
	} else {
		output.SetFromString(
			strings.Replace(input1.printrep, input2.printrep, input3.printrep, 1),
		)
	}
}

// ================================================================
// TODO: make a variant which allows compiling the regexp once and reusing it
// on each record
func MlrvalSub(output, input1, input2, input3 *Mlrval) {
	if input1.IsErrorOrAbsent() {
		output.CopyFrom(input1)
		return
	}
	if input2.IsErrorOrAbsent() {
		output.CopyFrom(input2)
		return
	}
	if input3.IsErrorOrAbsent() {
		output.CopyFrom(input3)
		return
	}
	if !input1.IsStringOrVoid() {
		output.SetFromError()
		return
	}
	if !input2.IsStringOrVoid() {
		output.SetFromError()
		return
	}
	if !input3.IsStringOrVoid() {
		output.SetFromError()
		return
	}
	// TODO: better exception-handling
	re := lib.CompileMillerRegexOrDie(input2.printrep)

	replacement := lib.RegexReplaceOnce(re, input1.printrep, input3.printrep)

	output.SetFromString(replacement)
}

// ================================================================
// TODO: make a variant which allows compiling the regexp once and reusing it
// on each record
func MlrvalGsub(output, input1, input2, input3 *Mlrval) {
	if input1.IsErrorOrAbsent() {
		output.CopyFrom(input1)
		return
	}
	if input2.IsErrorOrAbsent() {
		output.CopyFrom(input2)
		return
	}
	if input3.IsErrorOrAbsent() {
		output.CopyFrom(input3)
		return
	}
	if !input1.IsStringOrVoid() {
		output.SetFromError()
		return
	}
	if !input2.IsStringOrVoid() {
		output.SetFromError()
		return
	}
	if !input3.IsStringOrVoid() {
		output.SetFromError()
		return
	}
	// TODO: better exception-handling
	re := lib.CompileMillerRegexOrDie(input2.printrep)
	output.SetFromString(
		re.ReplaceAllString(input1.printrep, input3.printrep),
	)
}

// ================================================================
// TODO: make a variant which allows compiling the regexp once and reusing it
// on each record
func MlrvalStringMatchesRegexp(output, input1, input2 *Mlrval) {
	if !input1.IsLegit() {
		output.CopyFrom(input1)
		return
	}
	if !input2.IsLegit() {
		output.CopyFrom(input2)
		return
	}
	if !input1.IsStringOrVoid() {
		output.SetFromError()
		return
	}
	if !input2.IsStringOrVoid() {
		output.SetFromError()
		return
	}
	// TODO: better exception-handling
	re := lib.CompileMillerRegexOrDie(input2.printrep)
	output.SetFromBool(
		re.MatchString(input1.printrep),
	)
}

func MlrvalStringDoesNotMatchRegexp(output, input1, input2 *Mlrval) {
	MlrvalStringMatchesRegexp(output, input1, input2)
	if output.mvtype == MT_BOOL {
		output.SetFromBool(!output.boolval)
	}
	// else leave it as error, absent, etc.
}

// TODO: find a way to keep and stash a precompiled regex, somewhere in the CST ...
func MlrvalRegextract(output, input1, input2 *Mlrval) {
	if !input1.IsString() {
		output.SetFromError()
		return
	}
	if !input2.IsString() {
		output.SetFromError()
		return
	}
	regex := lib.CompileMillerRegexOrDie(input2.printrep)
	// TODO: See if we need FindStringIndex or FindStringSubmatch to distinguish from matching "".
	match := regex.FindString(input1.printrep)
	if match != "" {
		output.SetFromString(match)
	} else {
		output.SetFromAbsent()
	}
}

func MlrvalRegextractOrElse(output, input1, input2, input3 *Mlrval) {
	if !input1.IsString() {
		output.SetFromError()
		return
	}
	if !input2.IsString() {
		output.SetFromError()
		return
	}
	regex := lib.CompileMillerRegexOrDie(input2.printrep)
	// TODO: See if we need FindStringIndex or FindStringSubmatch to distinguish from matching "".
	found := regex.FindString(input1.printrep)
	if found != "" {
		output.SetFromString(found)
	} else {
		output.CopyFrom(input3)
	}
}
