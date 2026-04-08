package bifs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnkerl/miller/v6/pkg/mlrval"
)

func TestBIF_length(t *testing.T) {
	input1 := mlrval.FromInt(123)
	output := BIF_length(input1)
	intval, ok := output.GetIntValue()
	assert.True(t, ok)
	assert.Equal(t, int64(1), intval)
}

func TestBIF_depth(t *testing.T) {
	input1 := mlrval.FromInt(123)
	output := BIF_depth(input1)
	intval, ok := output.GetIntValue()
	assert.True(t, ok)
	assert.Equal(t, int64(0), intval)

	mapval := mlrval.NewMlrmap()
	mapval.PutReference("key", mlrval.FromString("value"))
	input1 = mlrval.FromMap(mapval)
	output = BIF_depth(input1)
	intval, ok = output.GetIntValue()
	assert.True(t, ok)
	assert.Equal(t, int64(1), intval)

	arrayval := make([]*mlrval.Mlrval, 1)
	arrayval[0] = mlrval.FromString("value")
	input1 = mlrval.FromArray(arrayval)
	output = BIF_depth(input1)
	intval, ok = output.GetIntValue()
	assert.True(t, ok)
	assert.Equal(t, int64(1), intval)
}

func TestBIF_hasvalue(t *testing.T) {
	t.Run("array with existing value", func(t *testing.T) {
		arrayval := make([]*mlrval.Mlrval, 3)
		arrayval[0] = mlrval.FromInt(1)
		arrayval[1] = mlrval.FromInt(2)
		arrayval[2] = mlrval.FromInt(3)
		input1 := mlrval.FromArray(arrayval)
		input2 := mlrval.FromInt(2)

		output := BIF_hasvalue(input1, input2)
		boolval, ok := output.GetBoolValue()
		assert.True(t, ok)
		assert.True(t, boolval)
	})

	t.Run("array without existing value", func(t *testing.T) {
		arrayval := make([]*mlrval.Mlrval, 3)
		arrayval[0] = mlrval.FromInt(1)
		arrayval[1] = mlrval.FromInt(2)
		arrayval[2] = mlrval.FromInt(3)
		input1 := mlrval.FromArray(arrayval)
		input2 := mlrval.FromInt(5)

		output := BIF_hasvalue(input1, input2)
		boolval, ok := output.GetBoolValue()
		assert.True(t, ok)
		assert.False(t, boolval)
	})

	t.Run("map with existing value", func(t *testing.T) {
		mapval := mlrval.NewMlrmap()
		mapval.PutReference("a", mlrval.FromString("apple"))
		mapval.PutReference("b", mlrval.FromString("banana"))
		mapval.PutReference("c", mlrval.FromString("cherry"))
		input1 := mlrval.FromMap(mapval)
		input2 := mlrval.FromString("banana")

		output := BIF_hasvalue(input1, input2)
		boolval, ok := output.GetBoolValue()
		assert.True(t, ok)
		assert.True(t, boolval)
	})

	t.Run("map without existing value", func(t *testing.T) {
		mapval := mlrval.NewMlrmap()
		mapval.PutReference("a", mlrval.FromString("apple"))
		mapval.PutReference("b", mlrval.FromString("banana"))
		input1 := mlrval.FromMap(mapval)
		input2 := mlrval.FromString("orange")

		output := BIF_hasvalue(input1, input2)
		boolval, ok := output.GetBoolValue()
		assert.True(t, ok)
		assert.False(t, boolval)
	})

	t.Run("not a collection - should error", func(t *testing.T) {
		input1 := mlrval.FromInt(123)
		input2 := mlrval.FromInt(1)

		output := BIF_hasvalue(input1, input2)
		assert.False(t, output.IsBool())
	})
}

// TODO: copy in more unit-test cases from existing regression-test data

// func leafcount_from_array(input1 *mlrval.Mlrval) *mlrval.Mlrval
// func leafcount_from_map(input1 *mlrval.Mlrval) *mlrval.Mlrval
// func leafcount_from_scalar(input1 *mlrval.Mlrval) *mlrval.Mlrval
// func BIF_leafcount(input1 *mlrval.Mlrval) *mlrval.Mlrval
// func has_key_in_array(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
// func has_key_in_map(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
// func BIF_haskey(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
// func BIF_mapselect(mlrvals []*mlrval.Mlrval) *mlrval.Mlrval
// func BIF_mapexcept(mlrvals []*mlrval.Mlrval) *mlrval.Mlrval
// func BIF_mapsum(mlrvals []*mlrval.Mlrval) *mlrval.Mlrval
// func BIF_mapdiff(mlrvals []*mlrval.Mlrval) *mlrval.Mlrval
// func BIF_joink(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
// func BIF_joinv(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
// func BIF_joinkv(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval
// func BIF_splitkv(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval
// func BIF_splitkvx(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval
// func BIF_splitnv(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
// func BIF_splitnvx(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
// func BIF_splita(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
// func BIF_splitax(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
// func mlrvalSplitAXHelper(input string, separator string) *mlrval.Mlrval
// func BIF_get_keys(input1 *mlrval.Mlrval) *mlrval.Mlrval
// func BIF_get_values(input1 *mlrval.Mlrval) *mlrval.Mlrval
// func BIF_append(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
// func BIF_flatten(input1, input2, input3 *mlrval.Mlrval) *mlrval.Mlrval
// func BIF_flatten_binary(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
// func BIF_unflatten(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
// func BIF_arrayify(input1 *mlrval.Mlrval) *mlrval.Mlrval
// func BIF_json_parse(input1 *mlrval.Mlrval) *mlrval.Mlrval
// func BIF_json_stringify_unary(input1 *mlrval.Mlrval) *mlrval.Mlrval
// func BIF_json_stringify_binary(input1, input2 *mlrval.Mlrval) *mlrval.Mlrval
