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
