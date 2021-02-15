package types

import (
	"strings"

	"miller/src/lib"
)

// ================================================================
func MlrvalSsub(ma, mb, mc *Mlrval) Mlrval {
	if ma.IsErrorOrAbsent() {
		return *ma
	}
	if mb.IsErrorOrAbsent() {
		return *mb
	}
	if mc.IsErrorOrAbsent() {
		return *mc
	}
	if !ma.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !mb.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !mc.IsStringOrVoid() {
		return MlrvalFromError()
	}
	return MlrvalFromString(
		strings.Replace(ma.printrep, mb.printrep, mc.printrep, 1),
	)
}

// ================================================================
// TODO: make a variant which allows compiling the regexp once and reusing it
// on each record
func MlrvalSub(ma, mb, mc *Mlrval) Mlrval {
	if ma.IsErrorOrAbsent() {
		return *ma
	}
	if mb.IsErrorOrAbsent() {
		return *mb
	}
	if mc.IsErrorOrAbsent() {
		return *mc
	}
	if !ma.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !mb.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !mc.IsStringOrVoid() {
		return MlrvalFromError()
	}
	// TODO: better exception-handling
	re := lib.CompileMillerRegexOrDie(mb.printrep)

	output := lib.RegexReplaceOnce(re, ma.printrep, mc.printrep)

	return MlrvalFromString(output)
}

// ================================================================
// TODO: make a variant which allows compiling the regexp once and reusing it
// on each record
func MlrvalGsub(ma, mb, mc *Mlrval) Mlrval {
	if ma.IsErrorOrAbsent() {
		return *ma
	}
	if mb.IsErrorOrAbsent() {
		return *mb
	}
	if mc.IsErrorOrAbsent() {
		return *mc
	}
	if !ma.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !mb.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !mc.IsStringOrVoid() {
		return MlrvalFromError()
	}
	// TODO: better exception-handling
	re := lib.CompileMillerRegexOrDie(mb.printrep)
	return MlrvalFromString(
		re.ReplaceAllString(ma.printrep, mc.printrep),
	)
}

// ================================================================
// TODO: make a variant which allows compiling the regexp once and reusing it
// on each record
func MlrvalStringMatchesRegexp(ma, mb *Mlrval) Mlrval {
	if !ma.IsLegit() {
		return *ma
	}
	if !mb.IsLegit() {
		return *mb
	}
	if !ma.IsStringOrVoid() {
		return MlrvalFromError()
	}
	if !mb.IsStringOrVoid() {
		return MlrvalFromError()
	}
	// TODO: better exception-handling
	re := lib.CompileMillerRegexOrDie(mb.printrep)
	return MlrvalFromBool(
		re.MatchString(ma.printrep),
	)
}

func MlrvalStringDoesNotMatchRegexp(ma, mb *Mlrval) Mlrval {
	matches := MlrvalStringMatchesRegexp(ma, mb)
	if matches.mvtype == MT_BOOL {
		return MlrvalFromBool(!matches.boolval)
	} else {
		return matches // error, absent, etc.
	}
}

// TODO: find a way to keep and stash a precompiled regex, somewhere in the CST ...
func MlrvalRegextract(ma, mb *Mlrval) Mlrval {
	if !ma.IsString() {
		return MlrvalFromError()
	}
	if !mb.IsString() {
		return MlrvalFromError()
	}
	regex := lib.CompileMillerRegexOrDie(mb.printrep)
	// TODO: See if we need FindStringIndex or FindStringSubmatch to distinguish from matching "".
	output := regex.FindString(ma.printrep)
	if output != "" {
		return MlrvalFromString(output)
	} else {
		return MlrvalFromAbsent()
	}
}

func MlrvalRegextractOrElse(ma, mb, mc *Mlrval) Mlrval {
	if !ma.IsString() {
		return MlrvalFromError()
	}
	if !mb.IsString() {
		return MlrvalFromError()
	}
	regex := lib.CompileMillerRegexOrDie(mb.printrep)
	// TODO: See if we need FindStringIndex or FindStringSubmatch to distinguish from matching "".
	output := regex.FindString(ma.printrep)
	if output != "" {
		return MlrvalFromString(output)
	} else {
		return *mc
	}
}
