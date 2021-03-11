package types

import (
	"fmt"
	"math"
	"regexp"
	"strings"
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

// ================================================================
func MlrvalDHMS2Sec(input1 *Mlrval) *Mlrval {
	if input1.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}
	var d, h, m, s int

	if strings.HasPrefix(input1.printrep, "-") {

		n, err := fmt.Sscanf(input1.printrep, "-%dd%dh%dm%ds", &d, &h, &m, &s)
		if n == 4 && err == nil {
			return MlrvalPointerFromInt(-(s + m*60 + h*60*60 + d*60*60*24))
		}
		n, err = fmt.Sscanf(input1.printrep, "-%dh%dm%ds", &h, &m, &s)
		if n == 3 && err == nil {
			return MlrvalPointerFromInt(-(s + m*60 + h*60*60))
		}
		n, err = fmt.Sscanf(input1.printrep, "-%dm%ds", &m, &s)
		if n == 2 && err == nil {
			return MlrvalPointerFromInt(-(s + m*60))
		}
		n, err = fmt.Sscanf(input1.printrep, "-%ds", &s)
		if n == 1 && err == nil {
			return MlrvalPointerFromInt(-(s))
		}

	} else {

		n, err := fmt.Sscanf(input1.printrep, "%dd%dh%dm%ds", &d, &h, &m, &s)
		if n == 4 && err == nil {
			return MlrvalPointerFromInt(s + m*60 + h*60*60 + d*60*60*24)
		}
		n, err = fmt.Sscanf(input1.printrep, "%dh%dm%ds", &h, &m, &s)
		if n == 3 && err == nil {
			return MlrvalPointerFromInt(s + m*60 + h*60*60)
		}
		n, err = fmt.Sscanf(input1.printrep, "%dm%ds", &m, &s)
		if n == 2 && err == nil {
			return MlrvalPointerFromInt(s + m*60)
		}
		n, err = fmt.Sscanf(input1.printrep, "%ds", &s)
		if n == 1 && err == nil {
			return MlrvalPointerFromInt(s)
		}

	}
	return MLRVAL_ERROR
}

func MlrvalDHMS2FSec(input1 *Mlrval) *Mlrval {
	if input1.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}

	var d, h, m int
	var s float64

	if strings.HasPrefix(input1.printrep, "-") {

		n, err := fmt.Sscanf(input1.printrep, "-%dd%dh%dm%fs", &d, &h, &m, &s)
		if n == 4 && err == nil {
			return MlrvalPointerFromFloat64(-(s + float64(m*60+h*60*60+d*60*60*24)))
		}
		n, err = fmt.Sscanf(input1.printrep, "-%dh%dm%fs", &h, &m, &s)
		if n == 3 && err == nil {
			return MlrvalPointerFromFloat64(-(s + float64(m*60+h*60*60)))
		}
		n, err = fmt.Sscanf(input1.printrep, "-%dm%fs", &m, &s)
		if n == 2 && err == nil {
			return MlrvalPointerFromFloat64(-(s + float64(m*60)))
		}
		n, err = fmt.Sscanf(input1.printrep, "-%fs", &s)
		if n == 1 && err == nil {
			return MlrvalPointerFromFloat64(-(s))
		}

	} else {

		n, err := fmt.Sscanf(input1.printrep, "%dd%dh%dm%fs", &d, &h, &m, &s)
		if n == 4 && err == nil {
			return MlrvalPointerFromFloat64(s + float64(m*60+h*60*60+d*60*60*24))
		}
		n, err = fmt.Sscanf(input1.printrep, "%dh%dm%fs", &h, &m, &s)
		if n == 3 && err == nil {
			return MlrvalPointerFromFloat64(s + float64(m*60+h*60*60))
		}
		n, err = fmt.Sscanf(input1.printrep, "%dm%fs", &m, &s)
		if n == 2 && err == nil {
			return MlrvalPointerFromFloat64(s + float64(m*60))
		}
		n, err = fmt.Sscanf(input1.printrep, "%fs", &s)
		if n == 1 && err == nil {
			return MlrvalPointerFromFloat64(s)
		}

	}
	return MLRVAL_ERROR
}

func MlrvalHMS2Sec(input1 *Mlrval) *Mlrval {
	if input1.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}
	if input1.printrep == "" {
		return MLRVAL_ERROR
	}
	var h, m, s int

	if strings.HasPrefix(input1.printrep, "-") {
		n, err := fmt.Sscanf(input1.printrep, "-%d:%d:%d", &h, &m, &s)
		if n == 3 && err == nil {
			return MlrvalPointerFromInt(-(s + m*60 + h*60*60))
		}
	} else {
		n, err := fmt.Sscanf(input1.printrep, "%d:%d:%d", &h, &m, &s)
		if n == 3 && err == nil {
			return MlrvalPointerFromInt(s + m*60 + h*60*60)
		}
	}

	return MLRVAL_ERROR
}

