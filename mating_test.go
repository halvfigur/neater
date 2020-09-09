package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func genCleanOrganism(inputs, outputs, nodes []nodeID, genes []*gene) *organism {
	o := newOrganism(len(inputs), len(outputs))
	o.inputs = inputs
	o.outputs = outputs

	for _, id := range inputs {
		o.nodes[id] = 0
	}

	for _, id := range outputs {
		o.nodes[id] = 0
	}

	for _, id := range nodes {
		o.nodes[id] = 0
	}

	for _, g := range genes {
		o.add(g)
	}

	return o
}

func hasGene(o *organism, g *gene) bool {
	return false
}

func TestMating(t *testing.T) {
	tests := []struct {
		name        string
		inputs      []nodeID
		outputs     []nodeID
		nodes       []nodeID
		alphaScore  float64
		betaScore   float64
		commonGenes []*gene
		alphaGenes  []*gene
		betaGenes   []*gene
	}{
		{
			name:       "sune",
			inputs:     []nodeID{nodeID(1)},
			outputs:    []nodeID{nodeID(4)},
			nodes:      []nodeID{nodeID(2), nodeID(3)},
			alphaScore: 1,
			betaScore:  0,
			commonGenes: []*gene{
				newGene(1, 2),
				newGene(1, 3),
				newGene(2, 4),
				newGene(3, 4),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := genCleanOrganism(test.inputs, test.outputs, test.nodes, test.commonGenes)
			b := a.copy()

			c := mate(a, b)
			require.Equal(t, a.recurrence, c.recurrence)
			require.Equal(t, a.strategy, c.strategy)

			for i, g := range a.oinnov {
				require.True(t, g.equalTo(c.oinnov[i]))
			}
		})
	}
}
