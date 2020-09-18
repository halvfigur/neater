package neat

import (
	"math"
)

type (
	species struct {
		conf       *Configuration
		rep        *organism
		champ      *organism
		population []*organism
	}
)

func newSpecies(c *Configuration) *species {
	s := &species{
		conf:       c,
		population: make([]*organism, c.InitialPopulationSize),
	}

	o := newOrganism(c)
	s.population[0] = o

	for i := range s.population[1:] {
		s.population[i+i] = o.copy()
	}

	return s
}

func (s *species) mutate() []*organism {
	cache := make(map[nodePair]*gene)
	//TODO set size of rejects based on configuration value
	rejectIdx := make([]int, 0, 64)

	for i, o := range s.population {
		o.mutate(cache)

		if !s.belongs(o) {
			rejectIdx = append(rejectIdx, i)
		}
	}

	population := make([]*organism, 0, len(s.population)-len(rejectIdx))
	rejects := make([]*organism, 0, len(rejectIdx))

	for i, o := range s.population {
		if i == rejectIdx[0] {
			rejects = append(rejects, o)
			rejectIdx = rejectIdx[1:]
		} else {
			population = append(population, o)
		}
	}

	s.population = population

	return rejects
}

func (s *species) train(tf TrainerFactory, cf FitnessCalculatorFactory) {
	s.champ = s.population[0]

	for _, o := range s.population {
		t := tf.New()
		c := cf.New()
		for input, ok := t.Next(); ok; input, ok = t.Next() {
			output := o.Eval(input)
			c.AddResult(input, output)
		}

		o.fitness = c.CalculateFitness()

		if o.fitness > s.champ.fitness {
			s.champ = o
		}
	}

	// Normalize the species fitness
	s.normalize()

	// Chose a species representative
	r := randIntn(len(s.population))
	s.rep = s.population[r]
}

// normalize normalizes the fitness of the population
func (s *species) normalize() {
	l := float64(len(s.population))

	for _, o := range s.population {
		o.fitness /= l
	}
}

func (s *species) belongs(o *organism) bool {
	return s.distance(s.rep, o) < s.conf.CompatibilityThreshold
}

func (s *species) add(o *organism) {
	s.population = append(s.population, o)
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

func (s *species) mate(a, b *organism) *organism {

	// Switch if necessary so that `a` has the best performance
	if a.fitness < b.fitness {
		a, b = b, a
	}

	o := newCleanOrganism(a.conf)
	copy(o.inputs, a.inputs)
	copy(o.outputs, a.outputs)

	i, j := 0, 0

	// Copy genes and hidden nodes
	for i < len(a.oinnov) && j < len(b.oinnov) {
		var g gene

		if a.oinnov[i].innov == b.oinnov[j].innov {
			// ´a´ has the better performance so copy the gene from from `a`
			g = *a.oinnov[i]
			i = min(i+1, len(a.oinnov))
			j = min(j+1, len(b.oinnov))
		} else if a.oinnov[i].innov < b.oinnov[j].innov {
			// `a` has a gene not present in ´b´
			g = *a.oinnov[i]
			i = min(i+1, len(a.oinnov))
		} else {
			// `b` has a gene not present in ´a´
			g = *b.oinnov[i]
			j = min(j+1, len(b.oinnov))
		}

		// Create the nodes in the target organism if the don't already exist.
		// TODO: figure out if we need a function for creating nodes.
		o.nodes[g.p.input] = 0
		o.nodes[g.p.output] = 0
		o.add(&g)
	}

	// Handle trailing genes (if any)
	for ; i < len(a.oinnov); i++ {
		g := *a.oinnov[i]
		o.add(&g)
	}

	// Handle trailing genes (if any)
	for ; j < len(b.oinnov); j++ {
		g := *b.oinnov[i]
		o.add(&g)
	}

	return o
}
