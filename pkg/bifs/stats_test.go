package bifs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnkerl/miller/pkg/mlrval"
)

func stats_test_array(n int) *mlrval.Mlrval {
	a := make([]*mlrval.Mlrval, n)
	for i := 0; i < n; i++ {
		a[i] = mlrval.FromInt(int64(i))
	}
	return mlrval.FromArray(a)
}

func array_to_map_for_test(a *mlrval.Mlrval) *mlrval.Mlrval {
	array := a.AcquireArrayValue()
	m := mlrval.NewMlrmap()
	for i := 0; i < len(array); i++ {
		key := fmt.Sprint(i)
		val := array[i]
		m.PutCopy(key, val)
	}
	return mlrval.FromMap(m)
}

func TestBIF_count(t *testing.T) {
	// Needs array or map
	input := mlrval.FromInt(3)
	output := BIF_count(input)
	assert.True(t, output.IsError())

	for n := 0; n < 5; n++ {
		input = stats_test_array(n)
		assert.True(t, mlrval.Equals(BIF_count(input), mlrval.FromInt(int64(n))))

		input = array_to_map_for_test(input)
		assert.True(t, mlrval.Equals(BIF_count(input), mlrval.FromInt(int64(n))))
	}
}

func TestBIF_distinct_count(t *testing.T) {
	// Needs array or map
	input := mlrval.FromInt(3)
	output := BIF_count(input)
	assert.True(t, output.IsError())

	input = mlrval.FromArray([]*mlrval.Mlrval{
		mlrval.FromInt(1),
		mlrval.FromInt(2),
		mlrval.FromInt(3),
		mlrval.FromInt(1),
		mlrval.FromInt(2),
	})
	assert.True(t, mlrval.Equals(BIF_distinct_count(input), mlrval.FromInt(3)))

	input = array_to_map_for_test(input)
	assert.True(t, mlrval.Equals(BIF_distinct_count(input), mlrval.FromInt(3)))
}

func TestBIF_null_count(t *testing.T) {
	// Needs array or map
	input := mlrval.FromInt(3)
	output := BIF_count(input)
	assert.True(t, output.IsError())

	input = mlrval.FromArray([]*mlrval.Mlrval{
		mlrval.FromInt(1),
		mlrval.FromString("two"),
		mlrval.FromString(""), // this counts
		mlrval.FromAnonymousError(),
		mlrval.ABSENT,
		mlrval.NULL, // this counts
	})
	assert.True(t, mlrval.Equals(BIF_null_count(input), mlrval.FromInt(2)))

	input = array_to_map_for_test(input)
	assert.True(t, mlrval.Equals(BIF_null_count(input), mlrval.FromInt(2)))

}

func TestBIF_mode_and_antimode(t *testing.T) {
	// Needs array or map
	input := mlrval.FromInt(3)
	output := BIF_count(input)
	assert.True(t, output.IsError())

	// Empty array
	input = mlrval.FromArray([]*mlrval.Mlrval{})
	assert.True(t, mlrval.Equals(BIF_mode(input), mlrval.VOID))
	assert.True(t, mlrval.Equals(BIF_antimode(input), mlrval.VOID))

	// Empty map
	input = array_to_map_for_test(input)
	assert.True(t, mlrval.Equals(BIF_mode(input), mlrval.VOID))
	assert.True(t, mlrval.Equals(BIF_antimode(input), mlrval.VOID))

	// Clear winner as array
	input = mlrval.FromArray([]*mlrval.Mlrval{
		mlrval.FromInt(1),
		mlrval.FromInt(2),
		mlrval.FromInt(3),
		mlrval.FromInt(1),
		mlrval.FromInt(1),
		mlrval.FromInt(2),
	})
	assert.True(t, mlrval.Equals(BIF_mode(input), mlrval.FromInt(1)))
	assert.True(t, mlrval.Equals(BIF_antimode(input), mlrval.FromInt(3)))

	// Clear winner as map
	input = array_to_map_for_test(input)
	assert.True(t, mlrval.Equals(BIF_mode(input), mlrval.FromInt(1)))
	assert.True(t, mlrval.Equals(BIF_antimode(input), mlrval.FromInt(3)))

	// Ties as array -- first-found breaks the tie
	input = mlrval.FromArray([]*mlrval.Mlrval{
		mlrval.FromInt(1),
		mlrval.FromInt(1),
		mlrval.FromInt(1),
		mlrval.FromInt(2),
		mlrval.FromInt(2),
		mlrval.FromInt(2),
	})
	assert.True(t, mlrval.Equals(BIF_mode(input), mlrval.FromInt(1)))
	assert.True(t, mlrval.Equals(BIF_antimode(input), mlrval.FromInt(1)))

	// Clear winner as map
	input = array_to_map_for_test(input)
	assert.True(t, mlrval.Equals(BIF_mode(input), mlrval.FromInt(1)))
	assert.True(t, mlrval.Equals(BIF_antimode(input), mlrval.FromInt(1)))
}

func TestBIF_sum(t *testing.T) {
	// Needs array or map
	input := mlrval.FromInt(3)
	output := BIF_count(input)
	assert.True(t, output.IsError())

	// TODO: test empty array/map
	for n := 1; n < 5; n++ {
		input = stats_test_array(n)
		var isum1 int64
		var isum2 int64
		var isum3 int64
		var isum4 int64
		for _, e := range input.AcquireArrayValue() {
			v := e.AcquireIntValue()
			isum1 += v
			isum2 += v * v
			isum3 += v * v * v
			isum4 += v * v * v * v
		}
		assert.True(t, mlrval.Equals(BIF_sum(input), mlrval.FromInt(isum1)))
		assert.True(t, mlrval.Equals(BIF_sum2(input), mlrval.FromInt(isum2)))
		assert.True(t, mlrval.Equals(BIF_sum3(input), mlrval.FromInt(isum3)))
		assert.True(t, mlrval.Equals(BIF_sum4(input), mlrval.FromInt(isum4)))

		input = array_to_map_for_test(input)
		assert.True(t, mlrval.Equals(BIF_sum(input), mlrval.FromInt(isum1)))
		assert.True(t, mlrval.Equals(BIF_sum2(input), mlrval.FromInt(isum2)))
		assert.True(t, mlrval.Equals(BIF_sum3(input), mlrval.FromInt(isum3)))
		assert.True(t, mlrval.Equals(BIF_sum4(input), mlrval.FromInt(isum4)))
	}
}

// More easily tested (much lower keystroking) within the regression-test framework:

// BIF_mean
// BIF_meaneb
// BIF_variance
// BIF_stddev
// BIF_skewness
// BIF_kurtosis

// BIF_min
// BIF_max

// BIF_minlen
// BIF_maxlen

// BIF_median
// BIF_median_with_options
// BIF_percentile
// BIF_percentile_with_options
// BIF_percentiles
// BIF_percentiles_with_options

// BIF_sort_collection
