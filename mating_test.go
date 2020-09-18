package neat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func genCleanOrganism(inputs, outputs, nodes []nodeID, genes []*gene) *organism {
	o := newOrganism(&Configuration{
		Inputs:  len(inputs),
		Outputs: len(outputs),
	})
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

func TestRecombinate(t *testing.T) {
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
				newGene(nodePair{1, 2}, defaultWeight, unit),
				newGene(nodePair{1, 3}, defaultWeight, unit),
				newGene(nodePair{2, 4}, defaultWeight, unit),
				newGene(nodePair{3, 4}, defaultWeight, unit),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var s *species
			a := genCleanOrganism(test.inputs, test.outputs, test.nodes, test.commonGenes)
			b := a.copy()

			c := s.recombinate(a, b)
			require.Equal(t, a.strategy, c.strategy)

			for i, g := range a.oinnov {
				require.True(t, g.equalTo(c.oinnov[i]))
			}
		})
	}
}
