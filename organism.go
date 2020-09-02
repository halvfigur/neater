package main

import "fmt"

type (
	connectStrategy int

	organism struct {
		inputs  []geneID
		outputs []geneID

		connectStrategy connectStrategy
		genome          map[geneID]*gene
		score           float64

		activate activationFunction
	}

	organismOpt func(*organism)
)

const (
	connectNone = connectStrategy(iota)
	connectFlow
	connectRandom

	defaultConnectStrategy = connectNone
)

func newOrganism(nInputs, nOutputs int, opts ...organismOpt) *organism {
	if nInputs <= 0 {
		panic("Number of inputs must be greater than 0")
	}

	if nOutputs <= 0 {
		panic("Number of outputs must be greater than 0")
	}

	o := &organism{
		connectStrategy: defaultConnectStrategy,
	}

	o.genome = make(map[geneID]*gene, len(o.inputs)*len(o.outputs))

	for _, opt := range opts {
		opt(o)
	}

	o.inputs = make([]geneID, nInputs)
	for i := 0; i < nInputs; i++ {
		g := newGene(withWeight(float64(1.0)), withActivationFunction(unit))
		o.inputs[i] = g.innov
		o.genome[g.innov] = g
	}

	o.outputs = make([]geneID, nOutputs)
	for i := 0; i < nOutputs; i++ {
		g := newGene(withWeight(float64(1.0)),
			withActivationFunction(o.activate),
			withOutput(terminal))
		o.outputs[i] = g.innov
		o.genome[g.innov] = g
	}

	o.connectTerminals()

	return o
}

func (o *organism) connectTerminals() {
	switch o.connectStrategy {
	case connectNone:
		panic("Not implemented")
	case connectFlow:
		o.connectFlow()
	case connectRandom:
		panic("Not implemented")
	}
}

func (o *organism) connectFlow() {
	m := max(len(o.inputs), len(o.outputs))

	for i := 0; i < m; i++ {
		inputID := o.inputs[i%len(o.inputs)]
		outputID := o.outputs[i%len(o.outputs)]

		o.gene(inputID).output = outputID
	}
}

func withGlobalActivationFunction(f activationFunction) organismOpt {
	return func(o *organism) {
		o.activate = f
	}
}

func withConnectStrategy(s connectStrategy) organismOpt {
	return func(o *organism) {
		o.connectStrategy = s
	}
}

func (o *organism) gene(id geneID) *gene {
	g, ok := o.genome[id]
	if !ok {
		panic(fmt.Sprintf("Gene not found %d", id))
	}

	return g
}

func (o *organism) inputGenes() []*gene {
	genes := make([]*gene, 0, len(o.inputs))

	for _, id := range o.inputs {
		genes = append(genes, o.gene(id))
	}

	return genes
}

func (o *organism) outputGenes() []*gene {
	genes := make([]*gene, 0, len(o.inputs))

	for _, id := range o.outputs {
		genes = append(genes, o.gene(id))
	}

	return genes
}

func (o *organism) clearGenome() {
	for _, g := range o.genome {
		g.clear()
	}
}

func (o *organism) Eval(inputs []float64) []float64 {
	if len(inputs) != len(o.inputs) {
		panic("Length of input vector must equal number of input nodes")
	}

	o.clearGenome()

	q := newSliceQueue(1024)

	// Going into this loop the genome must first be cleared
	for i, g := range o.inputGenes() {
		g.add(inputs[i])
		q.put(g)
	}

	for q.len() > 0 {
		g := q.get()
		h := o.gene(g.output)
		h.add(g.val())

		if !h.terminal() {
			q.put(h)
		}
	}

	outputs := make([]float64, len(o.outputs))
	for i, g := range o.outputGenes() {
		outputs[i] = g.val()
	}

	return outputs
}
