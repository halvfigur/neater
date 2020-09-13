package main

import (
	"fmt"
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

func TestAdd(t *testing.T) {
	var nCount uint64
	nodeIDGenerator = func() nodeID {
		nCount++
		return nodeID(nCount)
	}

	tests := []struct {
		name     string
		conf     *Configuration
		activate activationFunction
		pairs    []nodePair
		input    []float64
		output   []float64
		expect   []float64
	}{
		{
			//
			//      +--> 3 --+
			//		|        |
			//		|        v
			//	1 --+------> 2
			name: "I(1), O(2), 1->3, 1->2, 3->2",
			conf: &Configuration{
				Inputs:  1,
				Outputs: 1,
			},
			pairs: []nodePair{
				nodePair{1, 2},
				nodePair{1, 3},
				nodePair{3, 2},
			},
			input:  []float64{1},
			expect: []float64{2},
		},
		{
			//	1 ---+        +---> 3
			//		 |        |
			//		 +--> 5 --+
			//		 |        |
			//	2 ---+        +---> 4
			name: "testing",
			conf: &Configuration{
				Inputs:  2,
				Outputs: 2,
			},
			pairs: []nodePair{
				nodePair{1, 5},
				nodePair{5, 3},
				nodePair{2, 5},
				nodePair{5, 4},
			},
			input:  []float64{1, 2},
			expect: []float64{3, 3},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := newOrganism(test.conf,
				withGlobalActivationFunction(unit),
				withConnectStrategy(connectNone))

			// Reset node ID counter
			nCount = 0
			for _, p := range test.pairs {

				o.nodes[p.input] = 0
				o.nodes[p.output] = 0

				g := newGene(p, withActivationFunction(unit))
				fmt.Printf("Insert: %v\n", g.p)
				o.add(g)

				fmt.Print("After: ")
				for _, g := range o.oeval {
					fmt.Printf("%v ", g.p)
				}
				fmt.Println()
				fmt.Println()
			}

			output := o.Eval(test.input)
			require.Equal(t, test.expect, output)
		})
	}
}
