package neater

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecombinate(t *testing.T) {
	newGene := func(p nodePair, w float64, f activationFunction, innov geneID) *gene {
		return &gene{
			innov:    innov,
			p:        p,
			weight:   w,
			disabled: defaultDisabled,
			activate: f,
		}
	}

	tests := []struct {
		name         string
		conf         *Configuration
		alphaGenes   []*gene
		alphaFitness float64
		betaGenes    []*gene
		betaFitness  float64

		expect []*gene
	}{
		{
			name: "No unshared genes",
			conf: &Configuration{
				Inputs:  1,
				Outputs: 1,
			},
			alphaGenes: []*gene{
				newGene(nodePair{1, 10}, 1, sigmoid, geneID(1)),
				newGene(nodePair{2, 10}, 1, sigmoid, geneID(2)),
			},
			alphaFitness: 1.0,
			betaGenes: []*gene{
				newGene(nodePair{1, 10}, 1, sigmoid, geneID(1)),
				newGene(nodePair{2, 10}, 1, sigmoid, geneID(2)),
			},
			betaFitness: 0.9,

			expect: []*gene{
				newGene(nodePair{1, 10}, 1, sigmoid, geneID(1)),
				newGene(nodePair{2, 10}, 1, sigmoid, geneID(2)),
			},
		},
		{
			name: "Disjoint genes in alpha",
			conf: &Configuration{
				Inputs:  1,
				Outputs: 1,
			},
			alphaGenes: []*gene{
				newGene(nodePair{1, 10}, 1, sigmoid, geneID(1)),
				newGene(nodePair{2, 10}, 1, sigmoid, geneID(2)),
				newGene(nodePair{3, 10}, 1, sigmoid, geneID(3)),
				newGene(nodePair{4, 10}, 1, sigmoid, geneID(4)),
				newGene(nodePair{5, 10}, 1, sigmoid, geneID(5)),
			},
			alphaFitness: 1.0,
			betaGenes: []*gene{
				newGene(nodePair{1, 10}, 1, sigmoid, geneID(1)),
				newGene(nodePair{2, 10}, 1, sigmoid, geneID(2)),
			},
			betaFitness: 0.9,

			expect: []*gene{
				newGene(nodePair{1, 10}, 1, sigmoid, geneID(1)),
				newGene(nodePair{2, 10}, 1, sigmoid, geneID(2)),
				newGene(nodePair{3, 10}, 1, sigmoid, geneID(3)),
				newGene(nodePair{4, 10}, 1, sigmoid, geneID(4)),
				newGene(nodePair{5, 10}, 1, sigmoid, geneID(5)),
			},
		},
		{
			name: "Disjoint genes in beta",
			conf: &Configuration{
				Inputs:  1,
				Outputs: 1,
			},
			alphaGenes: []*gene{
				newGene(nodePair{1, 10}, 1, sigmoid, geneID(1)),
				newGene(nodePair{2, 10}, 1, sigmoid, geneID(2)),
			},
			alphaFitness: 1.0,
			betaGenes: []*gene{
				newGene(nodePair{1, 10}, 1, sigmoid, geneID(1)),
				newGene(nodePair{2, 10}, 1, sigmoid, geneID(2)),
				newGene(nodePair{3, 10}, 1, sigmoid, geneID(3)),
				newGene(nodePair{4, 10}, 1, sigmoid, geneID(4)),
				newGene(nodePair{5, 10}, 1, sigmoid, geneID(5)),
			},
			betaFitness: 0.9,

			expect: []*gene{
				newGene(nodePair{1, 10}, 1, sigmoid, geneID(1)),
				newGene(nodePair{2, 10}, 1, sigmoid, geneID(2)),
				newGene(nodePair{3, 10}, 1, sigmoid, geneID(3)),
				newGene(nodePair{4, 10}, 1, sigmoid, geneID(4)),
				newGene(nodePair{5, 10}, 1, sigmoid, geneID(5)),
			},
		},
		{
			name: "Disjoint genes in alpha and beta",
			conf: &Configuration{
				Inputs:  1,
				Outputs: 1,
			},
			alphaGenes: []*gene{
				newGene(nodePair{1, 10}, 1, sigmoid, geneID(1)),
				newGene(nodePair{2, 10}, 1, sigmoid, geneID(2)),
				newGene(nodePair{4, 10}, 1, sigmoid, geneID(4)),
			},
			alphaFitness: 1.0,
			betaGenes: []*gene{
				newGene(nodePair{1, 10}, 1, sigmoid, geneID(1)),
				newGene(nodePair{3, 10}, 1, sigmoid, geneID(3)),
				newGene(nodePair{5, 10}, 1, sigmoid, geneID(5)),
			},
			betaFitness: 0.9,

			expect: []*gene{
				newGene(nodePair{1, 10}, 1, sigmoid, geneID(1)),
				newGene(nodePair{2, 10}, 1, sigmoid, geneID(2)),
				newGene(nodePair{3, 10}, 1, sigmoid, geneID(3)),
				newGene(nodePair{4, 10}, 1, sigmoid, geneID(4)),
				newGene(nodePair{5, 10}, 1, sigmoid, geneID(5)),
			},
		},
		{
			name: "Gene gap in beta",
			conf: &Configuration{
				Inputs:  1,
				Outputs: 1,
			},
			alphaGenes: []*gene{
				newGene(nodePair{1, 10}, 1, sigmoid, geneID(1)),
				newGene(nodePair{2, 10}, 1, sigmoid, geneID(2)),
				newGene(nodePair{3, 10}, 1, sigmoid, geneID(3)),
				newGene(nodePair{4, 10}, 1, sigmoid, geneID(4)),
				newGene(nodePair{5, 10}, 1, sigmoid, geneID(5)),
			},
			alphaFitness: 1.0,
			betaGenes: []*gene{
				newGene(nodePair{1, 10}, 1, sigmoid, geneID(1)),
				newGene(nodePair{3, 10}, 1, sigmoid, geneID(3)),
				newGene(nodePair{5, 10}, 1, sigmoid, geneID(5)),
			},
			betaFitness: 0.9,

			expect: []*gene{
				newGene(nodePair{1, 10}, 1, sigmoid, geneID(1)),
				newGene(nodePair{2, 10}, 1, sigmoid, geneID(2)),
				newGene(nodePair{3, 10}, 1, sigmoid, geneID(3)),
				newGene(nodePair{4, 10}, 1, sigmoid, geneID(4)),
				newGene(nodePair{5, 10}, 1, sigmoid, geneID(5)),
			},
		},
		{
			name: "Gene gap in alpha",
			conf: &Configuration{
				Inputs:  1,
				Outputs: 1,
			},
			alphaGenes: []*gene{
				newGene(nodePair{1, 10}, 1, sigmoid, geneID(1)),
				newGene(nodePair{3, 10}, 1, sigmoid, geneID(3)),
				newGene(nodePair{5, 10}, 1, sigmoid, geneID(5)),
			},
			alphaFitness: 1.0,
			betaGenes: []*gene{
				newGene(nodePair{1, 10}, 1, sigmoid, geneID(1)),
				newGene(nodePair{2, 10}, 1, sigmoid, geneID(2)),
				newGene(nodePair{3, 10}, 1, sigmoid, geneID(3)),
				newGene(nodePair{4, 10}, 1, sigmoid, geneID(4)),
				newGene(nodePair{5, 10}, 1, sigmoid, geneID(5)),
			},
			betaFitness: 0.9,

			expect: []*gene{
				newGene(nodePair{1, 10}, 1, sigmoid, geneID(1)),
				newGene(nodePair{2, 10}, 1, sigmoid, geneID(2)),
				newGene(nodePair{3, 10}, 1, sigmoid, geneID(3)),
				newGene(nodePair{4, 10}, 1, sigmoid, geneID(4)),
				newGene(nodePair{5, 10}, 1, sigmoid, geneID(5)),
			},
		},
		{
			name: "Higher fitness in alpha",
			conf: &Configuration{
				Inputs:  1,
				Outputs: 1,
			},
			alphaGenes: []*gene{
				newGene(nodePair{1, 10}, 2, sigmoid, geneID(1)),
				newGene(nodePair{2, 10}, 2, sigmoid, geneID(2)),
				newGene(nodePair{3, 10}, 2, sigmoid, geneID(3)),
				newGene(nodePair{4, 10}, 2, sigmoid, geneID(4)),
				newGene(nodePair{5, 10}, 2, sigmoid, geneID(5)),
			},
			alphaFitness: 1.0,
			betaGenes: []*gene{
				newGene(nodePair{1, 10}, 1, sigmoid, geneID(1)),
				newGene(nodePair{2, 10}, 1, sigmoid, geneID(2)),
				newGene(nodePair{3, 10}, 1, sigmoid, geneID(3)),
				newGene(nodePair{4, 10}, 1, sigmoid, geneID(4)),
				newGene(nodePair{5, 10}, 1, sigmoid, geneID(5)),
			},
			betaFitness: 0.9,

			expect: []*gene{
				newGene(nodePair{1, 10}, 2, sigmoid, geneID(1)),
				newGene(nodePair{2, 10}, 2, sigmoid, geneID(2)),
				newGene(nodePair{3, 10}, 2, sigmoid, geneID(3)),
				newGene(nodePair{4, 10}, 2, sigmoid, geneID(4)),
				newGene(nodePair{5, 10}, 2, sigmoid, geneID(5)),
			},
		},
		{
			name: "Higher fitness in beta",
			conf: &Configuration{
				Inputs:  1,
				Outputs: 1,
			},
			alphaGenes: []*gene{
				newGene(nodePair{1, 10}, 1, sigmoid, geneID(1)),
				newGene(nodePair{2, 10}, 1, sigmoid, geneID(2)),
				newGene(nodePair{3, 10}, 1, sigmoid, geneID(3)),
				newGene(nodePair{4, 10}, 1, sigmoid, geneID(4)),
				newGene(nodePair{5, 10}, 1, sigmoid, geneID(5)),
			},
			alphaFitness: 0.9,
			betaGenes: []*gene{
				newGene(nodePair{1, 10}, 2, sigmoid, geneID(1)),
				newGene(nodePair{2, 10}, 2, sigmoid, geneID(2)),
				newGene(nodePair{3, 10}, 2, sigmoid, geneID(3)),
				newGene(nodePair{4, 10}, 2, sigmoid, geneID(4)),
				newGene(nodePair{5, 10}, 2, sigmoid, geneID(5)),
			},
			betaFitness: 1.0,

			expect: []*gene{
				newGene(nodePair{1, 10}, 2, sigmoid, geneID(1)),
				newGene(nodePair{2, 10}, 2, sigmoid, geneID(2)),
				newGene(nodePair{3, 10}, 2, sigmoid, geneID(3)),
				newGene(nodePair{4, 10}, 2, sigmoid, geneID(4)),
				newGene(nodePair{5, 10}, 2, sigmoid, geneID(5)),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := newCleanSpecies(test.conf)
			a := newCleanOrganism(test.conf)
			b := newCleanOrganism(test.conf)

			for _, g := range test.alphaGenes {
				a.nodes[g.p.input] = 0
				a.nodes[g.p.output] = 0
				a.add(g)
			}
			a.fitness = test.alphaFitness

			for _, g := range test.betaGenes {
				b.nodes[g.p.input] = 0
				b.nodes[g.p.output] = 0
				b.add(g)
			}
			b.fitness = test.betaFitness

			c := s.recombinate(a, b)

			require.Equal(t, len(test.expect), len(c.oinnov))
			for i, x := range test.expect {
				y := c.oinnov[i]
				require.True(t, x.equalTo(y))
			}
		})
	}
}
