package bifs

import (
	"encoding/base64"

	"github.com/johnkerl/miller/v6/pkg/mlrval"
)

func BIF_base64_encode(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if input1.IsBytes() {
		return mlrval.FromString(
			base64.StdEncoding.EncodeToString(input1.AcquireBytesValue()),
		)
	}
	if !input1.IsStringOrVoid() {
		return mlrval.FromNotStringError("base64_encode", input1)
	}
	return mlrval.FromString(
		base64.StdEncoding.EncodeToString([]byte(input1.AcquireStringValue())),
	)
}

func BIF_base64_decode(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsStringOrVoid() {
		return mlrval.FromNotStringError("base64_decode", input1)
	}
	decoded, err := base64.StdEncoding.DecodeString(input1.AcquireStringValue())
	if err != nil {
		return mlrval.FromError(err)
	}
	return mlrval.FromBytes(decoded)
}
