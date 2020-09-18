package neat

import (
	"math/rand"
	"testing"
)

func BenchmarkIntn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = rand.Intn(100)
	}
}

func BenchmarkFloat64AndMultiplication(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = rand.Float64() * 2.5
	}
}
