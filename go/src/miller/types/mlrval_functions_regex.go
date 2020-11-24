package types

import (
	"regexp"
	"strings"
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
	re := regexp.MustCompile(mb.printrep)

	onFirst := true
	output := re.ReplaceAllStringFunc(ma.printrep, func(s string) string {
		if !onFirst {
			return s
		}
		onFirst = false
		return re.ReplaceAllString(s, mc.printrep)
	})

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
	re := regexp.MustCompile(mb.printrep)
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
	re := regexp.MustCompile(mb.printrep)
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
