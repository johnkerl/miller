package types

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

func MlrvalMD5(output, input1 *Mlrval) {
	if !input1.IsStringOrVoid() {
		output.SetFromError()
	} else {
		output.SetFromString(
			fmt.Sprintf(
				"%x",
				md5.Sum([]byte(input1.printrep)),
			),
		)
	}
}

func MlrvalSHA1(output, input1 *Mlrval) {
	if !input1.IsStringOrVoid() {
		output.SetFromError()
	} else {
		output.SetFromString(
			fmt.Sprintf(
				"%x",
				sha1.Sum([]byte(input1.printrep)),
			),
		)
	}
}

func MlrvalSHA256(output, input1 *Mlrval) {
	if !input1.IsStringOrVoid() {
		output.SetFromError()
	} else {
		output.SetFromString(
			fmt.Sprintf(
				"%x",
				sha256.Sum256([]byte(input1.printrep)),
			),
		)
	}
}

func MlrvalSHA512(output, input1 *Mlrval) {
	if !input1.IsStringOrVoid() {
		output.SetFromError()
	} else {
		output.SetFromString(
			fmt.Sprintf(
				"%x",
				sha512.Sum512([]byte(input1.printrep)),
			),
		)
	}
}
