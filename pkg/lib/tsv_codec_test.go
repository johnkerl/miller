package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTSVDecodeField(t *testing.T) {
	assert.Equal(t, "", TSVDecodeField(""))
	assert.Equal(t, "a", TSVDecodeField("a"))
	assert.Equal(t, "abc", TSVDecodeField("abc"))
	assert.Equal(t, `\`, TSVDecodeField(`\`))
	assert.Equal(t, "\n", TSVDecodeField(`\n`))
	assert.Equal(t, "\r", TSVDecodeField(`\r`))
	assert.Equal(t, "\t", TSVDecodeField(`\t`))
	assert.Equal(t, "\\", TSVDecodeField(`\\`))
	assert.Equal(t, `\n`, TSVDecodeField(`\\n`))
	assert.Equal(t, "\\\n", TSVDecodeField(`\\\n`))
	assert.Equal(t, "abc\r\ndef\r\n", TSVDecodeField(`abc\r\ndef\r\n`))
}

func TestTSVEncodeField(t *testing.T) {
	assert.Equal(t, "", TSVEncodeField(""))
	assert.Equal(t, "a", TSVEncodeField("a"))
	assert.Equal(t, "abc", TSVEncodeField("abc"))
	assert.Equal(t, `\\`, TSVEncodeField(`\`))
	assert.Equal(t, `\n`, TSVEncodeField("\n"))
	assert.Equal(t, `\r`, TSVEncodeField("\r"))
	assert.Equal(t, `\t`, TSVEncodeField("\t"))
	assert.Equal(t, `\\`, TSVEncodeField("\\"))
	assert.Equal(t, `\\n`, TSVEncodeField("\\n"))
	assert.Equal(t, `\\\n`, TSVEncodeField("\\\n"))
	assert.Equal(t, `abc\r\ndef\r\n`, TSVEncodeField("abc\r\ndef\r\n"))
}
