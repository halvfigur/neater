package neat

type (
	Neat struct {
		conf *Configuration
	}
)

func NewNeat(c *Configuration) (*Neat, error) {
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

	sl := make([]*species, 0, n.conf.MaxPopulationSize)
	sl = append(sl, newSpecies(n.conf))

	rejects := make([]*organism, n.conf.MaxPopulationSize)
	for {
		for _, s := range sl {
			s.train(tf, cf)
		}

		if !condition {
			break
		}

		// Mutate organisms
		for _, s := range sl {
			rejected := s.mutate()
			rejects = append(rejects, rejected...)
		}

		// Handle organisms that where rejected by their species
		for _, o := range rejects {
			for _, s := range sl {
				if s.belongs(o) {
					s.add(o)
					break
				}
			}

			// Couldn't find a suitable species for organism, time to create a new species
			s := newSpecies(n.conf)
			s.add(o)
			sl = append(sl, s)
		}

		rejects = rejects[:0]
	}
}
