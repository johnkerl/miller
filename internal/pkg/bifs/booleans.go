// ================================================================
// Boolean expressions for ==, !=, >, >=, <, <=
// ================================================================

package bifs

import (
	"github.com/johnkerl/miller/internal/pkg/mlrval"
)

func BIF_logical_NOT(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsBool() {
		return mlrval.FromBool(!input1.AcquireBoolValue())
	} else {
		return mlrval.FromTypeErrorUnary("!", input1)
	}
}

func BIF_logical_AND(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsBool() && input2.IsBool() {
		return mlrval.FromBool(input1.AcquireBoolValue() && input2.AcquireBoolValue())
	} else {
		return mlrval.FromTypeErrorUnary("&&", input1)
	}
}

func BIF_logical_OR(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsBool() && input2.IsBool() {
		return mlrval.FromBool(input1.AcquireBoolValue() || input2.AcquireBoolValue())
	} else {
		return mlrval.FromTypeErrorUnary("||", input1)
	}
}

func BIF_logical_XOR(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsBool() && input2.IsBool() {
		return mlrval.FromBool(input1.AcquireBoolValue() != input2.AcquireBoolValue())
	} else {
		return mlrval.FromTypeErrorUnary("^^", input1)
	}
}
