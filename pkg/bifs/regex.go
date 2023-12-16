package bifs

import (
	"strings"

	"github.com/johnkerl/miller/pkg/lib"
	"github.com/johnkerl/miller/pkg/mlrval"
)

// BIF_ssub implements the ssub function -- no-frills string-replace, no
// regexes, no escape sequences.
func BIF_ssub(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	return bif_ssub_gssub(input1, input2, input3, false, "ssub")
}

// BIF_gssub implements the gssub function -- no-frills string-replace, no
// regexes, no escape sequences.
func BIF_gssub(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	return bif_ssub_gssub(input1, input2, input3, true, "gssub")
}

// bif_ssub_gssub is shared code for BIF_ssub and BIF_gssub.
func bif_ssub_gssub(input1, input2, input3 *mlrval.Mlrval, doAll bool, funcname string) *mlrval.Mlrval {
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
		return mlrval.FromNotStringError(funcname, input1)
	}
	if !input2.IsStringOrVoid() {
		return mlrval.FromNotStringError(funcname, input2)
	}
	if !input3.IsStringOrVoid() {
		return mlrval.FromNotStringError(funcname, input3)
	}
	if doAll {
		return mlrval.FromString(
			strings.ReplaceAll(input1.AcquireStringValue(), input2.AcquireStringValue(), input3.AcquireStringValue()),
		)
	} else {
		return mlrval.FromString(
			strings.Replace(input1.AcquireStringValue(), input2.AcquireStringValue(), input3.AcquireStringValue(), 1),
		)
	}
}

// BIF_sub implements the sub function, with support for regexes and regex captures
// of the form "\1" .. "\9".
//
// TODO: make a variant which allows compiling the regexp once and reusing it
// on each record. Likewise for other regex-using functions in this file.  But
// first, do a profiling run to see how much time would be saved, and if this
// precomputing+caching would be worthwhile.
func BIF_sub(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
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
		return mlrval.FromNotStringError("sub", input1)
	}
	if !input2.IsStringOrVoid() {
		return mlrval.FromNotStringError("sub", input2)
	}
	if !input3.IsStringOrVoid() {
		return mlrval.FromNotStringError("sub", input3)
	}

	input := input1.AcquireStringValue()
	sregex := input2.AcquireStringValue()
	replacement := input3.AcquireStringValue()

	stringOutput := lib.RegexStringSub(input, sregex, replacement)
	return mlrval.FromString(stringOutput)
}

// BIF_gsub implements the gsub function, with support for regexes and regex captures
// of the form "\1" .. "\9".
func BIF_gsub(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
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
		return mlrval.FromNotStringError("gsub", input1)
	}
	if !input2.IsStringOrVoid() {
		return mlrval.FromNotStringError("gsub", input2)
	}
	if !input3.IsStringOrVoid() {
		return mlrval.FromNotStringError("gsub", input3)
	}

	input := input1.AcquireStringValue()
	sregex := input2.AcquireStringValue()
	replacement := input3.AcquireStringValue()

	stringOutput := lib.RegexStringGsub(input, sregex, replacement)
	return mlrval.FromString(stringOutput)
}

func BIF_match(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsLegit() {
		return mlrval.FromNotStringError("match", input1) // TODO: CHANGE FLAVOR
	}
	if !input2.IsLegit() {
		return mlrval.FromNotStringError("match", input2) // TODO: CHANGE FLAVOR
	}
	input1string := input1.String()
	if !input2.IsStringOrVoid() {
		return mlrval.FromNotStringError("match", input2)
	}

	boolOutput := lib.RegexStringMatchSimple(input1string, input2.AcquireStringValue())

	return mlrval.FromBool(boolOutput)
}

