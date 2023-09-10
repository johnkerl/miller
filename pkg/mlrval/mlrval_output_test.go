package mlrval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert.Equal(t, "234", FromInferredType("234").String())
	assert.Equal(t, "234", FromDeferredType("234").String())
	assert.Equal(t, "234.5", FromInferredType("234.5").String())
	assert.Equal(t, "234.5", FromDeferredType("234.5").String())
	assert.Equal(t, "abc", FromInferredType("abc").String())
	assert.Equal(t, "abc", FromDeferredType("abc").String())
	assert.Equal(t, "", FromInferredType("").String())
	assert.Equal(t, "", FromDeferredType("").String())
}
