package main

import "fmt"

type (
	organism struct {
		inputs  []geneID
		outputs []geneID

		genome map[geneID]*gene
	}

	organismOpt func(*organism)
)

func newOrganism(opts ...organismOpt) *organism {
	o := new(organism)

	for _, opt := range opts {
		opt(o)
	}

	if o.inputs == nil {
		panic("the number of inputs must be spefified")
	}

	if o.outputs == nil {
		panic("the number of ouputs must be specified")
	}

	o.genome = make(map[geneID]*gene, len(o.inputs)*len(o.outputs))

	nInputs := len(o.inputs)
	for i := range o.inputs {
		g := newGene(withWeight(float64(1.0)), withActivationFunction(unit))
		o.inputs[i] = g.innov
		o.genome[g.innov] = g
	}

	nOutputs := len(o.outputs)
	for i := range o.outputs {
		g := newGene(withWeight(float64(0.5)), withActivationFunction(sigmoid))
		o.outputs[i] = terminal
		o.genome[g.innov] = g
	}

	m := max(nInputs, nOutputs)

	for i := 0; i < m; i++ {
		inputID := o.inputs[i%nInputs]
		outputID := o.outputs[i%nOutputs]

		o.genome[inputID].output = outputID
	}

	return o
}

func withNbrInputs(n int) organismOpt {
	return func(o *organism) {
		o.inputs = make([]geneID, n)
	}
}

func withNbrOutputs(n int) organismOpt {
	return func(o *organism) {
		o.outputs = make([]geneID, n)
	}
}

func (o *organism) getGene(id geneID) *gene {
	g, ok := o.genome[id]
	if !ok {
		panic(fmt.Sprintf("Gene not found %d", id))
	}

	return g
}

func (o *organism) inputGenes() []*gene {
	genes := make([]*gene, 0, len(o.inputs))

	for _, id := range o.inputs {
		g := o.getGene(id)

		genes = append(genes, g)
	}

	return genes
}

func (o *organism) outputGenes() []*gene {
	genes := make([]*gene, 0, len(o.inputs))

	for _, id := range o.outputs {
		g := o.getGene(id)

		genes = append(genes, g)
	}

	return genes
}

func (o *organism) clearGenome() {
	for _, g := range o.genome {
		g.clear()
	}
}

func (o *organism) feed(inputs []float64) []float64 {
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
		h := o.getGene(g.output)
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
