package neat

type (
	XORTrainer struct {
		input [][]float64
		n     int
	}

	XORFitnessCalculator struct {
		aggrErr float64
	}
)

func NewXORTrainer() *XORTrainer {
	return &XORTrainer{
		input: [][]float64{
			[]float64{0, 0},
			[]float64{0, 1},
			[]float64{1, 0},
			[]float64{1, 1},
		},
	}
}

func (t *XORTrainer) Next() ([]float64, bool) {
	if t.n == 4 {
		return nil, false
	}

	input := t.input[t.n]
	t.n++

	return input, true
}

func (t *XORTrainer) Reset() {
	t.n = 0
}

func NewXORFitnessCalculator() *XORFitnessCalculator {
	return new(XORFitnessCalculator)
}

func (c *XORFitnessCalculator) AddResult(input, output []float64) {
	if len(input) != 2 {
		panic("invalid input")
	}

	if len(output) != 1 {
		panic("invalid output")
	}

	i0, i1 := input[0], input[1]
	o := output[0]
	v0, v1 := float64(0), float64(1)

	if i0 == v0 && i1 == v0 {
		c.aggrErr += (v0 - o) * (v0 - o)
		return
	}

	if i0 == v0 && i1 == v1 {
		c.aggrErr += (v1 - o) * (v1 - o)
		return
	}

	if i0 == v1 && i1 == v0 {
		c.aggrErr += (v1 - o) * (v1 - o)
		return
	}

	if i0 == v1 && i1 == v1 {
		c.aggrErr += (v0 - o) * (v0 - o)
		return
	}

	panic("invalid input")
}

func (c *XORFitnessCalculator) CalculateFitness() float64 {
	return float64(1) - c.aggrErr
}

func (c *XORFitnessCalculator) Reset() {
	c.aggrErr = float64(0)
}

func NewXORTrainerFactory() TrainerFactory {
	return TrainerFactory{
		New: func() Trainer {
			return NewXORTrainer()
		},
		Inputs: func() int {
			return 2
		},
		Outputs: func() int {
			return 1
		},
	}
}

func NewXORFitnessCalculatorFactory() FitnessCalculatorFactory {
	return FitnessCalculatorFactory{
		New: func() FitnessCalculator {
			return NewXORFitnessCalculator()
		},
	}
}
