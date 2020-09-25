package neat

import (
	"math"
	"sort"
)

type (
	species struct {
		conf       *Configuration
		rep        *organism
		champ      *organism
		population []*organism
		generation int
	}
)

func newCleanSpecies(c *Configuration) *species {
	return &species{
		conf:       c,
		population: make([]*organism, c.InitialPopulationSize),
	}
}

func newSpecies(c *Configuration) *species {
	s := newCleanSpecies(c)
	o := newOrganism(c)
	s.population[0] = o

	for i := 1; i < len(s.population); i++ {
		s.population[i] = o.copy()
	}

	s.choseRepresentative()

	return s
}

func (s *species) choseRepresentative() {
	// Chose a species representative
	r := randIntn(len(s.population))
	s.rep = s.population[r]
}

func (s *species) mutate() []*organism {
	s.generation++

	// A cache to hold new connection innovations that have already been made
	// in this generation.
	connCache := make(map[nodePair]*gene)

	// A cache to hold new node innovations that have already been made in this
	// generation.
	nodeCache := make(map[nodePair]genePair)

	// rejectIdx stores the indices of the organisms that are no longer
	// compatible with the species after mutation
	//TODO set size of rejects based on configuration value
	rejectIdx := make([]int, 0, len(s.population))

	// Iteratate over the population and mutate all organisms except the
	// champion. Any organism that is no longer compatible with the species
	// representative is marked as rejected
	for i, o := range s.population {
		// Spare the champ from mutation
		if o == s.champ {
			continue
		}

		o.mutate(connCache, nodeCache)

		// If o no longer belongs, mark it as rejected
		if !s.belongs(o) {
			rejectIdx = append(rejectIdx, i)
		}
	}

	// If no organism was rejected we are done
	if len(rejectIdx) == 0 {
		return nil
	}

	// Time to separate out the rejected organisms.

	// population stores all organisms that should remain in the species
	population := make([]*organism, 0, len(s.population)-len(rejectIdx))

	// rejects stores all organisms that should be removed from the species
	rejects := make([]*organism, 0, len(rejectIdx))

	for i, o := range s.population {
		if len(rejectIdx) > 0 && i == rejectIdx[0] {
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

	// Adjust the population according to the SurvivalThreshold
	sort.Slice(s.population, func(i, j int) bool {
		return s.population[i].fitness > s.population[j].fitness
	})

	// Normalize the species fitness
	s.normalize()

	// Chose a new species representative
	s.choseRepresentative()
}

// normalize normalizes the fitness of the population
func (s *species) normalize() {
	l := float64(len(s.population))

	for _, o := range s.population {
		o.fitness /= l
	}
}

func (s *species) belongs(o *organism) bool {
	return s.distance(s.rep, o) < s.conf.CompatibilityThreshold+(s.conf.CompatibilityModifier*float64(s.generation-1))
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
			weightDiff += math.Abs(a.oinnov[i].weight - b.oinnov[j].weight)
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
	if s.conf.NormalizeFitness {
		largest := float64(max(len(a.oinnov), len(b.oinnov)))
		if largest > float64(s.conf.FitnessNormalizationThreshold) {
			// 'n' normalizes for genome size ('n' can be set to 1
			// if both genomes are small, i.e., consist of fewer than 20 genes)
			n = largest
		}
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

// mate mates the top species of the population, the organisms must be sorted
// in order of ascending fitness before entering this function.
func (s *species) mate() {
	/*
		// Let only the top performers survive and mate
		topCutOffIndex := int(float64(s.conf.PopulationThreshold) * s.conf.SurvivalThreshold)
		s.population = s.population[:min(topCutOffIndex, len(s.population))]
	*/

	cutOffIdx := min(len(s.population), int(math.Sqrt(float64(s.conf.PopulationThreshold+4))))
	s.population = s.population[:cutOffIdx]

	n := len(s.population)
	children := make([]*organism, 0, n*(n+1)/2)

	for i, a := range s.population {
		for j, b := range s.population {
			if i != j {
				c := s.recombinate(a, b)
				children = append(children, c)
			}
		}
	}

	s.population = append(s.population, children...)

	survivalIdx := min(s.conf.PopulationThreshold, len(s.population))

	s.population = s.population[:survivalIdx]
}

func (s *species) recombinate(a, b *organism) *organism {

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
		var g *gene

		if a.oinnov[i].innov == b.oinnov[j].innov {
			// ´a´ has the better performance so copy the gene from from `a`
			g = a.oinnov[i]
			i = min(i+1, len(a.oinnov))
			j = min(j+1, len(b.oinnov))
		} else if a.oinnov[i].innov < b.oinnov[j].innov {
			// `a` has a gene not present in ´b´
			g = a.oinnov[i]
			i = min(i+1, len(a.oinnov))
		} else {
			// `b` has a gene not present in ´a´
			g = b.oinnov[j]
			j = min(j+1, len(b.oinnov))
		}

		// Create the nodes in the target organism if the don't already exist.
		// TODO: figure out if we need a function for creating nodes.
		o.nodes[g.p.input] = 0
		o.nodes[g.p.output] = 0
		o.add(g)
	}

	// Handle trailing genes (if any)
	for ; i < len(a.oinnov); i++ {
		g := a.oinnov[i]
		o.nodes[g.p.input] = 0
		o.nodes[g.p.output] = 0
		o.add(g)
	}

	// Handle trailing genes (if any)
	for ; j < len(b.oinnov); j++ {
		g := b.oinnov[j]
		o.nodes[g.p.input] = 0
		o.nodes[g.p.output] = 0
		o.add(g)
	}

	return o
}
