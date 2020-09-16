package neat

import (
	"math"
	"math/rand"
)

type (
	species struct {
		conf       *Configuration
		rep        *organism
		population []*organism
	}
)

func newSpecies(conf *Configuration, o *organism) *species {
	s := &species{
		conf:       conf,
		population: make([]*organism, conf.PopulationThreshold),
	}

	s.population[0] = o

	for i := range s.population[1:] {
		s.population[i+i] = o.copy()
	}

	s.mutate()

	return s
}

// normalize normalizes the fitness of the population
func (s *species) normalize() {
	l := float64(len(s.population))

	for _, o := range s.population {
		o.fitness /= l
	}
}

func (s *species) mutate() {
	innovCache := make(map[nodePair]*gene)

	connectPair := func(o *organism, p nodePair) {
		if x, ok := innovCache[p]; ok {
			// This innovation has already been made
			o.add(x.copy())
		} else {
			x := newGene(p)
			innovCache[p] = x
			o.add(x)
		}
	}

	for _, o := range s.population {
		for _, g := range o.oinnov {
			if rand.Intn(o.conf.WeightMutationProb) == 0 {
				g.weight *= rand.Float64() * s.conf.WeightMutationPower
			}

			if rand.Intn(o.conf.AddNodeMutationProb) == 0 {
				p := s.getRandUnconnectedNodePair(o)

				connectPair(o, p)
			}

			if rand.Intn(o.conf.ConnectNodesMutationProb) == 0 {
				g.disabled = true

				id := nodeIDGenerator()
				o.nodes[id] = 0

				connectPair(o, nodePair{g.p.input, id})
				connectPair(o, nodePair{id, g.p.output})
			}
		}
	}
}

func (s *species) getRandUnconnectedNodePair(o *organism) nodePair {

	var p nodePair

	// Assume the nodes are already connected and keep going until we find
	// a pair that aren't connected. This may get us stuck in an infinite loop.
	alreadyConnected := true

	for alreadyConnected {
		p.input = o.randomNode()
		p.output = o.randomNode()

		// Make sure the input and output are different
		if p.input == p.output {
			continue
		}

		if !s.conf.Recurrent {
			// If reccurent connections aren't allowed then the first gene that
			// takes ´p.output´ as input must appear after the last gene that
			// outputs to 'p.input' in the evaluation order
			firstIdx := -1
			lastIdx := -1
			for i, g := range o.oeval {
				if firstIdx == -1 {
					if g.p.input == p.output {
						firstIdx = i
					}
				}

				if g.p.output == p.input {
					lastIdx = i
				}
			}

			// If there is a patch from ´p.input' to ´p.output' make sure that
			// firstIdx > lastIdx
			if firstIdx > lastIdx {
				continue
			}
		}

		// Make sure the nodes aren't already connected
		alreadyConnected = o.connected(p.input, p.output)
	}

	return p
}

func (s *species) belongs(o *organism) bool {
	return s.distance(s.rep, o) < s.conf.CompatibilityThreshold
}

func (s *species) distance(a, b *organism) float64 {
	var (
		commonGenes   int
		disjointGenes int
		excessGenes   int
		weightDiff    float64
	)

	i, j := 0, 0
	for i < len(a.oinnov) && j < len(b.oinnov) {
		if a.oinnov[i].innov == b.oinnov[j].innov {
			// ´a´ and ´b´ have a gene in common
			weightDiff += math.Abs(a.oinnov[i].weight - b.oinnov[i].weight)
			commonGenes++
			i = min(i+1, len(a.oinnov))
			j = min(j+1, len(b.oinnov))
		} else if a.oinnov[i].innov < b.oinnov[j].innov {
			// `a` has a gene not present in ´b´
			i = min(i+1, len(a.oinnov))
			disjointGenes++
		} else {
			// `b` has a gene not present in ´a´
			j = min(j+1, len(b.oinnov))
			disjointGenes++
		}
	}

	// Account for excess genes in ´a´ (if any)
	excessGenes += len(a.oinnov) - i - 1

	// Account for excess genes in ´b´ (if any)
	excessGenes += len(b.oinnov) - j - 1

	n := float64(1)
	largest := float64(max(len(a.oinnov), len(b.oinnov)))
	if largest > float64(20) {
		// 'n' normalizes for genome size ('n' can be set to 1
		// if both genomes are small, i.e., consist of fewer than 20 genes)
		n = largest
	}

	// Shorten the names so that the calculation is readable
	c1 := s.conf.ExcessCoefficient
	c2 := s.conf.DisjointCoefficient
	c3 := s.conf.WeightDifferenceCoefficient
	e := float64(excessGenes)
	d := float64(disjointGenes)
	w := weightDiff / float64(commonGenes)

	return ((c1*e)+(c2*d))/n + c3*w
}
