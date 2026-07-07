package bifs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnkerl/miller/v6/pkg/mlrval"
)

func TestBIF_sparkline(t *testing.T) {
	// Needs array or map
	input := mlrval.FromInt(3)
	output := BIF_sparkline(input)
	assert.True(t, output.IsError())

	// Non-numeric element
	input = mlrval.FromArray([]*mlrval.Mlrval{
		mlrval.FromInt(1),
		mlrval.FromString("abc"),
	})
	output = BIF_sparkline(input)
	assert.True(t, output.IsError())

	// Empty array is void
	input = mlrval.FromArray([]*mlrval.Mlrval{})
	output = BIF_sparkline(input)
	assert.True(t, output.IsVoid())

	// Ascending values span the full tick range, low to high
	input = mlrval.FromArray([]*mlrval.Mlrval{
		mlrval.FromInt(1),
		mlrval.FromInt(2),
		mlrval.FromInt(3),
		mlrval.FromInt(4),
		mlrval.FromInt(5),
		mlrval.FromInt(6),
		mlrval.FromInt(7),
		mlrval.FromInt(8),
	})
	output = BIF_sparkline(input)
	assert.True(t, mlrval.Equals(output, mlrval.FromString("▁▂▃▄▅▆▇█")))

	// Same value throughout maps to the lowest tick
	input = mlrval.FromArray([]*mlrval.Mlrval{
		mlrval.FromInt(3),
		mlrval.FromInt(3),
		mlrval.FromInt(3),
	})
	output = BIF_sparkline(input)
	assert.True(t, mlrval.Equals(output, mlrval.FromString("▁▁▁")))

	// Map input follows insertion order
	input = array_to_map_for_test(mlrval.FromArray([]*mlrval.Mlrval{
		mlrval.FromInt(1),
		mlrval.FromInt(8),
	}))
	output = BIF_sparkline(input)
	assert.True(t, mlrval.Equals(output, mlrval.FromString("▁█")))
}
