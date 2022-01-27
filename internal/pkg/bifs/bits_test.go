package bifs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnkerl/miller/internal/pkg/mlrval"
)

func TestBIF_bitcount(t *testing.T) {
	input1 := mlrval.FromDeferredType("0xcafe")
	output := BIF_bitcount(input1)
	intval, ok := output.GetIntValue()
	assert.True(t, ok)
	assert.Equal(t, int64(11), intval)
}

// TODO: copy in more unit-test cases from existing regression-test data
