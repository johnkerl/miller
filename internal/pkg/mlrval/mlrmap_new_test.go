package mlrval

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMlrmapAsRecord(t *testing.T) {
	mlrmap := newMlrmapUnhashed()
	assert.Equal(t, false, mlrmap.isHashed())
}

func TestNewMlrmap(t *testing.T) {
	mlrmap := NewMlrmap()
	assert.Equal(t, true, mlrmap.isHashed())
}
