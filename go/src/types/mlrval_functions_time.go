package types

import (
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
	return MlrvalStrftime(input1, MlrvalPointerFromString("%Y-%m-%d"))
}

// ----------------------------------------------------------------
// Argument 1 is int/float seconds since the epoch.
// Argument 2 is format string like "%Y-%m-%d %H:%M:%S".
func MlrvalStrftime(input1, input2 *Mlrval) *Mlrval {
	epochSeconds, ok := input1.GetNumericToFloatValue()
	if !ok {
		return MLRVAL_ERROR
	}
	if input2.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}
	formatString := input2.printrep

	inputTime := lib.EpochSecondsToTime(epochSeconds)

	outputString, err := strftime.Format(formatString, inputTime)
	if err != nil {
		return MLRVAL_ERROR
	}

	return MlrvalPointerFromString(outputString)
}

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

// Argument 1 is formatted date string like "2021-03-04T02:59:50Z".
func MlrvalGMT2Sec(input1 *Mlrval) *Mlrval {
	return MlrvalStrptime(input1, ptr_ISO8601_TIME_FORMAT)
}
