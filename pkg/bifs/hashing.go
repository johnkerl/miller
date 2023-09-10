package bifs

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"

	"github.com/johnkerl/miller/pkg/mlrval"
)

func BIF_md5(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsStringOrVoid() {
		return mlrval.FromNotStringError("md5", input1)
	} else {
		return mlrval.FromString(
			fmt.Sprintf(
				"%x",
				md5.Sum([]byte(input1.AcquireStringValue())),
			),
		)
	}
}

func BIF_sha1(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsStringOrVoid() {
		return mlrval.FromNotStringError("sha1", input1)
	} else {
		return mlrval.FromString(
			fmt.Sprintf(
				"%x",
				sha1.Sum([]byte(input1.AcquireStringValue())),
			),
		)
	}
}

func BIF_sha256(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsStringOrVoid() {
		return mlrval.FromNotStringError("sha256", input1)
	} else {
		return mlrval.FromString(
			fmt.Sprintf(
				"%x",
				sha256.Sum256([]byte(input1.AcquireStringValue())),
			),
		)
	}
}

func BIF_sha512(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsStringOrVoid() {
		return mlrval.FromNotStringError("sha512", input1)
	} else {
		return mlrval.FromString(
			fmt.Sprintf(
				"%x",
				sha512.Sum512([]byte(input1.AcquireStringValue())),
			),
		)
	}
}