// TODO: WRITE ME
func BIF_matchx(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsLegit() {
		return mlrval.FromNotStringError("match", input1) // TODO: CHANGE FLAVOR
	}
	if !input2.IsLegit() {
		return mlrval.FromNotStringError("match", input2) // TODO: CHANGE FLAVOR
	}
	input1string := input1.String()
	if !input2.IsStringOrVoid() {
		return mlrval.FromNotStringError("match", input2)
	}

	boolOutput, captures, starts, ends := lib.RegexStringMatchWithMapResults(input1string, input2.AcquireStringValue())

	results := mlrval.NewMlrmap()
	results.PutReference("matched", mlrval.FromBool(boolOutput))

	// XXX assert all lengths equal

	captures_array := make([]*mlrval.Mlrval, len(captures))
	starts_array := make([]*mlrval.Mlrval, len(captures))
	ends_array := make([]*mlrval.Mlrval, len(captures))

	if len(captures) > 0 {
		for i, _ := range captures {
			if i == 0 {
				results.PutReference("full_capture", mlrval.FromString(captures[i]))
			} else {
				captures_array[i] = mlrval.FromString(captures[i])
			}
		}

		starts_array := make([]*mlrval.Mlrval, len(starts))
		for i, _ := range starts {
			if i == 0 {
				results.PutReference("full_start", mlrval.FromInt(int64(starts[i])))
			} else {
				starts_array[i] = mlrval.FromInt(int64(starts[i]))
			}
		}

		ends_array := make([]*mlrval.Mlrval, len(ends))
		for i, _ := range ends {
			if i == 0 {
				results.PutReference("full_end", mlrval.FromInt(int64(ends[i])))
			} else {
				ends_array[i] = mlrval.FromInt(int64(ends[i]))
			}
		}

		if len(captures) > 1 {
			results.PutReference("captures", mlrval.FromArray(captures_array[1:]))
			results.PutReference("starts", mlrval.FromArray(starts_array[1:]))
			results.PutReference("ends", mlrval.FromArray(ends_array[1:]))
		}
	} else {
		results.PutReference("captures", mlrval.FromArray(captures_array))
		results.PutReference("starts", mlrval.FromArray(starts_array))
		results.PutReference("ends", mlrval.FromArray(ends_array))
	}

	return mlrval.FromMap(results)
}

// BIF_string_matches_regexp implements the =~ operator, with support for
// setting regex-captures for later expressions to access using "\1" .. "\9".
func BIF_string_matches_regexp(input1, input2 *mlrval.Mlrval) (retval *mlrval.Mlrval, captures []string) {
	if !input1.IsLegit() {
		return input1, nil
	}
	if !input2.IsLegit() {
		return input2, nil
	}
	input1string := input1.String()
	if !input2.IsStringOrVoid() {
		return mlrval.FromNotStringError("=~", input2), nil
	}

	boolOutput, captures := lib.RegexStringMatchWithCaptures(input1string, input2.AcquireStringValue())
	return mlrval.FromBool(boolOutput), captures
}

// BIF_string_matches_regexp implements the !=~ operator.
func BIF_string_does_not_match_regexp(input1, input2 *mlrval.Mlrval) (retval *mlrval.Mlrval, captures []string) {
	output, captures := BIF_string_matches_regexp(input1, input2)
	if output.IsBool() {
		return mlrval.FromBool(!output.AcquireBoolValue()), captures
	} else {
		// else leave it as error, absent, etc.
		return output, captures
	}
}

func BIF_regextract(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsString() {
		return mlrval.FromNotStringError("!=~", input1)
	}
	if !input2.IsString() {
		return mlrval.FromNotStringError("!=~", input2)
	}
	regex := lib.CompileMillerRegexOrDie(input2.AcquireStringValue())
	match := regex.FindStringIndex(input1.AcquireStringValue())
	if match != nil {
		return mlrval.FromString(input1.AcquireStringValue()[match[0]:match[1]])
	} else {
		return mlrval.ABSENT
	}
}

func BIF_regextract_or_else(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsString() {
		return mlrval.FromNotStringError("regextract_or_else", input1)
	}
	if !input2.IsString() {
		return mlrval.FromNotStringError("regextract_or_else", input2)
	}
	regex := lib.CompileMillerRegexOrDie(input2.AcquireStringValue())
	match := regex.FindStringIndex(input1.AcquireStringValue())
	if match != nil {
		return mlrval.FromString(input1.AcquireStringValue()[match[0]:match[1]])
	} else {
		return input3
	}
}
