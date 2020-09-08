package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEval(t *testing.T) {
	tests := []struct {
		name     string
		nInputs  int
		nOutputs int
		activate activationFunction

		input  []float64
		expect []float64
	}{
		{
			name:     "single unit",
			nInputs:  1,
			nOutputs: 1,
			activate: unit,

			input:  []float64{1},
			expect: []float64{defaultWeight * unit(1)},
		},
		{
			name:     "single sigmoid",
			nInputs:  1,
			nOutputs: 1,
			activate: sigmoid,

			input:  []float64{1},
			expect: []float64{defaultWeight * sigmoid(1)},
		},
		{
			name:     "double unit",
			nInputs:  2,
			nOutputs: 2,
			activate: unit,

			input:  []float64{1, 2},
			expect: []float64{defaultWeight * unit(1), defaultWeight * unit(2)},
		},
		{
			name:     "double sigmoid",
			nInputs:  2,
			nOutputs: 2,
			activate: sigmoid,

			input:  []float64{1, 2},
			expect: []float64{defaultWeight * sigmoid(1), defaultWeight * sigmoid(2)},
		},
		{
			name:     "single split",
			nInputs:  1,
			nOutputs: 2,
			activate: sigmoid,

			input:  []float64{1},
			expect: []float64{defaultWeight * sigmoid(1), defaultWeight * sigmoid(1)},
		},
		{
			name:     "single join",
			nInputs:  2,
			nOutputs: 1,
			activate: sigmoid,

			input:  []float64{1, 2},
			expect: []float64{defaultWeight*sigmoid(1) + defaultWeight*sigmoid(2)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := newOrganism(test.nInputs, test.nOutputs,
				withGlobalActivationFunction(test.activate),
				withConnectStrategy(connectFlow))

			output := o.Eval(test.input)
			require.Equal(t, test.expect, output)
		})
	}
}
