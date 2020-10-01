package neater

type (
	FitnessCalculator interface {
		// AddResult adds a new input/output result
		AddResult(input, output []float64)

		// CalculateFitness calculates fitness
		CalculateFitness() float64

		// Reset resets the FitnessCalculator
		Reset()
	}

	Trainer interface {
		// Next returns the next input or nil if the set is exhausted
		Next() ([]float64, bool)

		// Reset resets the Trainer
		Reset()
	}

	TrainerFactory struct {
		// Inputs is the number of inputs
		Inputs func() int

		// Outputs is the number of outputs
		Outputs func() int

		// New creates a new Trainer
		New func() Trainer
	}

	FitnessCalculatorFactory struct {
		// New creates a new FitnessCalculatorFactory
		New func() FitnessCalculator
	}
)
