package types

import (
	"fmt"
	"regexp"
	"time"

	"github.com/lestrrat-go/strftime"
	strptime "mlr/internal/pkg/pbnjay-strptime"

	"mlr/internal/pkg/lib"
)

const ISO8601_TIME_FORMAT = "%Y-%m-%dT%H:%M:%SZ"

var ptr_ISO8601_TIME_FORMAT = MlrvalFromString("%Y-%m-%dT%H:%M:%SZ")
var ptr_ISO8601_LOCAL_TIME_FORMAT = MlrvalFromString("%Y-%m-%d %H:%M:%S")
var ptr_YMD_FORMAT = MlrvalFromString("%Y-%m-%d")

// ================================================================
func BIF_systime() *Mlrval {
	return MlrvalFromFloat64(
		float64(time.Now().UnixNano()) / 1.0e9,
	)
}
func BIF_systimeint() *Mlrval {
	return MlrvalFromInt(int(time.Now().Unix()))
}

var startTime float64

func init() {
	startTime = float64(time.Now().UnixNano()) / 1.0e9
}
func BIF_uptime() *Mlrval {
	return MlrvalFromFloat64(
		float64(time.Now().UnixNano())/1.0e9 - startTime,
	)
}

// ================================================================

func BIF_sec2gmt_unary(input1 *Mlrval) *Mlrval {
	floatValue, isNumeric := input1.GetNumericToFloatValue()
	if !isNumeric {
		return input1
	}

	numDecimalPlaces := 0

	return MlrvalFromString(lib.Sec2GMT(floatValue, numDecimalPlaces))
}

func BIF_sec2gmt_binary(input1, input2 *Mlrval) *Mlrval {
	floatValue, isNumeric := input1.GetNumericToFloatValue()
	if !isNumeric {
		return input1
	}

	numDecimalPlaces, isInt := input2.GetIntValue()
	if !isInt {
		return MLRVAL_ERROR
	}

	return MlrvalFromString(lib.Sec2GMT(floatValue, numDecimalPlaces))
}

func BIF_sec2localtime_unary(input1 *Mlrval) *Mlrval {
	floatValue, isNumeric := input1.GetNumericToFloatValue()
	if !isNumeric {
		return input1
	}

	numDecimalPlaces := 0

	return MlrvalFromString(lib.Sec2LocalTime(floatValue, numDecimalPlaces))
}

func BIF_sec2localtime_binary(input1, input2 *Mlrval) *Mlrval {
	floatValue, isNumeric := input1.GetNumericToFloatValue()
	if !isNumeric {
		return input1
	}

	numDecimalPlaces, isInt := input2.GetIntValue()
	if !isInt {
		return MLRVAL_ERROR
	}

	return MlrvalFromString(lib.Sec2LocalTime(floatValue, numDecimalPlaces))
}

func BIF_sec2localtime_ternary(input1, input2, input3 *Mlrval) *Mlrval {
	floatValue, isNumeric := input1.GetNumericToFloatValue()
	if !isNumeric {
		return input1
	}

	numDecimalPlaces, isInt := input2.GetIntValue()
	if !isInt {
		return MLRVAL_ERROR
	}

	locationString, isString := input3.GetString()
	if !isString {
		return MLRVAL_ERROR
	}

	location, err := time.LoadLocation(locationString)
	if err != nil {
		return MLRVAL_ERROR
	}

	return MlrvalFromString(lib.Sec2LocationTime(floatValue, numDecimalPlaces, location))
}

func BIF_sec2gmtdate(input1 *Mlrval) *Mlrval {
	if !input1.IsNumeric() {
		return input1
	}
	return BIF_strftime(input1, ptr_YMD_FORMAT)
}

func BIF_sec2localdate_unary(input1 *Mlrval) *Mlrval {
	if !input1.IsNumeric() {
		return input1
	}
	return BIF_strftime_local_binary(input1, ptr_YMD_FORMAT)
}

func BIF_sec2localdate_binary(input1, input2 *Mlrval) *Mlrval {
	if !input1.IsNumeric() {
		return input1
	}
	return BIF_strftime_local_ternary(input1, ptr_YMD_FORMAT, input2)
}

// ----------------------------------------------------------------

func BIF_localtime2gmt_unary(input1 *Mlrval) *Mlrval {
	if !input1.IsString() {
		return MLRVAL_ERROR
	}
	return BIF_sec2gmt_unary(BIF_localtime2sec_unary(input1))
}

