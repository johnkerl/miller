package types

import (
	"time"

	"miller/lib"
)

// ================================================================
func MlrvalSystime() Mlrval {
	return MlrvalFromFloat64(
		float64(time.Now().UnixNano()) / 1.0e9,
	)
}
func MlrvalSystimeInt() Mlrval {
	return MlrvalFromInt64(time.Now().Unix())
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
