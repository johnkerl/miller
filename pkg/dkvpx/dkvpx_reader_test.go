package dkvpx

import (
	"io"
	"strings"
	"testing"

	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/stretchr/testify/assert"
)

// orderedMapToMap converts OrderedMap to map for easier assertion (order-agnostic).
func orderedMapToMap(om *lib.OrderedMap[string]) map[string]string {
	if om == nil {
		return nil
	}
	m := make(map[string]string)
	for pe := om.Head; pe != nil; pe = pe.Next {
		m[pe.Key] = pe.Value
	}
	return m
}

func TestRead_Basic(t *testing.T) {
	r := NewReader(strings.NewReader("x=1,y=2,z=3\n"))
	rec, err := r.Read()
	assert.NoError(t, err)
	assert.NotNil(t, rec)
	assert.Equal(t, map[string]string{"x": "1", "y": "2", "z": "3"}, orderedMapToMap(rec))

	rec, err = r.Read()
	assert.Nil(t, rec)
	assert.ErrorIs(t, err, io.EOF)
}

func TestRead_ImplicitKeys(t *testing.T) {
	r := NewReader(strings.NewReader("x=1,2,z=3\n"))
	rec, err := r.Read()
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"x": "1", "2": "2", "z": "3"}, orderedMapToMap(rec))
}

func TestRead_QuotedKeyAndValue(t *testing.T) {
	r := NewReader(strings.NewReader(`"x,y"="a,b,c",z=3` + "\n"))
	rec, err := r.Read()
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"x,y": "a,b,c", "z": "3"}, orderedMapToMap(rec))
}

func TestRead_EscapedQuotes(t *testing.T) {
	r := NewReader(strings.NewReader(`x="the ""word""",y=2` + "\n"))
	rec, err := r.Read()
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"x": `the "word"`, "y": "2"}, orderedMapToMap(rec))
}

func TestRead_MultilineValue(t *testing.T) {
	r := NewReader(strings.NewReader("x=\"a\n b\",y=2\n"))
	rec, err := r.Read()
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"x": "a\n b", "y": "2"}, orderedMapToMap(rec))
}

func TestRead_EmptyLine(t *testing.T) {
	r := NewReader(strings.NewReader("\n"))
	rec, err := r.Read()
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{}, orderedMapToMap(rec))
}

func TestRead_EmptyEOF(t *testing.T) {
	r := NewReader(strings.NewReader(""))
	rec, err := r.Read()
	assert.Nil(t, rec)
	assert.ErrorIs(t, err, io.EOF)
}

func TestRead_ReadAll(t *testing.T) {
	r := NewReader(strings.NewReader("a=1,b=2\nx=10,y=20\n"))
	rec1, err := r.Read()
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"a": "1", "b": "2"}, orderedMapToMap(rec1))

	rec2, err := r.Read()
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"x": "10", "y": "20"}, orderedMapToMap(rec2))

	rec3, err := r.Read()
	assert.Nil(t, rec3)
	assert.ErrorIs(t, err, io.EOF)
}

func TestRead_ValueWithEquals(t *testing.T) {
	r := NewReader(strings.NewReader("a=b=c,d=1\n"))
	rec, err := r.Read()
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"a": "b=c", "d": "1"}, orderedMapToMap(rec))
}

func TestRead_CommentSkipped(t *testing.T) {
	r := NewReader(strings.NewReader("# comment\na=1,b=2\n"))
	r.Comment = '#'
	rec, err := r.Read()
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"a": "1", "b": "2"}, orderedMapToMap(rec))
}

func TestRead_OrderPreserved(t *testing.T) {
	r := NewReader(strings.NewReader("z=3,x=1,y=2\n"))
	rec, err := r.Read()
	assert.NoError(t, err)
	// OrderedMap should preserve insertion order
	keys := make([]string, 0)
	for pe := rec.Head; pe != nil; pe = pe.Next {
		keys = append(keys, pe.Key)
	}
	assert.Equal(t, []string{"z", "x", "y"}, keys)
}
