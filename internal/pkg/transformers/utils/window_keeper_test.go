package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test00(t *testing.T) {
	wk := NewWindowKeeper(0, 0)

	wk.IngestRecord("a")
	assert.Equal(t, "a", wk.GetRecord(0).(string))

	wk.IngestRecord("b")
	assert.Equal(t, "b", wk.GetRecord(0).(string))
}

func Test10(t *testing.T) {
	wk := NewWindowKeeper(1, 0)

	wk.IngestRecord("a")
	assert.Equal(t, "a", wk.GetRecord(0).(string))
	assert.Equal(t, nil, wk.GetRecord(-1))

	wk.IngestRecord("b")
	assert.Equal(t, "b", wk.GetRecord(0).(string))
	assert.Equal(t, "a", wk.GetRecord(-1).(string))

	wk.IngestRecord("c")
	assert.Equal(t, "c", wk.GetRecord(0).(string))
	assert.Equal(t, "b", wk.GetRecord(-1).(string))
}

func Test20(t *testing.T) {
	wk := NewWindowKeeper(2, 0)

	wk.IngestRecord("a")
	assert.Equal(t, "a", wk.GetRecord(0).(string))
	assert.Equal(t, nil, wk.GetRecord(-1))
	assert.Equal(t, nil, wk.GetRecord(-2))

	wk.IngestRecord("b")
	assert.Equal(t, "b", wk.GetRecord(0).(string))
	assert.Equal(t, "a", wk.GetRecord(-1).(string))
	assert.Equal(t, nil, wk.GetRecord(-2))

	wk.IngestRecord("c")
	assert.Equal(t, "c", wk.GetRecord(0).(string))
	assert.Equal(t, "b", wk.GetRecord(-1).(string))
	assert.Equal(t, "a", wk.GetRecord(-2).(string))

	wk.IngestRecord("d")
	assert.Equal(t, "d", wk.GetRecord(0).(string))
	assert.Equal(t, "c", wk.GetRecord(-1).(string))
	assert.Equal(t, "b", wk.GetRecord(-2).(string))
}

func Test01(t *testing.T) {
	wk := NewWindowKeeper(0, 1)

	wk.IngestRecord("a")
	assert.Equal(t, "a", wk.GetRecord(1).(string))
	assert.Equal(t, nil, wk.GetRecord(0))

	wk.IngestRecord("b")
	assert.Equal(t, "b", wk.GetRecord(1).(string))
	assert.Equal(t, "a", wk.GetRecord(0).(string))

	wk.IngestRecord("c")
	assert.Equal(t, "c", wk.GetRecord(1).(string))
	assert.Equal(t, "b", wk.GetRecord(0).(string))
}

func Test02(t *testing.T) {
	wk := NewWindowKeeper(0, 2)

	wk.IngestRecord("a")
	assert.Equal(t, "a", wk.GetRecord(2).(string))
	assert.Equal(t, nil, wk.GetRecord(1))
	assert.Equal(t, nil, wk.GetRecord(0))

	wk.IngestRecord("b")
	assert.Equal(t, "b", wk.GetRecord(2).(string))
	assert.Equal(t, "a", wk.GetRecord(1).(string))
	assert.Equal(t, nil, wk.GetRecord(0))

	wk.IngestRecord("c")
	assert.Equal(t, "c", wk.GetRecord(2).(string))
	assert.Equal(t, "b", wk.GetRecord(1).(string))
	assert.Equal(t, "a", wk.GetRecord(0).(string))

	wk.IngestRecord("d")
	assert.Equal(t, "d", wk.GetRecord(2).(string))
	assert.Equal(t, "c", wk.GetRecord(1).(string))
	assert.Equal(t, "b", wk.GetRecord(0).(string))
}

func Test32(t *testing.T) {
	wk := NewWindowKeeper(3, 2)

	wk.IngestRecord("a")
	assert.Equal(t, "a", wk.GetRecord(2).(string))
	assert.Equal(t, nil, wk.GetRecord(1))
	assert.Equal(t, nil, wk.GetRecord(0))
	assert.Equal(t, nil, wk.GetRecord(-1))
	assert.Equal(t, nil, wk.GetRecord(-2))
	assert.Equal(t, nil, wk.GetRecord(-3))

	wk.IngestRecord("b")
	assert.Equal(t, "b", wk.GetRecord(2).(string))
	assert.Equal(t, "a", wk.GetRecord(1).(string))
	assert.Equal(t, nil, wk.GetRecord(0))
	assert.Equal(t, nil, wk.GetRecord(-1))
	assert.Equal(t, nil, wk.GetRecord(-2))
	assert.Equal(t, nil, wk.GetRecord(-3))

	wk.IngestRecord("c")
	assert.Equal(t, "c", wk.GetRecord(2).(string))
	assert.Equal(t, "b", wk.GetRecord(1).(string))
	assert.Equal(t, "a", wk.GetRecord(0).(string))
	assert.Equal(t, nil, wk.GetRecord(-1))
	assert.Equal(t, nil, wk.GetRecord(-2))
	assert.Equal(t, nil, wk.GetRecord(-3))

	wk.IngestRecord("d")
	assert.Equal(t, "d", wk.GetRecord(2).(string))
	assert.Equal(t, "c", wk.GetRecord(1).(string))
	assert.Equal(t, "b", wk.GetRecord(0).(string))
	assert.Equal(t, "a", wk.GetRecord(-1).(string))
	assert.Equal(t, nil, wk.GetRecord(-2))
	assert.Equal(t, nil, wk.GetRecord(-3))

	wk.IngestRecord("e")
	assert.Equal(t, "e", wk.GetRecord(2).(string))
	assert.Equal(t, "d", wk.GetRecord(1).(string))
	assert.Equal(t, "c", wk.GetRecord(0).(string))
	assert.Equal(t, "b", wk.GetRecord(-1).(string))
	assert.Equal(t, "a", wk.GetRecord(-2).(string))
	assert.Equal(t, nil, wk.GetRecord(-3))

	wk.IngestRecord("f")
	assert.Equal(t, "f", wk.GetRecord(2).(string))
	assert.Equal(t, "e", wk.GetRecord(1).(string))
	assert.Equal(t, "d", wk.GetRecord(0).(string))
	assert.Equal(t, "c", wk.GetRecord(-1).(string))
	assert.Equal(t, "b", wk.GetRecord(-2).(string))
	assert.Equal(t, "a", wk.GetRecord(-3).(string))

	wk.IngestRecord("g")
	assert.Equal(t, "g", wk.GetRecord(2).(string))
	assert.Equal(t, "f", wk.GetRecord(1).(string))
	assert.Equal(t, "e", wk.GetRecord(0).(string))
	assert.Equal(t, "d", wk.GetRecord(-1).(string))
	assert.Equal(t, "c", wk.GetRecord(-2).(string))
	assert.Equal(t, "b", wk.GetRecord(-3).(string))
}
