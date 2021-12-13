// ================================================================
// Tests mlrval-to-string logic
// ================================================================

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFoo(t *testing.T) {
	assert.Equal(t, false, true)
}

//func TestGetString(t *testing.T) {
//	mv := FromInferredType("234")
//	stringval, ok := mv.GetString()
//	assert.False(t, ok)
//
//	mv = FromDeferredType("234")
//	stringval, ok = mv.GetString()
//	assert.False(t, ok)
//
//	mv = FromInferredType("234.5")
//	stringval, ok = mv.GetString()
//	assert.False(t, ok)
//
//	mv = FromDeferredType("234.5")
//	stringval, ok = mv.GetString()
//	assert.False(t, ok)
//
//	mv = FromInferredType("abc")
//	stringval, ok = mv.GetString()
//	assert.Equal(t, "abc", stringval)
//	assert.True(t, ok)
//
//	mv = FromDeferredType("abc")
//	stringval, ok = mv.GetString()
//	assert.Equal(t, "abc", stringval)
//	assert.True(t, ok)
//
//	mv = FromInferredType("")
//	stringval, ok = mv.GetString()
//	assert.Equal(t, "", stringval)
//	assert.True(t, ok)
//
//	mv = FromDeferredType("")
//	stringval, ok = mv.GetString()
//	assert.Equal(t, "", stringval)
//	assert.True(t, ok)
//}
