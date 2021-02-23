// ================================================================
// TODO: comment how these are part of the copy-reduction project.

package types

var value_MLRVAL_ERROR = MlrvalFromError()
var value_MLRVAL_ABSENT = MlrvalFromAbsent()
var value_MLRVAL_VOID = MlrvalFromVoid()
var value_MLRVAL_INT_0 = MlrvalFromInt(0)
var value_MLRVAL_INT_1 = MlrvalFromInt(1)
var value_MLRVAL_FLOAT_0 = MlrvalFromFloat64(0)
var value_MLRVAL_FLOAT_1 = MlrvalFromFloat64(1)
var value_MLRVAL_TRUE = MlrvalFromTrue()
var value_MLRVAL_FALSE = MlrvalFromFalse()

var MLRVAL_ERROR = &value_MLRVAL_ERROR
var MLRVAL_ABSENT = &value_MLRVAL_ABSENT
var MLRVAL_VOID = &value_MLRVAL_VOID
var MLRVAL_INT_0 = &value_MLRVAL_INT_0
var MLRVAL_INT_1 = &value_MLRVAL_INT_1
var MLRVAL_FLOAT_0 = &value_MLRVAL_FLOAT_0
var MLRVAL_FLOAT_1 = &value_MLRVAL_FLOAT_1
var MLRVAL_TRUE = &value_MLRVAL_TRUE
var MLRVAL_FALSE = &value_MLRVAL_FALSE
