package bifs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johnkerl/miller/internal/pkg/mlrval"
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
