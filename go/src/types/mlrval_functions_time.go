package types

import (
	"time"

	"miller/src/lib"
)

// ================================================================
func MlrvalSystime(output *Mlrval) {
	output.SetFromFloat64(
		float64(time.Now().UnixNano()) / 1.0e9,
	)
}
func MlrvalSystimeInt(output *Mlrval) {
	output.SetFromInt(int(time.Now().Unix()))
}

var startTime float64

func init() {
	startTime = float64(time.Now().UnixNano()) / 1.0e9
}
func MlrvalUptime(output *Mlrval) {
	output.SetFromFloat64(
		float64(time.Now().UnixNano())/1.0e9 - startTime,
	)
}

// ================================================================
func MlrvalSec2GMTUnary(output, input1 *Mlrval) {
	if input1.mvtype == MT_FLOAT {
		output.SetFromString(lib.Sec2GMT(input1.floatval, 0))
	} else if input1.mvtype == MT_INT {
		output.SetFromString(lib.Sec2GMT(float64(input1.intval), 0))
	} else {
		output.CopyFrom(input1)
	}
}

// ----------------------------------------------------------------
func MlrvalSec2GMTBinary(input1, input2 *Mlrval) Mlrval {
	// xxx temp
	output := MlrvalFromAbsent()
	if input2.mvtype != MT_INT {
		output.SetFromError()
	} else if input1.mvtype == MT_FLOAT {
		output.SetFromString(lib.Sec2GMT(input1.floatval, int(input2.intval)))
	} else if input1.mvtype == MT_INT {
		output.SetFromString(lib.Sec2GMT(float64(input1.intval), int(input2.intval)))
	} else {
		output.CopyFrom(input1)
	}
	return output
}
