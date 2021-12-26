package scan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDecimalDigit(t *testing.T) {
	var c byte
	for c = 0x00; c < 0xff; c++ {
		if c >= '0' && c <= '9' {
			assert.True(t, isDecimalDigit(c))
		} else {
			assert.False(t, isDecimalDigit(c))
		}
	}
}

func TestIsHexDigit(t *testing.T) {
	var c byte
	for c = 0x00; c < 0xff; c++ {
		if c >= '0' && c <= '9' {
			assert.True(t, isHexDigit(c))
		} else if c >= 'a' && c <= 'f' {
			assert.True(t, isHexDigit(c))
		} else if c >= 'A' && c <= 'F' {
			assert.True(t, isHexDigit(c))
		} else {
			assert.False(t, isHexDigit(c))
		}
	}
}

func TestIsFloatDigit(t *testing.T) {
	var c byte
	for c = 0x00; c < 0xff; c++ {
		if c >= '0' && c <= '9' {
			assert.True(t, isFloatDigit(c))
		} else if c == '.' || c == '-' || c == '+' || c == 'e' || c == 'E' {
			assert.True(t, isFloatDigit(c))
		} else {
			assert.False(t, isFloatDigit(c))
		}
	}
}
