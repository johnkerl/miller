package lib

import (
	"math/rand"
	"os"
	"time"
)

// ================================================================
// Thinly wraps Go's rand library, with seed-function support
// ================================================================

// By default, Miller random numbers are different on every run.
var defaultSeed = time.Now().UnixNano() ^ int64(os.Getpid())
var source = rand.NewSource(defaultSeed)
var generator = rand.New(source)

// Users can request specific seeds if they want the same random-number
// sequence on each run.
func SeedRandom(seed int64) {
	source = rand.NewSource(seed)
	generator = rand.New(source)
}

func RandFloat64() float64 {
	return generator.Float64()
}
func RandUint32() uint32 {
	return generator.Uint32()
}
