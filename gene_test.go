package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewGene(t *testing.T) {
	tests := []struct {
		name     string
		output   geneID
		weight   float64
		disabled bool
		sum      float64
		activate activationFunction
	}{
		{
			name:     "output",
			output:   defaultOutput + 1,
			weight:   defaultWeight,
			disabled: defaultDisabled,
			activate: defaultActivationFunction,
		},
		{
			name:     "weight",
			output:   defaultOutput,
			weight:   defaultWeight + 1,
			disabled: defaultDisabled,
			activate: defaultActivationFunction,
		},
		{
			name:     "disabled",
			output:   defaultOutput,
			weight:   defaultWeight,
			disabled: !defaultDisabled,
			activate: defaultActivationFunction,
		},
		{
			name:     "activate",
			output:   defaultOutput,
			weight:   defaultWeight,
			disabled: defaultDisabled,
			activate: unit,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := newGene(withOutput(test.output),
				withWeight(test.weight),
				withDisabled(test.disabled),
				withActivationFunction(test.activate))

			require.Equal(t, g.output, test.output)
		})
	}
}

func TestGeneAdd(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		expect []float64
	}{
		{
			name:   "simple",
			values: []float64{0, 1, 2, 3, 4, 5},
			expect: []float64{0, 1, 3, 6, 10, 15},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := newGene(withActivationFunction(unit))

			for i, v := range test.values {
				g.add(v)
				require.Equal(t, g.val(), test.expect[i])
			}
		})
	}
}

func TestGeneClear(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "simple",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := newGene(withActivationFunction(unit))
			g.add(1)
			g.clear()

			require.Equal(t, g.val(), unit(0))
		})
	}
}

func TestGeneVal(t *testing.T) {
	tests := []struct {
		name     string
		activate activationFunction
		values   []float64
	}{
		{
			name:     "unit",
			activate: unit,
			values:   []float64{0, 1, 2, 3, 4, 5},
		},
		{
			name:     "sigmoid",
			activate: sigmoid,
			values:   []float64{0, 1, 2, 3, 4, 5},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := newGene(withActivationFunction(test.activate))

			for _, v := range test.values {
				g.clear()
				g.add(v)
				require.Equal(t, g.val(), test.activate(v))
			}
		})
	}
}

func TestGeneTerminal(t *testing.T) {
	tests := []struct {
		name   string
		output geneID
		expect bool
	}{
		{
			name:   "terminal",
			output: terminal,
			expect: true,
		},
		{
			name:   "not terminal",
			output: terminal + 1,
			expect: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := newGene(withOutput(test.output))

			require.Equal(t, g.terminal(), test.expect)
		})
	}
}
