package lib

import (
	"fmt"
	"os"
	"strconv"
)

// Constructors

func MlrvalFromError() Mlrval {
	return Mlrval{
		MT_ERROR,
		"(error)", // xxx const somewhere
		true,
		0, 0.0, false,
	}
}

func MlrvalFromAbsent() Mlrval {
	return Mlrval{
		MT_ABSENT,
		"(absent)",
		true,
		0, 0.0, false,
	}
}

func MlrvalFromVoid() Mlrval {
	return Mlrval{
		MT_VOID,
		"(void)",
		true,
		0, 0.0, false,
	}
}

func MlrvalFromString(input string) Mlrval {
	return Mlrval{
		MT_STRING,
		input,
		true,
		0, 0.0, false,
	}
}

// xxx comment why two -- one for from parsed user data; other for from math ops
func MlrvalFromInt64String(input string) Mlrval {
	ival, ok := tryInt64FromString(input)
	// xxx comment assummption is input-string already deemed parseable so no error return
	if !ok {
		// xxx get file/line info here .......
		fmt.Fprintf(os.Stderr, "Internal coding error detected\n")
		os.Exit(1)
	}
	return Mlrval{
		MT_INT,
		input,
		true,
		ival,
		0.0,
		false,
	}
}

func MlrvalFromInt64(input int64) Mlrval {
	return Mlrval{
		MT_INT,
		"(bug-if-you-see-this)",
		false,
		input,
		0.0,
		false,
	}
}

func tryInt64FromString(input string) (int64, bool) {
	// xxx need to handle octal, hex, ......
	ival, err := strconv.ParseInt(input, 10, 64)
	if err == nil {
		return ival, true
	} else {
		return 0, false
	}
}

// xxx comment why two -- one for from parsed user data; other for from math ops
// xxx comment assummption is input-string already deemed parseable so no error return
func MlrvalFromFloat64String(input string) Mlrval {
	fval, ok := tryFloat64FromString(input)
	// xxx comment assummption is input-string already deemed parseable so no error return
	if !ok {
		// xxx get file/line info here .......
		fmt.Fprintf(os.Stderr, "Internal coding error detected\n")
		os.Exit(1)
	}
	return Mlrval{
		MT_FLOAT,
		input,
		true,
		0,
		fval,
		false,
	}
}

func MlrvalFromFloat64(input float64) Mlrval {
	return Mlrval{
		MT_FLOAT,
		"(bug-if-you-see-this)",
		false,
		0,
		input,
		false,
	}
}

func tryFloat64FromString(input string) (float64, bool) {
	ival, err := strconv.ParseFloat(input, 64)
	if err == nil {
		return ival, true
	} else {
		return 0, false
	}
}

func MlrvalFromTrue() Mlrval {
	return Mlrval{
		MT_BOOL,
		"true",
		true,
		0,
		0.0,
		true,
	}
}

func MlrvalFromFalse() Mlrval {
	return Mlrval{
		MT_BOOL,
		"false",
		true,
		0,
		0.0,
		false,
	}
}

func MlrvalFromBool(input bool) Mlrval {
	if input == true {
		return MlrvalFromTrue()
	} else {
		return MlrvalFromFalse()
	}
}

func MlrvalFromBoolString(input string) Mlrval {
	if input == "true" {
		return MlrvalFromTrue()
	} else {
		return MlrvalFromFalse()
	}
	// else panic
}

func tryBoolFromBoolString(input string) (bool, bool) {
	if input == "true" {
		return true, true
	} else if input == "false" {
		return false, true
	} else {
		return false, false
	}
}

func MlrvalFromInferredType(input string) Mlrval {
	// xxx the parsing has happened so stash it ...
	// xxx emphasize the invariant that a non-invalid printrep always
	// matches the nval ...
	_, iok := tryInt64FromString(input)
	if iok {
		return MlrvalFromInt64String(input)
	}

	_, fok := tryFloat64FromString(input)
	if fok {
		return MlrvalFromFloat64String(input)
	}

	_, bok := tryBoolFromBoolString(input)
	if bok {
		return MlrvalFromBoolString(input)
	}

	return MlrvalFromString(input)
}
