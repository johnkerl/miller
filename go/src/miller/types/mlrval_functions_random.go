package types

import (
	"miller/lib"
)

// ================================================================
func MlrvalUrand() Mlrval {
	return MlrvalFromFloat64(
		lib.RandFloat64(),
	)
}

func MlrvalUrand32() Mlrval {
	return MlrvalFromInt64(
		int64(
			lib.RandUint32(),
		),
	)
}
