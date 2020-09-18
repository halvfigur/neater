package main

import (
	"log"
	"neat"
)

func main() {

	tf := neat.NewXORTrainerFactory()
	cf := neat.NewXORFitnessCalculatorFactory()

	c := &neat.Configuration{
		// Inputs is the number of inputs
		Inputs: tf.Inputs(),

		// Outputs is the number of outputs
		Outputs: tf.Outputs(),

		// WeightMutationProb is the probability that a given gene's weight is mutated
		WeightMutationProb: 0.1,

		WeightMutationPower: 2.5,

		// AddNodeMutationProb is the probability that a gene is disabled and a new Node is inserted
		AddNodeMutationProb: 0.1,

		// ConnectNodesMutationProb is the probability that a new gene connecting two nodes hkk
		ConnectNodesMutationProb: 0.1,

		// PopulationThreshold is the maximum size of a species population
		PopulationThreshold: 32,

		// Recurrent controls whether recurrent connections are allowed
		Recurrent: false,

		// MaxPopulationSize is the maximum number of different species
		MaxPopulationSize: 32,

		// DisjointCoefficient
		DisjointCoefficient: 2.0,

		// ExcessCoefficient
		ExcessCoefficient: 2.0,

		// WeightDifferenceCoefficient
		WeightDifferenceCoefficient: 1.0,

		// CompatibilityThreshold
		CompatibilityThreshold: 6.0,

		// CompatibilityModifier
		CompatibilityModifier: 0.3,

		// DropOffAge
		DropOffAge: 15,

		// SurvivalThreshold controls how many percent of the population top
		// performers survive and reproduce, range (0, 1]
		SurvivalThreshold: 0.2,

		// MutationPower
		MutationPower: 2.5,

		// InitialPopulationSize
		InitialPopulationSize: 10,
	}

	n, err := neat.NewNeat(c)
	if err != nil {
		log.Fatal(err)
	}

	n.Train(tf, cf)
}
