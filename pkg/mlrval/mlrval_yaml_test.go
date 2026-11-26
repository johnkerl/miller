// Tests for YAML decode/encode (MlrvalDecodeFromYAML, MlrmapToYAMLNative).

package mlrval

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.yaml.in/yaml/v3"
)

func TestMlrvalFromYAMLNativeScalars(t *testing.T) {
	mv, err := mlrvalFromYAMLNative("hello")
	assert.NoError(t, err)
	assert.True(t, mv.IsString())
	s, _ := mv.GetStringValue()
	assert.Equal(t, "hello", s)

	mv, err = mlrvalFromYAMLNative(42)
	assert.NoError(t, err)
	assert.True(t, mv.IsInt())
	i, _ := mv.GetIntValue()
	assert.Equal(t, int64(42), i)

	mv, err = mlrvalFromYAMLNative(3.14)
	assert.NoError(t, err)
	assert.True(t, mv.IsFloat())
	f, _ := mv.GetFloatValue()
	assert.Equal(t, 3.14, f)

	mv, err = mlrvalFromYAMLNative(true)
	assert.NoError(t, err)
	assert.True(t, mv.IsBool())
	b, _ := mv.GetBoolValue()
	assert.True(t, b)

	mv, err = mlrvalFromYAMLNative(nil)
	assert.NoError(t, err)
	assert.True(t, mv.IsNull())
}

func TestMlrvalFromYAMLNativeMap(t *testing.T) {
	// gopkg.in/yaml.v3 unmarshals into map[string]interface{}
	m := map[string]interface{}{
		"a": 1,
		"b": "two",
		"c": true,
	}
	mv, err := mlrvalFromYAMLNative(m)
	assert.NoError(t, err)
	assert.True(t, mv.IsMap())
	rec := mv.GetMap()
	assert.NotNil(t, rec)
	assert.Equal(t, int64(3), rec.FieldCount)
	av := rec.Get("a")
	assert.True(t, av.IsInt())
	ai, _ := av.GetIntValue()
	assert.Equal(t, int64(1), ai)
	bv := rec.Get("b")
	assert.True(t, bv.IsString())
	bs, _ := bv.GetStringValue()
	assert.Equal(t, "two", bs)
}

func TestMlrvalDecodeFromYAML(t *testing.T) {
	input := "a: 1\nb: 2\n"
	decoder := yaml.NewDecoder(bytes.NewBufferString(input))
	mv, eof, err := MlrvalDecodeFromYAML(decoder)
	assert.NoError(t, err)
	assert.False(t, eof)
	assert.NotNil(t, mv)
	assert.True(t, mv.IsMap())
	rec := mv.GetMap()
	assert.Equal(t, int64(2), rec.FieldCount)
	v := rec.Get("a")
	vi, _ := v.GetIntValue()
	assert.Equal(t, int64(1), vi)
	v = rec.Get("b")
	vi, _ = v.GetIntValue()
	assert.Equal(t, int64(2), vi)

	// EOF on next decode
	_, eof, err = MlrvalDecodeFromYAML(decoder)
	assert.NoError(t, err)
	assert.True(t, eof)
}

func TestMlrmapToYAMLNativeRoundTrip(t *testing.T) {
	rec := NewMlrmapAsRecord()
	rec.PutReference("x", FromInt(10))
	rec.PutReference("y", FromString("hello"))
	rec.PutReference("z", FromBool(true))

	native, err := MlrmapToYAMLNative(rec)
	assert.NoError(t, err)
	assert.NotNil(t, native)

	out, err := yaml.Marshal(native)
	assert.NoError(t, err)
	decoder := yaml.NewDecoder(bytes.NewReader(out))
	mv, eof, err := MlrvalDecodeFromYAML(decoder)
	assert.NoError(t, err)
	assert.False(t, eof)
	assert.True(t, mv.IsMap())
	back := mv.GetMap()
	xv := back.Get("x")
	xi, _ := xv.GetIntValue()
	assert.Equal(t, int64(10), xi)
	yv := back.Get("y")
	ys, _ := yv.GetStringValue()
	assert.Equal(t, "hello", ys)
	zv := back.Get("z")
	zb, _ := zv.GetBoolValue()
	assert.True(t, zb)
}

func TestMlrvalDecodeFromYAMLArrayOfMaps(t *testing.T) {
	input := "- a: 1\n  b: 2\n- a: 3\n  b: 4\n"
	decoder := yaml.NewDecoder(bytes.NewBufferString(input))
	mv, eof, err := MlrvalDecodeFromYAML(decoder)
	assert.NoError(t, err)
	assert.False(t, eof)
	assert.True(t, mv.IsArray())
	arr := mv.GetArray()
	assert.Len(t, arr, 2)
	assert.True(t, arr[0].IsMap())
	assert.True(t, arr[1].IsMap())
	r0 := arr[0].GetMap()
	av := r0.Get("a")
	ai, _ := av.GetIntValue()
	assert.Equal(t, int64(1), ai)
	r1 := arr[1].GetMap()
	av = r1.Get("a")
	ai, _ = av.GetIntValue()
	assert.Equal(t, int64(3), ai)
}

func TestMlrvalDecodeFromYAMLMultiDoc(t *testing.T) {
	input := "a: 1\nb: 2\n---\nc: 3\nd: 4\n"
	decoder := yaml.NewDecoder(bytes.NewBufferString(input))

	mv1, eof, err := MlrvalDecodeFromYAML(decoder)
	assert.NoError(t, err)
	assert.False(t, eof)
	assert.True(t, mv1.IsMap())
	assert.Equal(t, int64(2), mv1.GetMap().FieldCount)

	mv2, eof, err := MlrvalDecodeFromYAML(decoder)
	assert.NoError(t, err)
	assert.False(t, eof)
	assert.True(t, mv2.IsMap())
	assert.Equal(t, int64(2), mv2.GetMap().FieldCount)

	_, eof, err = MlrvalDecodeFromYAML(decoder)
	assert.NoError(t, err)
	assert.True(t, eof)
}

