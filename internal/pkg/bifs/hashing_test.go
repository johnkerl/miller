package bifs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnkerl/miller/internal/pkg/mlrval"
)

func TestBIF_md5(t *testing.T) {
	input1 := mlrval.FromDeferredType("")
	output := BIF_md5(input1)
	stringval, ok := output.GetStringValue()
	assert.True(t, ok)
	assert.Equal(t, "d41d8cd98f00b204e9800998ecf8427e", stringval)

	input1 = mlrval.FromDeferredType("miller")
	output = BIF_md5(input1)
	stringval, ok = output.GetStringValue()
	assert.True(t, ok)
	assert.Equal(t, "f0af962ddbc82430e947390b2f3f6e49", stringval)
}

// TODO: copy in more unit-test cases from existing regression-test data
