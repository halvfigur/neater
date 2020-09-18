package neat

import (
	"math"
	"math/rand"
)

var (
	defaultRandIntn = rand.Intn
	randIntn        = defaultRandIntn

	defaultRandFloat64 = rand.Float64
	randFloat64        = defaultRandFloat64
)

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func unit(x float64) float64 {
	return x
}

func sigmoid(x float64) float64 {
	return float64(1) / (float64(1) + math.Exp(x))
}
