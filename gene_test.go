package neater

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewGene(t *testing.T) {
	tests := []struct {
		name     string
		p        nodePair
		weight   float64
		sum      float64
		activate activationFunction
	}{
		{
			name:     "input",
			p:        nodePair{1, 0},
			weight:   defaultWeight,
			activate: defaultActivationFunction,
		},
		{
			name:     "output",
			p:        nodePair{0, 1},
			weight:   defaultWeight,
			activate: defaultActivationFunction,
		},
		{
			name:     "weight",
			weight:   defaultWeight + 1,
			activate: defaultActivationFunction,
		},
		{
			name:     "disabled",
			weight:   defaultWeight,
			activate: defaultActivationFunction,
		},
		{
			name:     "activate",
			weight:   defaultWeight,
			activate: unit,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := newGene(test.p, test.weight, test.activate)

			require.Equal(t, g.p, test.p)
			require.Equal(t, g.weight, test.weight)
		})
	}
}
