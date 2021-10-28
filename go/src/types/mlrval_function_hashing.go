package types

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

func BIF_md5(input1 *Mlrval) *Mlrval {
	if !input1.IsStringOrVoid() {
		return MLRVAL_ERROR
	} else {
		return MlrvalFromString(
			fmt.Sprintf(
				"%x",
				md5.Sum([]byte(input1.printrep)),
			),
		)
	}
}

func BIF_sha1(input1 *Mlrval) *Mlrval {
	if !input1.IsStringOrVoid() {
		return MLRVAL_ERROR
	} else {
		return MlrvalFromString(
			fmt.Sprintf(
				"%x",
				sha1.Sum([]byte(input1.printrep)),
			),
		)
	}
}

func BIF_sha256(input1 *Mlrval) *Mlrval {
	if !input1.IsStringOrVoid() {
		return MLRVAL_ERROR
	} else {
		return MlrvalFromString(
			fmt.Sprintf(
				"%x",
				sha256.Sum256([]byte(input1.printrep)),
			),
		)
	}
}

func BIF_sha512(input1 *Mlrval) *Mlrval {
	if !input1.IsStringOrVoid() {
		return MLRVAL_ERROR
	} else {
		return MlrvalFromString(
			fmt.Sprintf(
				"%x",
				sha512.Sum512([]byte(input1.printrep)),
			),
		)
	}
}
