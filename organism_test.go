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

			input:  []float64{0},
			expect: []float64{unit(0)},
		},
		{
			name:     "single sigmoid",
			nInputs:  1,
			nOutputs: 1,
			activate: sigmoid,

			input:  []float64{0},
			expect: []float64{sigmoid(0)},
		},
		{
			name:     "double unit",
			nInputs:  2,
			nOutputs: 2,
			activate: unit,

			input:  []float64{0, 1},
			expect: []float64{unit(0), unit(1)},
		},
		{
			name:     "double sigmoid",
			nInputs:  2,
			nOutputs: 2,
			activate: sigmoid,

			input:  []float64{0, 1},
			expect: []float64{sigmoid(0), sigmoid(1)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := newOrganism(test.nInputs, test.nOutputs,
				withGlobalActivationFunction(test.activate),
				withConnectStrategy(connectFlow))

			output := o.Eval(test.input)
			require.Equal(t, output, test.expect)
		})
	}
}
