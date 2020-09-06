package main

import (
	"fmt"
)

type (
	connectStrategy int

	organism struct {
		// input holds input node IDs
		inputs []nodeID
		// output holds output node IDs
		outputs []nodeID
		// order holds the gene evalulauation order
		order []*gene
		// nodes holds all the nodes values
		nodes map[nodeID]float64

		// genome holds the organism genome
		genome map[geneID]*gene

		// strategy determines how to connect the nodes during the initial
		// setup
		strategy connectStrategy

		// score is the organisms score
		score float64

		// activate is the activation function to use when evealuating node
		// output values
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

var (
	nodeCount = uint64(0)
)

func newOrganism(inputs, outputs int, opts ...organismOpt) *organism {
	if inputs <= 0 {
		panic("Number of inputs must be greater than 0")
	}

	if outputs <= 0 {
		panic("Number of outputs must be greater than 0")
	}

	nNodes := inputs * outputs
	o := &organism{
		order:    make([]*gene, 0, nNodes),
		nodes:    make(map[nodeID]float64, nNodes),
		genome:   make(map[geneID]*gene, nNodes),
		strategy: defaultConnectStrategy,
	}

	for _, opt := range opts {
		opt(o)
	}

	o.inputs = make([]nodeID, inputs)
	for i := 0; i < inputs; i++ {
		id := nextNodeID()

		o.inputs[i] = id
		o.nodes[id] = 0
	}

	o.outputs = make([]nodeID, outputs)
	for i := 0; i < outputs; i++ {
		id := nextNodeID()

		o.outputs[i] = id
		o.nodes[id] = 0
	}

	o.connectTerminals()

	return o
}

func (o *organism) connectTerminals() {
	switch o.strategy {
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
		input := o.inputs[i%len(o.inputs)]
		output := o.outputs[i%len(o.outputs)]

		g := newGene(input, output,
			withActivationFunction(o.activate))
		o.genome[g.innov] = g
		o.order = append(o.order, g)
	}
}

func withGlobalActivationFunction(f activationFunction) organismOpt {
	return func(o *organism) {
		o.activate = f
	}
}

func withConnectStrategy(s connectStrategy) organismOpt {
	return func(o *organism) {
		o.strategy = s
	}
}

func (o *organism) gene(id geneID) *gene {
	g, ok := o.genome[id]
	if !ok {
		panic(fmt.Sprintf("Gene not found %d", id))
	}

	return g
}

func (o *organism) clearNodes() {
	for id := range o.nodes {
		o.nodes[id] = 0
	}
}

func (o *organism) Eval(input []float64) []float64 {
	if len(input) != len(o.inputs) {
		panic("Length of input vector must equal number of input nodes")
	}

	o.clearNodes()

	for i, id := range o.inputs {
		o.nodes[id] = input[i]
	}

	for _, g := range o.order {
		input := o.nodes[g.input]

		v := g.activate(input) * g.weight

		o.nodes[g.output] += v
	}

	output := make([]float64, len(o.outputs))
	for i, id := range o.outputs {
		output[i] = o.nodes[id]
	}

	return output
}
