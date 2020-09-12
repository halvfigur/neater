package main

type (
	Configuration struct {
		// Inputs is the number of inputs
		Inputs int

		// Outputs is the number of outputs
		Outputs int

		// WeightMutationProb is the probability that a given gene's weight is mutated
		WeightMutationProb int

		WeightMutationPower float64

		// AddNodeMutationProb is the probability that a gene is disabled and a new Node is inserted
		AddNodeMutationProb int

		// ConnectNodesMutationProb is the probability that a new gene connecting two nodes hkk
		ConnectNodesMutationProb int

		// PopulationThreshold is the maximum size of a species population
		PopulationThreshold int

		// Recurrent controlls whether recurrent connections are allowed
		Recurrent bool
	}
)