func BIF_localtime2gmt_binary(input1, input2 *Mlrval) *Mlrval {
	if !input1.IsString() {
		return MLRVAL_ERROR
	}
	return BIF_sec2gmt_unary(BIF_localtime2sec_binary(input1, input2))
}

func BIF_gmt2localtime_unary(input1 *Mlrval) *Mlrval {
	if !input1.IsString() {
		return MLRVAL_ERROR
	}
	return BIF_sec2localtime_unary(BIF_gmt2sec(input1))
}

func BIF_gmt2localtime_binary(input1, input2 *Mlrval) *Mlrval {
	if !input1.IsString() {
		return MLRVAL_ERROR
	}
	return BIF_sec2localtime_ternary(BIF_gmt2sec(input1), MlrvalFromInt(0), input2)
}

// ================================================================
// Argument 1 is int/float seconds since the epoch.
// Argument 2 is format string like "%Y-%m-%d %H:%M:%S".

var extensionRegex = regexp.MustCompile("([1-9])S")

func BIF_strftime(input1, input2 *Mlrval) *Mlrval {
	return strftimeHelper(input1, input2, false, nil)
}

func BIF_strftime_local_binary(input1, input2 *Mlrval) *Mlrval {
	return strftimeHelper(input1, input2, true, nil)
}

func BIF_strftime_local_ternary(input1, input2, input3 *Mlrval) *Mlrval {
	locationString, isString := input3.GetString()
	if !isString {
		return MLRVAL_ERROR
	}

	location, err := time.LoadLocation(locationString)
	if err != nil {
		return MLRVAL_ERROR
	}

	return strftimeHelper(input1, input2, true, location)
}

func strftimeHelper(input1, input2 *Mlrval, doLocal bool, location *time.Location) *Mlrval {
	if input1.mvtype == MT_VOID {
		return input1
	}
	epochSeconds, ok := input1.GetNumericToFloatValue()
	if !ok {
		return MLRVAL_ERROR
	}
	if input2.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}

	// Convert argument1 from float seconds since the epoch to a Go time.
	var inputTime time.Time
	if doLocal {
		if location != nil {
			inputTime = lib.EpochSecondsToLocationTime(epochSeconds, location)
		} else {
			inputTime = lib.EpochSecondsToLocalTime(epochSeconds)
		}
	} else {
		inputTime = lib.EpochSecondsToGMT(epochSeconds)
	}

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

	return MlrvalFromString(outputString)
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
func BIF_strptime(input1, input2 *Mlrval) *Mlrval {
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

	return MlrvalFromFloat64(float64(t.UnixNano()) / 1.0e9)
}

func BIF_strptime_local_binary(input1, input2 *Mlrval) *Mlrval {
	if input1.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}
	if input2.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}
	timeString := input1.printrep
	formatString := input2.printrep

	t, err := strptime.ParseLocal(timeString, formatString)
	if err != nil {
		return MLRVAL_ERROR
	}

	return MlrvalFromFloat64(float64(t.UnixNano()) / 1.0e9)
}

func BIF_strptime_local_ternary(input1, input2, input3 *Mlrval) *Mlrval {
	if input1.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}
	if input2.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}
	if input3.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}

	timeString := input1.printrep
	formatString := input2.printrep
	locationString := input3.printrep

	location, err := time.LoadLocation(locationString)
	if err != nil {
		return MLRVAL_ERROR
	}

	// TODO: use location

	t, err := strptime.ParseLocation(timeString, formatString, location)
	if err != nil {
		return MLRVAL_ERROR
	}

	return MlrvalFromFloat64(float64(t.UnixNano()) / 1.0e9)
}

// ================================================================
// Argument 1 is formatted date string like "2021-03-04T02:59:50Z".
func BIF_gmt2sec(input1 *Mlrval) *Mlrval {
	return BIF_strptime(input1, ptr_ISO8601_TIME_FORMAT)
}

func BIF_localtime2sec_unary(input1 *Mlrval) *Mlrval {
	return BIF_strptime_local_binary(input1, ptr_ISO8601_LOCAL_TIME_FORMAT)
}

func BIF_localtime2sec_binary(input1, input2 *Mlrval) *Mlrval {
	return BIF_strptime_local_ternary(input1, ptr_ISO8601_LOCAL_TIME_FORMAT, input2)
}
