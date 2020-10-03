package neater

type (
	ActivationFunction string

	Configuration struct {
		// Inputs is the number of inputs
		Inputs int

		// Outputs is the number of outputs
		Outputs int

		// WeightMutationProb is the probability that a given gene's weight is mutated
		WeightMutationProb float64

		// WeightMutationPower is the threshold for weight mutations in one mutation
		WeightMutationPower float64

		// WeightMutationStandardDeviation
		WeightMutationStandardDeviation float64

		// AddNodeMutationProb is the probability that a gene is disabled and a new Node is inserted
		AddNodeMutationProb float64

		// ConnectNodesMutationProb is the probability that a new gene connecting two nodes hkk
		ConnectNodesMutationProb float64

		// PopulationThreshold is the maximum size of a species population
		PopulationThreshold int

		// Recurrent controls whether recurrent connections are allowed
		Recurrent bool

		// RecurrentConnProb the probability that a new connection is recurrent
		RecurrentConnProb float64

		// MaxPopulationSize is the maximum number of different species
		MaxPopulationSize int

		// DisjointCoefficient
		DisjointCoefficient float64

		// ExcessCoefficient
		ExcessCoefficient float64

		// WeightDifferenceCoefficient
		WeightDifferenceCoefficient float64

		// CompatibilityThreshold
		CompatibilityThreshold float64

		// CompatibilityModifier
		CompatibilityModifier float64

		// DropOffAge
		DropOffAge int

		// SurvivalThreshold controls how many percent of the population top
		// performers survive and reproduce, range (0, 1]
		SurvivalThreshold float64

		// MutationPower
		MutationPower float64

		// InitialPopulationSize
		InitialPopulationSize int

		// ActivationFunction
		ActivationFunction string

		// NormalizeFitness
		NormalizeFitness bool

		// FitnessNormalizationThreshold
		FitnessNormalizationThreshold int

		activate activationFunction
	}
)

const (
	ActivateSigmoid = "sigmoid"
	ActivateUnit    = "unit"
)
