package types

import (
	"time"

	"github.com/lestrrat-go/strftime"
	"github.com/pbnjay/strptime"

	"miller/src/lib"
)

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
