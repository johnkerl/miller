package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test00(t *testing.T) {
	wk := NewWindowKeeper(0, 0)

	wk.Ingest("a")
	assert.Equal(t, "a", wk.Get(0).(string))

	wk.Ingest("b")
	assert.Equal(t, "b", wk.Get(0).(string))
}

func Test10(t *testing.T) {
	wk := NewWindowKeeper(1, 0)

	wk.Ingest("a")
	assert.Equal(t, "a", wk.Get(0).(string))
	assert.Equal(t, nil, wk.Get(-1))

	wk.Ingest("b")
	assert.Equal(t, "b", wk.Get(0).(string))
	assert.Equal(t, "a", wk.Get(-1).(string))

	wk.Ingest("c")
	assert.Equal(t, "c", wk.Get(0).(string))
	assert.Equal(t, "b", wk.Get(-1).(string))
}

func Test20(t *testing.T) {
	wk := NewWindowKeeper(2, 0)

	wk.Ingest("a")
	assert.Equal(t, "a", wk.Get(0).(string))
	assert.Equal(t, nil, wk.Get(-1))
	assert.Equal(t, nil, wk.Get(-2))

	wk.Ingest("b")
	assert.Equal(t, "b", wk.Get(0).(string))
	assert.Equal(t, "a", wk.Get(-1).(string))
	assert.Equal(t, nil, wk.Get(-2))

	wk.Ingest("c")
	assert.Equal(t, "c", wk.Get(0).(string))
	assert.Equal(t, "b", wk.Get(-1).(string))
	assert.Equal(t, "a", wk.Get(-2).(string))

	wk.Ingest("d")
	assert.Equal(t, "d", wk.Get(0).(string))
	assert.Equal(t, "c", wk.Get(-1).(string))
	assert.Equal(t, "b", wk.Get(-2).(string))
}

func Test01(t *testing.T) {
	wk := NewWindowKeeper(0, 1)

	wk.Ingest("a")
	assert.Equal(t, "a", wk.Get(1).(string))
	assert.Equal(t, nil, wk.Get(0))

	wk.Ingest("b")
	assert.Equal(t, "b", wk.Get(1).(string))
	assert.Equal(t, "a", wk.Get(0).(string))

	wk.Ingest("c")
	assert.Equal(t, "c", wk.Get(1).(string))
	assert.Equal(t, "b", wk.Get(0).(string))
}

func Test02(t *testing.T) {
	wk := NewWindowKeeper(0, 2)

	wk.Ingest("a")
	assert.Equal(t, "a", wk.Get(2).(string))
	assert.Equal(t, nil, wk.Get(1))
	assert.Equal(t, nil, wk.Get(0))

	wk.Ingest("b")
	assert.Equal(t, "b", wk.Get(2).(string))
	assert.Equal(t, "a", wk.Get(1).(string))
	assert.Equal(t, nil, wk.Get(0))

	wk.Ingest("c")
	assert.Equal(t, "c", wk.Get(2).(string))
	assert.Equal(t, "b", wk.Get(1).(string))
	assert.Equal(t, "a", wk.Get(0).(string))

	wk.Ingest("d")
	assert.Equal(t, "d", wk.Get(2).(string))
	assert.Equal(t, "c", wk.Get(1).(string))
	assert.Equal(t, "b", wk.Get(0).(string))
}

func Test32(t *testing.T) {
	wk := NewWindowKeeper(3, 2)

	wk.Ingest("a")
	assert.Equal(t, "a", wk.Get(2).(string))
	assert.Equal(t, nil, wk.Get(1))
	assert.Equal(t, nil, wk.Get(0))
	assert.Equal(t, nil, wk.Get(-1))
	assert.Equal(t, nil, wk.Get(-2))
	assert.Equal(t, nil, wk.Get(-3))

	wk.Ingest("b")
	assert.Equal(t, "b", wk.Get(2).(string))
	assert.Equal(t, "a", wk.Get(1).(string))
	assert.Equal(t, nil, wk.Get(0))
	assert.Equal(t, nil, wk.Get(-1))
	assert.Equal(t, nil, wk.Get(-2))
	assert.Equal(t, nil, wk.Get(-3))

	wk.Ingest("c")
	assert.Equal(t, "c", wk.Get(2).(string))
	assert.Equal(t, "b", wk.Get(1).(string))
	assert.Equal(t, "a", wk.Get(0).(string))
	assert.Equal(t, nil, wk.Get(-1))
	assert.Equal(t, nil, wk.Get(-2))
	assert.Equal(t, nil, wk.Get(-3))

	wk.Ingest("d")
	assert.Equal(t, "d", wk.Get(2).(string))
	assert.Equal(t, "c", wk.Get(1).(string))
	assert.Equal(t, "b", wk.Get(0).(string))
	assert.Equal(t, "a", wk.Get(-1).(string))
	assert.Equal(t, nil, wk.Get(-2))
	assert.Equal(t, nil, wk.Get(-3))

	wk.Ingest("e")
	assert.Equal(t, "e", wk.Get(2).(string))
	assert.Equal(t, "d", wk.Get(1).(string))
	assert.Equal(t, "c", wk.Get(0).(string))
	assert.Equal(t, "b", wk.Get(-1).(string))
	assert.Equal(t, "a", wk.Get(-2).(string))
	assert.Equal(t, nil, wk.Get(-3))

	wk.Ingest("f")
	assert.Equal(t, "f", wk.Get(2).(string))
	assert.Equal(t, "e", wk.Get(1).(string))
	assert.Equal(t, "d", wk.Get(0).(string))
	assert.Equal(t, "c", wk.Get(-1).(string))
	assert.Equal(t, "b", wk.Get(-2).(string))
	assert.Equal(t, "a", wk.Get(-3).(string))

	wk.Ingest("g")
	assert.Equal(t, "g", wk.Get(2).(string))
	assert.Equal(t, "f", wk.Get(1).(string))
	assert.Equal(t, "e", wk.Get(0).(string))
	assert.Equal(t, "d", wk.Get(-1).(string))
	assert.Equal(t, "c", wk.Get(-2).(string))
	assert.Equal(t, "b", wk.Get(-3).(string))
}
