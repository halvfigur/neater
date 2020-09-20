package neat

import (
	"fmt"
	"time"
)

type (
	Stats struct {
		Iterations int
		NbrSpecies int

		Champion *species
	}

	Neat struct {
		conf  *Configuration
		stats Stats
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
		conf: c,
	}

	return n, nil
}

func (n *Neat) Train(tf TrainerFactory, cf FitnessCalculatorFactory) {

	// We need some condition that lets us abort the traning, we could for
	// instance look at the fitness values and see if they have converged.
	// For now just iterate into oblivion.
	condition := true

	ss := make([]*species, 0, n.conf.MaxPopulationSize)
	ss = append(ss, newSpecies(n.conf))

	rejects := make([]*organism, 0, n.conf.MaxPopulationSize)

	for {
		n.stats.Iterations++

		for _, s := range ss {
			s.train(tf, cf)

			if n.stats.Champion == nil ||
				s.champ.fitness > n.stats.Champion.champ.fitness {
				n.stats.Champion = s
			}
		}

		if !condition {
			break
		}

		// Mutate & mate organisms
		for _, s := range ss {
			rejected := s.mutate()
			if rejected != nil {
				rejects = append(rejects, rejected...)
			}

			s.mate()
		}

		// Handle organisms that where rejected by their species
		for _, o := range rejects {
			for _, s := range ss {
				if s.belongs(o) {
					s.add(o)
					break
				}
			}

			// Couldn't find a suitable species for organism, time to create a new species
			s := newSpecies(n.conf)
			s.add(o)
			ss = append(ss, s)
		}

		n.stats.NbrSpecies = len(ss)
		rejects = rejects[:0]

		n.printStats()
	}
}

func (n *Neat) printStats() {
	//fmt.Print("\033[2J")
	fmt.Printf("---General-------\n")
	fmt.Printf("Iterations:        %-3d\n", n.stats.Iterations)
	fmt.Printf("NbrSpecies:        %-3d\n", n.stats.NbrSpecies)

	fmt.Printf("---Top Species---\n")
	fmt.Printf("Generation:        %-3d\n", n.stats.Champion.generation)
	fmt.Printf("Population size:   %-3d\n", len(n.stats.Champion.population))
	fmt.Printf("Fitness:        %.2f\n", n.stats.Champion.champ.fitness)
	time.Sleep(time.Second)
	fmt.Printf("\n\n")
}