func TestYAMLKeyStringNonStringKeys(t *testing.T) {
	// map[interface{}]interface{} with int key
	m := map[interface{}]interface{}{
		1:     "one",
		"two": 2,
	}
	mv, err := mlrvalFromYAMLMap(m)
	assert.NoError(t, err)
	assert.True(t, mv.IsMap())
	rec := mv.GetMap()
	v := rec.Get("1")
	assert.True(t, v.IsString())
	s, _ := v.GetStringValue()
	assert.Equal(t, "one", s)
	v = rec.Get("two")
	assert.True(t, v.IsInt())
}

func TestMlrmapToYAMLNativeKeyOrder(t *testing.T) {
	rec := NewMlrmapAsRecord()
	rec.PutReference("b", FromInt(22))
	rec.PutReference("c", FromInt(33))
	rec.PutReference("a", FromInt(11))

	native, err := MlrmapToYAMLNative(rec)
	assert.NoError(t, err)
	out, err := yaml.Marshal(native)
	assert.NoError(t, err)

	var node yaml.Node
	err = yaml.Unmarshal(out, &node)
	assert.NoError(t, err)
	if !assert.Len(t, node.Content, 1, "expected a single YAML document node") {
		return
	}

	mapping := node.Content[0]
	if !assert.Equal(t, yaml.MappingNode, mapping.Kind, "expected top-level YAML mapping") {
		return
	}

	// Keys must appear in insertion order (b, c, a), not sorted (a, b, c).
	if !assert.GreaterOrEqual(t, len(mapping.Content), 6, "expected at least three key/value pairs") {
		return
	}
	assert.Equal(t, "b", mapping.Content[0].Value)
	assert.Equal(t, "c", mapping.Content[2].Value)
	assert.Equal(t, "a", mapping.Content[4].Value)
}

func TestMlrmapToYAMLNativeNested(t *testing.T) {
	inner := NewMlrmapAsRecord()
	inner.PutReference("p", FromInt(1))
	inner.PutReference("q", FromInt(2))
	rec := NewMlrmapAsRecord()
	rec.PutReference("a", FromString("x"))
	rec.PutReference("nested", FromMap(inner))

	native, err := MlrmapToYAMLNative(rec)
	assert.NoError(t, err)
	out, err := yaml.Marshal(native)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(out), "nested:"))
	assert.True(t, strings.Contains(string(out), "p: 1"))
}

func yamlDecodedKeys(t *testing.T, mv *Mlrval) []string {
	t.Helper()
	if !assert.True(t, mv.IsMap()) {
		return nil
	}
	keys := []string{}
	for pe := mv.GetMap().Head; pe != nil; pe = pe.Next {
		keys = append(keys, pe.Key)
	}
	return keys
}

func TestMlrvalDecodeFromYAMLKeyOrder(t *testing.T) {
	input := "b: 22\nc:\n  z: 1\n  y: 2\na: [{q: 1, p: 2}]\n"
	decoder := yaml.NewDecoder(bytes.NewBufferString(input))
	mv, eof, err := MlrvalDecodeFromYAML(decoder)
	assert.NoError(t, err)
	assert.False(t, eof)
	assert.Equal(t, []string{"b", "c", "a"}, yamlDecodedKeys(t, mv))
	assert.Equal(t, []string{"z", "y"}, yamlDecodedKeys(t, mv.GetMap().Get("c")))
	arr := mv.GetMap().Get("a").GetArray()
	if assert.Len(t, arr, 1) {
		assert.Equal(t, []string{"q", "p"}, yamlDecodedKeys(t, arr[0]))
	}
}

func TestMlrvalDecodeFromYAMLKeyOrderArrayOfMaps(t *testing.T) {
	input := "- z: 1\n  a: 2\n- y: 3\n  b: 4\n"
	decoder := yaml.NewDecoder(bytes.NewBufferString(input))
	mv, _, err := MlrvalDecodeFromYAML(decoder)
	assert.NoError(t, err)
	arr := mv.GetArray()
	if assert.Len(t, arr, 2) {
		assert.Equal(t, []string{"z", "a"}, yamlDecodedKeys(t, arr[0]))
		assert.Equal(t, []string{"y", "b"}, yamlDecodedKeys(t, arr[1]))
	}
}

func TestMlrvalDecodeFromYAMLKeyOrderWithAliasAndMerge(t *testing.T) {
	input := "base: &base\n  y: 1\n  x: 2\ncopy: *base\nmerged:\n  w: 0\n  <<: *base\n  x: 9\n"
	decoder := yaml.NewDecoder(bytes.NewBufferString(input))
	mv, _, err := MlrvalDecodeFromYAML(decoder)
	assert.NoError(t, err)
	rec := mv.GetMap()
	assert.Equal(t, []string{"base", "copy", "merged"}, yamlDecodedKeys(t, mv))
	assert.Equal(t, []string{"y", "x"}, yamlDecodedKeys(t, rec.Get("copy")))
	merged := rec.Get("merged")
	assert.Equal(t, []string{"w", "y", "x"}, yamlDecodedKeys(t, merged))
	x, _ := merged.GetMap().Get("x").GetIntValue()
	assert.Equal(t, int64(9), x)
}

func TestMlrvalDecodeFromYAMLDuplicateKeyIsError(t *testing.T) {
	decoder := yaml.NewDecoder(bytes.NewBufferString("a: 1\na: 2\n"))
	_, _, err := MlrvalDecodeFromYAML(decoder)
	assert.Error(t, err)
}
