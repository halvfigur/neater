package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestXorTrainer(t *testing.T) {
	inputs := [][]float64{
		{0, 0},
		{0, 1},
		{1, 0},
		{1, 1},
	}

	outputs := [][]float64{
		{0},
		{1},
		{1},
		{0},
	}

	c := NewXORFitnessCalculator()

	for i := range inputs {
		c.AddResult(inputs[i], outputs[i])
	}

	require.Equal(t, float64(1), c.CalculateFitness())
}
