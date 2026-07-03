package bifs

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"

	"github.com/johnkerl/miller/v6/pkg/mlrval"
)

// hashableBytes returns the raw payload to hash: for strings, their UTF-8
// bytes; for bytes values, the bytes themselves.
func hashableBytes(funcname string, input1 *mlrval.Mlrval) ([]byte, *mlrval.Mlrval) {
	if input1.IsBytes() {
		return input1.AcquireBytesValue(), nil
	}
	if !input1.IsStringOrVoid() {
		return nil, mlrval.FromNotStringError(funcname, input1)
	}
	return []byte(input1.AcquireStringValue()), nil
}

func BIF_md5(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	payload, errval := hashableBytes("md5", input1)
	if errval != nil {
		return errval
	}
	return mlrval.FromString(fmt.Sprintf("%x", md5.Sum(payload)))
}

func BIF_sha1(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	payload, errval := hashableBytes("sha1", input1)
	if errval != nil {
		return errval
	}
	return mlrval.FromString(fmt.Sprintf("%x", sha1.Sum(payload)))
}

func BIF_sha256(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	payload, errval := hashableBytes("sha256", input1)
	if errval != nil {
		return errval
	}
	return mlrval.FromString(fmt.Sprintf("%x", sha256.Sum256(payload)))
}

func BIF_sha512(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	payload, errval := hashableBytes("sha512", input1)
	if errval != nil {
		return errval
	}
	return mlrval.FromString(fmt.Sprintf("%x", sha512.Sum512(payload)))
}
