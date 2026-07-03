package bifs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnkerl/miller/v6/pkg/mlrval"
)

func TestBIF_base64_encode(t *testing.T) {
	input1 := mlrval.FromDeferredType("")
	output := BIF_base64_encode(input1)
	stringval, ok := output.GetStringValue()
	assert.True(t, ok)
	assert.Equal(t, "", stringval)

	input1 = mlrval.FromDeferredType("hello")
	output = BIF_base64_encode(input1)
	stringval, ok = output.GetStringValue()
	assert.True(t, ok)
	assert.Equal(t, "aGVsbG8=", stringval)

	// Bytes input encodes the raw bytes
	output = BIF_base64_encode(mlrval.FromBytes([]byte{0xff}))
	stringval, ok = output.GetStringValue()
	assert.True(t, ok)
	assert.Equal(t, "/w==", stringval)
}

func TestBIF_base64_decode(t *testing.T) {
	// Decode always yields bytes
	output := BIF_base64_decode(mlrval.FromDeferredType("aGVsbG8="))
	assert.True(t, output.IsBytes())
	assert.Equal(t, []byte("hello"), output.AcquireBytesValue())

	// string() recovers the text for UTF-8 payloads
	assert.Equal(t, "hello", BIF_string(output).String())

	// Non-UTF-8 payloads round-trip losslessly
	output = BIF_base64_decode(mlrval.FromDeferredType("/w=="))
	assert.True(t, output.IsBytes())
	assert.Equal(t, []byte{0xff}, output.AcquireBytesValue())
	assert.Equal(t, "ff", output.String())

	// Malformed base64 input results in an error mlrval
	output = BIF_base64_decode(mlrval.FromDeferredType("not!valid!base64"))
	assert.True(t, output.IsError())
}

func TestBIF_base64_round_trip(t *testing.T) {
	input1 := mlrval.FromBytes([]byte{0x00, 0x01, 0xfe, 0xff})
	output := BIF_base64_decode(BIF_base64_encode(input1))
	assert.True(t, output.IsBytes())
	assert.Equal(t, []byte{0x00, 0x01, 0xfe, 0xff}, output.AcquireBytesValue())
}

func TestBIF_hex_encode_decode(t *testing.T) {
	output := BIF_hex_encode(mlrval.FromBytes([]byte{0xde, 0xad, 0xbe, 0xef}))
	stringval, ok := output.GetStringValue()
	assert.True(t, ok)
	assert.Equal(t, "deadbeef", stringval)

	output = BIF_hex_encode(mlrval.FromDeferredType("hi"))
	stringval, ok = output.GetStringValue()
	assert.True(t, ok)
	assert.Equal(t, "6869", stringval)

	output = BIF_hex_decode(mlrval.FromDeferredType("deadbeef"))
	assert.True(t, output.IsBytes())
	assert.Equal(t, []byte{0xde, 0xad, 0xbe, 0xef}, output.AcquireBytesValue())

	// Odd length and non-hex characters are errors
	assert.True(t, BIF_hex_decode(mlrval.FromDeferredType("abc")).IsError())
	assert.True(t, BIF_hex_decode(mlrval.FromDeferredType("zz")).IsError())
}

func TestBIF_bytes_conversions(t *testing.T) {
	// bytes(string) → bytes; string(bytes) is the inverse
	b := BIF_bytes(mlrval.FromString("héllo"))
	assert.True(t, b.IsBytes())
	assert.Equal(t, []byte("héllo"), b.AcquireBytesValue())
	assert.Equal(t, "héllo", BIF_string(b).String())

	// bytes(bytes) is the identity
	assert.Equal(t, b, BIF_bytes(b))

	// bytes of a number is a type error
	assert.True(t, BIF_bytes(mlrval.FromInt(5)).IsError())
}

func TestBIF_bytes_strlen_substr(t *testing.T) {
	// héllo is 5 runes, 6 bytes
	s := mlrval.FromString("héllo")
	b := BIF_bytes(s)
	assert.Equal(t, "5", BIF_strlen(s).String())
	assert.Equal(t, "6", BIF_strlen(b).String())

	// Byte-oriented slicing
	payload := mlrval.FromBytes([]byte{0x00, 0x01, 0x02, 0x03})
	sliced := BIF_substr_0_up(payload, mlrval.FromInt(1), mlrval.FromInt(2))
	assert.True(t, sliced.IsBytes())
	assert.Equal(t, []byte{0x01, 0x02}, sliced.AcquireBytesValue())

	sliced = BIF_substr_1_up(payload, mlrval.FromInt(1), mlrval.FromInt(2))
	assert.True(t, sliced.IsBytes())
	assert.Equal(t, []byte{0x00, 0x01}, sliced.AcquireBytesValue())

	// Empty slice preserves the bytes type
	sliced = BIF_substr_0_up(payload, mlrval.FromInt(3), mlrval.FromInt(1))
	assert.True(t, sliced.IsBytes())
	assert.Equal(t, 0, len(sliced.AcquireBytesValue()))
}

func TestBIF_bytes_hashing(t *testing.T) {
	// Hash of bytes equals hash of the same payload as a string
	fromString := BIF_md5(mlrval.FromDeferredType("miller"))
	fromBytes := BIF_md5(BIF_bytes(mlrval.FromString("miller")))
	assert.Equal(t, fromString.String(), fromBytes.String())
	assert.Equal(t, "f0af962ddbc82430e947390b2f3f6e49", fromBytes.String())

	// Non-UTF-8 bytes hash without corruption
	output := BIF_sha256(mlrval.FromBytes([]byte{0xff}))
	stringval, ok := output.GetStringValue()
	assert.True(t, ok)
	assert.Equal(t, "a8100ae6aa1940d0b663bb31cd466142ebbdbd5187131b92d93818987832eb89", stringval)
}
