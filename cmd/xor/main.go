package main

import (
	"log"
	"math/rand"
	"neater"
	"os"
	"os/signal"
)

func main() {

	rand.Seed(0)
	tf := NewXORTrainerFactory()
	cf := NewXORFitnessCalculatorFactory()

	c := &neater.Configuration{
		// Inputs is the number of inputs
		Inputs: tf.Inputs(),

		// Outputs is the number of outputs
		Outputs: tf.Outputs(),

		// WeightMutationProb is the probability that a given gene's weight is mutated
		WeightMutationProb: 0.1,

		// WeightMutationPower is the threshold for wait mutations in one mutation
		WeightMutationPower: 2.5,

		// WeightMutationStandardDeviation
		WeightMutationStandardDeviation: 0.5,

		// AddNodeMutationProb is the probability that a gene is disabled and a new Node is inserted
		AddNodeMutationProb: 0.1,

		// ConnectNodesMutationProb is the probability that a new gene connecting two nodes hkk
		ConnectNodesMutationProb: 0.1,

		// PopulationThreshold is the maximum size of a species population
		PopulationThreshold: 32,

		// Recurrent controls whether recurrent connections are allowed
		Recurrent: false,

		// MaxPopulationSize is the maximum number of different species
		MaxPopulationSize: 64,

		// DisjointCoefficient
		DisjointCoefficient: 2.0,

		// ExcessCoefficient
		ExcessCoefficient: 2.0,

		// WeightDifferenceCoefficient
		WeightDifferenceCoefficient: 1.0,

		// CompatibilityThreshold
		CompatibilityThreshold: 6.0,

		// CompatibilityModifier
		CompatibilityModifier: 0.01,

		// DropOffAge
		DropOffAge: 15,

		// SurvivalThreshold controls how many percent of the population top
		// performers survive and reproduce, range (0, 1]
		SurvivalThreshold: 0.25,

		// MutationPower
		MutationPower: 2.5,

		// InitialPopulationSize
		InitialPopulationSize: 8,

		// ActivationFunction
		ActivationFunction: neater.ActivateSigmoid,
	}

	n, err := neater.NewNeat(c)
	if err != nil {
		log.Fatal(err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)

	for {
		select {
		case <-sigc:
			if n.Champion() != nil {
				return
			}
		default:
			n.Train(tf, cf)
		}
	}
}
