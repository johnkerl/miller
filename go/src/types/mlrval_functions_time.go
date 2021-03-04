package types

import (
	"time"

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
func MlrvalStrftime(input1, input2 *Mlrval) *Mlrval {
	return MlrvalPointerFromString("foo")
}

func MlrvalStrptime(input1, input2 *Mlrval) *Mlrval {
	return MlrvalPointerFromString("foo")
}
