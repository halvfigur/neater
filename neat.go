package neater

import (
	"fmt"
	"sort"
)

type (
	Stats struct {
		Iterations int
		NbrSpecies int

		Champion *species
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

func (n *Neat) Champion() *organism {
	return n.stats.Champion.champ
}

func (n *Neat) Train(tf TrainerFactory, cf FitnessCalculatorFactory) float64 {

	rejects := make([]*organism, 0, n.conf.MaxPopulationSize)

	n.stats.Iterations++

	for _, s := range n.species {
		s.train(tf, cf)
	}

	// Adjust the population according to the SurvivalThreshold
	sort.Slice(n.species, func(i, j int) bool {
		return n.species[i].champ.fitness > n.species[j].champ.fitness
	})

	n.stats.Champion = n.species[0]

	if len(n.species) > n.conf.MaxPopulationSize {
		n.species = n.species[:n.conf.MaxPopulationSize]
	}

	// Mutate & mate organisms
	for _, s := range n.species {
		rejected := s.mutate()
		if rejected != nil {
			rejects = append(rejects, rejected...)
		}

		s.mate()
	}

	// Handle organisms that where rejected by their species
	for _, o := range rejects {
		for _, s := range n.species {
			if s.belongs(o) {
				s.add(o)
				break
			}
		}

		// Couldn't find a suitable species for organism, time to create a new species
		s := newSpecies(n.conf, n.inputs, n.outputs)
		s.add(o)
		n.species = append(n.species, s)
	}

	n.stats.NbrSpecies = len(n.species)
	rejects = rejects[:0]

	n.printStats()

	return n.Champion().fitness
}

func (n *Neat) printStats() {
	fmt.Print("\033[2J")
	fmt.Printf("---General-------\n")
	fmt.Printf("Iterations:        %-3d\n", n.stats.Iterations)
	fmt.Printf("NbrSpecies:        %-3d\n", n.stats.NbrSpecies)

	fmt.Printf("---Top Species---\n")
	fmt.Printf("Generation:        %-3d\n", n.stats.Champion.generation)
	fmt.Printf("Population size:   %-3d\n", len(n.stats.Champion.population))
	//fmt.Printf("Fitness:        %.2f\n", n.stats.Champion.champ.fitness)
	fmt.Printf("Fitness:        %f\n", n.stats.Champion.champ.fitness)
	fmt.Printf("Node count:        %-3d\n", len(n.stats.Champion.champ.nodes))
	fmt.Printf("Gene count:        %-3d\n", len(n.stats.Champion.champ.oinnov))
	fmt.Printf("\n\n")
}
