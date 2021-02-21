package types

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

func MlrvalMD5(input1 *Mlrval) Mlrval {
	if !input1.IsStringOrVoid() {
		return MlrvalFromError()
	}
	return MlrvalFromString(
		fmt.Sprintf(
			"%x",
			md5.Sum([]byte(input1.printrep)),
		),
	)
}

func MlrvalSHA1(input1 *Mlrval) Mlrval {
	if !input1.IsStringOrVoid() {
		return MlrvalFromError()
	}
	return MlrvalFromString(
		fmt.Sprintf(
			"%x",
			sha1.Sum([]byte(input1.printrep)),
		),
	)
}

func MlrvalSHA256(input1 *Mlrval) Mlrval {
	if !input1.IsStringOrVoid() {
		return MlrvalFromError()
	}
	return MlrvalFromString(
		fmt.Sprintf(
			"%x",
			sha256.Sum256([]byte(input1.printrep)),
		),
	)
}

func MlrvalSHA512(input1 *Mlrval) Mlrval {
	if !input1.IsStringOrVoid() {
		return MlrvalFromError()
	}
	return MlrvalFromString(
		fmt.Sprintf(
			"%x",
			sha512.Sum512([]byte(input1.printrep)),
		),
	)
}
