// ================================================================
// Tests for YAML decode/encode (MlrvalDecodeFromYAML, MlrmapToYAMLNative).
// ================================================================

package mlrval

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"gopkg.in/yaml.v3"
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
