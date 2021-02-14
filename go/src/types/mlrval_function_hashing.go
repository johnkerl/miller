package types

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

func MlrvalMD5(ma *Mlrval) Mlrval {
	if !ma.IsStringOrVoid() {
		return MlrvalFromError()
	}
	return MlrvalFromString(
		fmt.Sprintf(
			"%x",
			md5.Sum([]byte(ma.printrep)),
		),
	)
}

func MlrvalSHA1(ma *Mlrval) Mlrval {
	if !ma.IsStringOrVoid() {
		return MlrvalFromError()
	}
	return MlrvalFromString(
		fmt.Sprintf(
			"%x",
			sha1.Sum([]byte(ma.printrep)),
		),
	)
}

func MlrvalSHA256(ma *Mlrval) Mlrval {
	if !ma.IsStringOrVoid() {
		return MlrvalFromError()
	}
	return MlrvalFromString(
		fmt.Sprintf(
			"%x",
			sha256.Sum256([]byte(ma.printrep)),
		),
	)
}

func MlrvalSHA512(ma *Mlrval) Mlrval {
	if !ma.IsStringOrVoid() {
		return MlrvalFromError()
	}
	return MlrvalFromString(
		fmt.Sprintf(
			"%x",
			sha512.Sum512([]byte(ma.printrep)),
		),
	)
}
