package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEval(t *testing.T) {
	tests := []struct {
		name     string
		conf     *Configuration
		activate activationFunction

		input  []float64
		expect []float64
	}{
		{
			name: "single unit",
			conf: &Configuration{
				Inputs:  1,
				Outputs: 1,
			},
			activate: unit,

			input:  []float64{1},
			expect: []float64{defaultWeight * unit(1)},
		},
		{
			name: "single sigmoid",
			conf: &Configuration{
				Inputs:  1,
				Outputs: 1,
			},
			activate: sigmoid,

			input:  []float64{1},
			expect: []float64{defaultWeight * sigmoid(1)},
		},
		{
			name: "double unit",
			conf: &Configuration{
				Inputs:  2,
				Outputs: 2,
			},
			activate: unit,

			input:  []float64{1, 2},
			expect: []float64{defaultWeight * unit(1), defaultWeight * unit(2)},
		},
		{
			name: "double sigmoid",
			conf: &Configuration{
				Inputs:  2,
				Outputs: 2,
			},
			activate: sigmoid,

			input:  []float64{1, 2},
			expect: []float64{defaultWeight * sigmoid(1), defaultWeight * sigmoid(2)},
		},
		{
			name: "single split",
			conf: &Configuration{
				Inputs:  1,
				Outputs: 2,
			},
			activate: sigmoid,

			input:  []float64{1},
			expect: []float64{defaultWeight * sigmoid(1), defaultWeight * sigmoid(1)},
		},
		{
			name: "single join",
			conf: &Configuration{
				Inputs:  2,
				Outputs: 1,
			},
			activate: sigmoid,

			input:  []float64{1, 2},
			expect: []float64{defaultWeight*sigmoid(1) + defaultWeight*sigmoid(2)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := newOrganism(test.conf,
				withGlobalActivationFunction(test.activate),
				withConnectStrategy(connectFlow))

			output := o.Eval(test.input)
			require.Equal(t, test.expect, output)
		})
	}
}
