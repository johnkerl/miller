package bifs

import (
	"encoding/hex"

	"github.com/johnkerl/miller/v6/pkg/mlrval"
)

func BIF_hex_encode(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsBytes() {
		return mlrval.FromString(
			hex.EncodeToString(input1.AcquireBytesValue()),
		)
	}
	if !input1.IsStringOrVoid() {
		return mlrval.FromNotStringError("hex_encode", input1)
	}
	return mlrval.FromString(
		hex.EncodeToString([]byte(input1.AcquireStringValue())),
	)
}

func BIF_hex_decode(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsStringOrVoid() {
		return mlrval.FromNotStringError("hex_decode", input1)
	}
	decoded, err := hex.DecodeString(input1.AcquireStringValue())
	if err != nil {
		return mlrval.FromError(err)
	}
	return mlrval.FromBytes(decoded)
}
