package neater

import (
	"fmt"
	"sort"
)

type (
	Stats struct {
		Iterations int
		NbrSpecies int

		BestSpecies  *species
		BestOrganism *organism
	}

	Neat struct {
		conf    *Configuration
		species []*species
		inputs  []nodeID
		outputs []nodeID
		stats   Stats
	}
)

func NewNeat(c *Configuration) (*Neat, error) {
	switch c.ActivationFunction {
	case ActivateSigmoid:
		c.activate = sigmoid
	case ActivateUnit:
		c.activate = unit
	default:
		panic("unknown activation function")
	}

	n := &Neat{
		conf:    c,
		species: make([]*species, 0, c.MaxPopulationSize),
	}

	n.inputs = make([]nodeID, n.conf.Inputs)
	for i := range n.inputs {
		n.inputs[i] = nodeIDGenerator()
	}
	n.outputs = make([]nodeID, n.conf.Outputs)
	for i := range n.outputs {
		n.outputs[i] = nodeIDGenerator()
	}

	n.species = append(n.species, newSpecies(n.conf, n.inputs, n.outputs))

	return n, nil
}

func (n *Neat) BestOrganism() *organism {
	return n.stats.BestOrganism
}

func (n *Neat) train(tf TrainerFactory, cf FitnessCalculatorFactory) {
	for _, s := range n.species {
		s.train(tf, cf)
	}

	// Sort the species according to their most fit organism
	sort.Slice(n.species, func(i, j int) bool {
		return n.species[i].champ.fitness > n.species[j].champ.fitness
	})

	n.stats.BestSpecies = n.species[0]
	n.stats.BestOrganism = n.stats.BestSpecies.champ
}

func (n *Neat) adjustPopulationSize() {
	// Adjust the population according to the SurvivalThreshold
	if len(n.species) > n.conf.MaxPopulationSize {
		n.species = n.species[:n.conf.MaxPopulationSize]
	}
}

func (n *Neat) mutate() []*organism {
	rejected := make([]*organism, 0, n.conf.MaxPopulationSize)

	// Mutate & mate organisms
	for _, s := range n.species {
		rejects := s.mutate()
		if rejects != nil {
			rejected = append(rejected, rejects...)
		}
	}

	return rejected

}

func (n *Neat) handleRejected(rejected []*organism) {
	for _, o := range rejected {
		inserted := false
		for _, s := range n.species {
			if s.belongs(o) {
				s.add(o)
				inserted = true
				break
			}
		}

		if !inserted {
			// Couldn't find a suitable species for organism, time to create a new species

			s := newSpecies(n.conf, n.inputs, n.outputs)
			s.add(o)
			n.species = append(n.species, s)
		}
	}
}

func (n *Neat) mate() {
	for _, s := range n.species {
		s.mate()
	}
}

func (n *Neat) Train(tf TrainerFactory, cf FitnessCalculatorFactory) float64 {

	n.stats.Iterations++

	n.train(tf, cf)

	n.adjustPopulationSize()

	n.stats.NbrSpecies = len(n.species)

	n.printStats()

	rejected := n.mutate()

	n.handleRejected(rejected)

	n.mate()

	return n.stats.BestOrganism.fitness
}

func (n *Neat) printStats() {
	fmt.Print("\033[2J")
	fmt.Printf("---General--------\n")
	fmt.Printf("Iterations:      %10d\n", n.stats.Iterations)
	fmt.Printf("NbrSpecies:      %10d\n", n.stats.NbrSpecies)

	fmt.Printf("---Top Species----\n")
	fmt.Printf("ID:              %10d\n", n.stats.BestSpecies.id)
	fmt.Printf("Generation:      %10d\n", n.stats.BestSpecies.generation)
	fmt.Printf("Population size: %10d\n", len(n.stats.BestSpecies.population))

	fmt.Printf("---Top Organism---\n")
	fmt.Printf("ID:              %10d\n", n.stats.BestOrganism.id)
	fmt.Printf("Fitness:         %10f\n", n.stats.BestOrganism.fitness)
	fmt.Printf("Node count:      %10d\n", len(n.stats.BestOrganism.nodes))
	fmt.Printf("Gene count:      %10d\n", len(n.stats.BestOrganism.oinnov))
	fmt.Printf("\n\n")
}
