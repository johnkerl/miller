package types

import (
	"time"

	"miller/src/lib"
)

// ================================================================
func MlrvalSystime() Mlrval {
	return MlrvalFromFloat64(
		float64(time.Now().UnixNano()) / 1.0e9,
	)
}
func MlrvalSystimeInt() Mlrval {
	return MlrvalFromInt(int(time.Now().Unix()))
}

var startTime float64

func init() {
	startTime = float64(time.Now().UnixNano()) / 1.0e9
}
func MlrvalUptime() Mlrval {
	return MlrvalFromFloat64(
		float64(time.Now().UnixNano())/1.0e9 - startTime,
	)
}

// ================================================================
func MlrvalSec2GMTUnary(ma *Mlrval) Mlrval {
	if ma.mvtype == MT_FLOAT {
		return MlrvalFromString(lib.Sec2GMT(ma.floatval, 0))
	} else if ma.mvtype == MT_INT {
		return MlrvalFromString(lib.Sec2GMT(float64(ma.intval), 0))
	} else {
		return *ma
	}
}

// ----------------------------------------------------------------
func MlrvalSec2GMTBinary(ma, mb *Mlrval) Mlrval {
	if mb.mvtype != MT_INT {
		return MlrvalFromError()
	}
	if ma.mvtype == MT_FLOAT {
		return MlrvalFromString(lib.Sec2GMT(ma.floatval, int(mb.intval)))
	} else if ma.mvtype == MT_INT {
		return MlrvalFromString(lib.Sec2GMT(float64(ma.intval), int(mb.intval)))
	} else {
		return *ma
	}
}