func MlrvalHMS2FSec(input1 *Mlrval) *Mlrval {
	if input1.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}

	var h, m int
	var s float64

	if strings.HasPrefix(input1.printrep, "-") {
		n, err := fmt.Sscanf(input1.printrep, "-%d:%d:%f", &h, &m, &s)
		if n == 3 && err == nil {
			return MlrvalPointerFromFloat64(-(s + float64(m*60+h*60*60)))
		}
	} else {
		n, err := fmt.Sscanf(input1.printrep, "%d:%d:%f", &h, &m, &s)
		if n == 3 && err == nil {
			return MlrvalPointerFromFloat64(s + float64(m*60+h*60*60))
		}
	}

	return MLRVAL_ERROR
}

func MlrvalSec2DHMS(input1 *Mlrval) *Mlrval {
	isec, ok := input1.GetIntValue()
	if !ok {
		return MLRVAL_ERROR
	}

	var d, h, m, s int

	splitIntToDHMS(isec, &d, &h, &m, &s)
	if d != 0 {
		return MlrvalPointerFromString(
			fmt.Sprintf("%dd%02dh%02dm%02ds", d, h, m, s),
		)
	} else if h != 0 {
		return MlrvalPointerFromString(
			fmt.Sprintf("%dh%02dm%02ds", h, m, s),
		)
	} else if m != 0 {
		return MlrvalPointerFromString(
			fmt.Sprintf("%02dm%02ds", m, s),
		)
	} else {
		return MlrvalPointerFromString(
			fmt.Sprintf("%02ds", s),
		)
	}

	return MLRVAL_ERROR
}

func MlrvalSec2HMS(input1 *Mlrval) *Mlrval {
	isec, ok := input1.GetIntValue()
	if !ok {
		return MLRVAL_ERROR
	}
	sign := ""
	if isec < 0 {
		sign = "-"
		isec = -isec
	}

	var d, h, m, s int

	splitIntToDHMS(isec, &d, &h, &m, &s)
	h += d * 24

	return MlrvalPointerFromString(
		fmt.Sprintf("%s%02d:%02d:%02d", sign, h, m, s),
	)

	return MLRVAL_ERROR
}

func MlrvalFSec2DHMS(input1 *Mlrval) *Mlrval {
	fsec, ok := input1.GetNumericToFloatValue()
	if !ok {
		return MLRVAL_ERROR
	}

	sign := 1
	if fsec < 0 {
		sign = -1
		fsec = -fsec
	}
	isec := int(math.Trunc(fsec))
	fractional := fsec - float64(isec)

	var d, h, m, s int

	splitIntToDHMS(isec, &d, &h, &m, &s)

	if d != 0 {
		d = sign * d
		return MlrvalPointerFromString(
			fmt.Sprintf(
				"%dd%02dh%02dm%09.6fs",
				d, h, m, float64(s)+fractional),
		)
	} else if h != 0 {
		h = sign * h
		return MlrvalPointerFromString(
			fmt.Sprintf(
				"%dh%02dm%09.6fs",
				h, m, float64(s)+fractional),
		)
	} else if m != 0 {
		m = sign * m
		return MlrvalPointerFromString(
			fmt.Sprintf(
				"%dm%09.6fs",
				m, float64(s)+fractional),
		)
	} else {
		s = sign * s
		fractional = float64(sign) * fractional
		return MlrvalPointerFromString(
			fmt.Sprintf(
				"%.6fs",
				float64(s)+fractional),
		)
	}
}

func MlrvalFSec2HMS(input1 *Mlrval) *Mlrval {
	fsec, ok := input1.GetNumericToFloatValue()
	if !ok {
		return MLRVAL_ERROR
	}

	sign := ""
	if fsec < 0 {
		sign = "-"
		fsec = -fsec
	}
	isec := int(math.Trunc(fsec))
	fractional := fsec - float64(isec)

	var d, h, m, s int

	splitIntToDHMS(isec, &d, &h, &m, &s)
	h += d * 24

	// "%02.6f" does not exist so we have to do our own zero-pad
	if s < 10 {
		return MlrvalPointerFromString(
			fmt.Sprintf("%s%02d:%02d:0%.6f", sign, h, m, float64(s)+fractional),
		)
	} else {
		return MlrvalPointerFromString(
			fmt.Sprintf("%s%02d:%02d:%.6f", sign, h, m, float64(s)+fractional),
		)
	}

	return MLRVAL_ERROR
}

// Helper function
func splitIntToDHMS(u int, pd, ph, pm, ps *int) {
	d := 0
	h := 0
	m := 0
	s := 0
	sign := 1
	if u < 0 {
		u = -u
		sign = -1
	}
	s = u % 60
	u = u / 60
	if u == 0 {
		s = s * sign
	} else {
		m = u % 60
		u = u / 60
		if u == 0 {
			m = m * sign
		} else {
			h = u % 24
			u = u / 24
			if u == 0 {
				h = h * sign
			} else {
				d = u * sign
			}
		}
	}
	*pd = d
	*ph = h
	*pm = m
	*ps = s
}
