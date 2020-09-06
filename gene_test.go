package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewGene(t *testing.T) {
	tests := []struct {
		name     string
		input    nodeID
		output   nodeID
		weight   float64
		disabled bool
		sum      float64
		activate activationFunction
	}{
		{
			name:     "input",
			input:    1,
			weight:   defaultWeight,
			disabled: defaultDisabled,
			activate: defaultActivationFunction,
		},
		{
			name:     "output",
			output:   1,
			weight:   defaultWeight,
			disabled: defaultDisabled,
			activate: defaultActivationFunction,
		},
		{
			name:     "weight",
			weight:   defaultWeight + 1,
			disabled: defaultDisabled,
			activate: defaultActivationFunction,
		},
		{
			name:     "disabled",
			weight:   defaultWeight,
			disabled: !defaultDisabled,
			activate: defaultActivationFunction,
		},
		{
			name:     "activate",
			weight:   defaultWeight,
			disabled: defaultDisabled,
			activate: unit,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := newGene(test.input, test.output,
				withWeight(test.weight),
				withDisabled(test.disabled),
				withActivationFunction(test.activate))

			require.Equal(t, g.output, test.output)
		})
	}
}
