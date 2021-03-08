package types

import (
	"fmt"
	"regexp"
	"time"

	"github.com/lestrrat-go/strftime"
	"github.com/pbnjay/strptime"

	"miller/src/lib"
)

const ISO8601_TIME_FORMAT = "%Y-%m-%dT%H:%M:%SZ"

var ptr_ISO8601_TIME_FORMAT = MlrvalPointerFromString("%Y-%m-%dT%H:%M:%SZ")

// ================================================================
func MlrvalSystime() *Mlrval {
	return MlrvalPointerFromFloat64(
		float64(time.Now().UnixNano()) / 1.0e9,
	)
}
func MlrvalSystimeInt() *Mlrval {
	return MlrvalPointerFromInt(int(time.Now().Unix()))
}

var startTime float64

func init() {
	startTime = float64(time.Now().UnixNano()) / 1.0e9
}
func MlrvalUptime() *Mlrval {
	return MlrvalPointerFromFloat64(
		float64(time.Now().UnixNano())/1.0e9 - startTime,
	)
}

// ================================================================
func MlrvalSec2GMTUnary(input1 *Mlrval) *Mlrval {
	if input1.mvtype == MT_FLOAT {
		return MlrvalPointerFromString(lib.Sec2GMT(input1.floatval, 0))
	} else if input1.mvtype == MT_INT {
		return MlrvalPointerFromString(lib.Sec2GMT(float64(input1.intval), 0))
	} else {
		return input1
	}
}

// ----------------------------------------------------------------
func MlrvalSec2GMTBinary(input1, input2 *Mlrval) *Mlrval {
	if input2.mvtype != MT_INT {
		return MLRVAL_ERROR
	} else if input1.mvtype == MT_FLOAT {
		return MlrvalPointerFromString(lib.Sec2GMT(input1.floatval, int(input2.intval)))
	} else if input1.mvtype == MT_INT {
		return MlrvalPointerFromString(lib.Sec2GMT(float64(input1.intval), int(input2.intval)))
	} else {
		return input1
	}
}

// ----------------------------------------------------------------
func MlrvalSec2GMTDate(input1 *Mlrval) *Mlrval {
	if input1.mvtype == MT_INT || input1.mvtype == MT_FLOAT {
		return MlrvalStrftime(input1, MlrvalPointerFromString("%Y-%m-%d"))
	} else {
		return input1
	}
}

// ================================================================
// Argument 1 is int/float seconds since the epoch.
// Argument 2 is format string like "%Y-%m-%d %H:%M:%S".

var extensionRegex = regexp.MustCompile("([1-9])S")

func MlrvalStrftime(input1, input2 *Mlrval) *Mlrval {
	epochSeconds, ok := input1.GetNumericToFloatValue()
	if !ok {
		return MLRVAL_ERROR
	}
	if input2.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}

	// Convert argument1 from float seconds since the epoch to a Go time.
	inputTime := lib.EpochSecondsToTime(epochSeconds)

	// Convert argument 2 to a strftime format string.
	//
	// Miller fractional-second formats are like "%6S", and were so in the C
	// implementation. However, in the strftime package we're using in the Go
	// port, extension-formats are only a single byte so we need to rewrite
	// them to "%6".
	formatString := extensionRegex.ReplaceAllString(input2.printrep, "$1")

	formatter, err := strftime.New(formatString, strftimeExtensions)
	if err != nil {
		return MLRVAL_ERROR
	}

	outputString := formatter.FormatString(inputTime)

	return MlrvalPointerFromString(outputString)
}

// ----------------------------------------------------------------
// This is support for %1S .. %9S in format strings, using github.com/lestrrat-go/strftime.

var strftimeExtensions strftime.Option

// This is a helper function for the appenders below, which let people get
// 1..9 decimal places in the seconds of their strftime format strings.
func specificationHelper(b []byte, t time.Time, sprintfFormat string, quotient int) []byte {
	seconds := int(t.Second())
	fractional := int(t.Nanosecond() / quotient)
	secondsString := fmt.Sprintf("%02d", seconds)
	b = append(b, secondsString...)
	b = append(b, '.')
	fractionalString := fmt.Sprintf(sprintfFormat, fractional)
	b = append(b, fractionalString...)
	return b
}

func init() {
	appender1 := strftime.AppendFunc(func(b []byte, t time.Time) []byte {
		return specificationHelper(b, t, "%01d", 100000000)
	})
	appender2 := strftime.AppendFunc(func(b []byte, t time.Time) []byte {
		return specificationHelper(b, t, "%02d", 10000000)
	})
	appender3 := strftime.AppendFunc(func(b []byte, t time.Time) []byte {
		return specificationHelper(b, t, "%03d", 1000000)
	})
	appender4 := strftime.AppendFunc(func(b []byte, t time.Time) []byte {
		return specificationHelper(b, t, "%04d", 100000)
	})
	appender5 := strftime.AppendFunc(func(b []byte, t time.Time) []byte {
		return specificationHelper(b, t, "%05d", 10000)
	})
	appender6 := strftime.AppendFunc(func(b []byte, t time.Time) []byte {
		return specificationHelper(b, t, "%06d", 1000)
	})
	appender7 := strftime.AppendFunc(func(b []byte, t time.Time) []byte {
		return specificationHelper(b, t, "%07d", 100)
	})
	appender8 := strftime.AppendFunc(func(b []byte, t time.Time) []byte {
		return specificationHelper(b, t, "%09d", 10)
	})
	appender9 := strftime.AppendFunc(func(b []byte, t time.Time) []byte {
		return specificationHelper(b, t, "%09d", 1)
	})

	ss := strftime.NewSpecificationSet()
	ss.Set('1', appender1)
	ss.Set('2', appender2)
	ss.Set('3', appender3)
	ss.Set('4', appender4)
	ss.Set('5', appender5)
	ss.Set('6', appender6)
	ss.Set('7', appender7)
	ss.Set('8', appender8)
	ss.Set('9', appender9)

	strftimeExtensions = strftime.WithSpecificationSet(ss)
}

// ================================================================
// Argument 1 is formatted date string like "2021-03-04 02:59:50".
// Argument 2 is format string like "%Y-%m-%d %H:%M:%S".
func MlrvalStrptime(input1, input2 *Mlrval) *Mlrval {
	if input1.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}
	if input2.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}
	timeString := input1.printrep
	formatString := input2.printrep

	t, err := strptime.Parse(timeString, formatString)
	if err != nil {
		return MLRVAL_ERROR
	}

	return MlrvalPointerFromFloat64(float64(t.UnixNano()) / 1.0e9)
}

// ================================================================
// Argument 1 is formatted date string like "2021-03-04T02:59:50Z".
func MlrvalGMT2Sec(input1 *Mlrval) *Mlrval {
	return MlrvalStrptime(input1, ptr_ISO8601_TIME_FORMAT)
}
