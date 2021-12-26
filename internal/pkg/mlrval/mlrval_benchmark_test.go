package mlrval

import (
	"testing"
)

// go test -run=nonesuch -bench=. github.com/johnkerl/miller/internal/pkg/mlrval/...

func BenchmarkFromDeferredType(b *testing.B) {
    for i := 0; i < b.N; i++ {
		_ = FromDeferredType("123")
    }
}

func BenchmarkInferIntFromDeferredType(b *testing.B) {
    for i := 0; i < b.N; i++ {
		mv := FromDeferredType("123")
		mv.Type()
    }
}

func BenchmarkInferFloatFromDeferredType(b *testing.B) {
    for i := 0; i < b.N; i++ {
		mv := FromDeferredType("123.4")
		mv.Type()
    }
}

func BenchmarkInferStringFromDeferredType(b *testing.B) {
    for i := 0; i < b.N; i++ {
		mv := FromDeferredType("abc")
		mv.Type()
    }
}

