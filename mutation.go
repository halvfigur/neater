package main

import "math/rand"

type (
	species struct {
		conf       *Configuration
		population []*organism
	}

	nodePair struct {
		input  nodeID
		output nodeID
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

func (s *species) mutate() {
	innovCache := make(map[nodePair]*gene)

	connectPair := func(o *organism, p nodePair) {
		if x, ok := innovCache[p]; ok {
			// This innovation has already been made
			o.add(x.copy())
		} else {
			x := newGene(p.input, p.output)
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

				connectPair(o, nodePair{g.input, id})
				connectPair(o, nodePair{id, g.output})
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
					if g.input == p.output {
						firstIdx = i
					}
				}

				if g.output == p.input {
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
