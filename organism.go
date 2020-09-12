package main

import (
	"fmt"
	"math/rand"
)

type (
	connectStrategy int

	organism struct {
		// conf is the global configuration
		conf *Configuration

		// input holds input node IDs
		inputs []nodeID
		// output holds output node IDs
		outputs []nodeID
		// oinnov holds the gene innovation order
		oinnov []*gene
		// oeval holds the gene evalulauation order
		oeval []*gene
		// nodes holds all the nodes values
		nodes map[nodeID]float64
		// connections holds all the input to output connections
		connections map[nodeID]map[nodeID]bool

		// recurrence determines if recurrent nodes are permitted
		recurrence bool

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
	connectFull
	connectFlow
	connectRandom

	defaultConnectStrategy = connectFull
)

func newCleanOrganism(conf *Configuration) *organism {
	if conf.Inputs <= 0 {
		panic("Number of inputs must be greater than 0")
	}

	if conf.Outputs <= 0 {
		panic("Number of outputs must be greater than 0")
	}

	nNodes := conf.Inputs * conf.Outputs
	return &organism{
		conf:        conf,
		inputs:      make([]nodeID, conf.Inputs),
		outputs:     make([]nodeID, conf.Outputs),
		oeval:       make([]*gene, 0, nNodes),
		oinnov:      make([]*gene, 0, nNodes),
		nodes:       make(map[nodeID]float64, nNodes),
		connections: make(map[nodeID]map[nodeID]bool),
		strategy:    defaultConnectStrategy,
	}
}

func newOrganism(conf *Configuration, opts ...organismOpt) *organism {
	o := newCleanOrganism(conf)

	for _, opt := range opts {
		opt(o)
	}

	for i := range o.inputs {
		id := nodeIDGenerator()

		o.inputs[i] = id
		o.nodes[id] = 0
	}

	for i := range o.outputs {
		id := nodeIDGenerator()

		o.outputs[i] = id
		o.nodes[id] = 0
	}

	o.connectTerminals()

	return o
}

func (o *organism) copy() *organism {
	x := *o

	(&x).oeval = make([]*gene, len(o.oeval))
	copy(x.oeval, o.oeval)

	(&x).oinnov = make([]*gene, len(o.oinnov))
	copy(x.oinnov, o.oinnov)

	(&x).nodes = make(map[nodeID]float64, len(o.nodes))
	for k, v := range o.nodes {
		x.nodes[k] = v
	}

	(&x).connections = make(map[nodeID]map[nodeID]bool)
	for i, m := range o.connections {
		outputSet := make(map[nodeID]bool)
		for o := range m {
			outputSet[o] = true
		}
		(&x).connections[i] = outputSet
	}

	return &x
}

func (o *organism) connectTerminals() {
	switch o.strategy {
	case connectNone:
		o.connectFull()
	case connectFull:
	case connectFlow:
		o.connectFlow()
	case connectRandom:
		panic("Not implemented")
	}
}

// connectFull connects each input node to every output node.
func (o *organism) connectFull() {
	for _, in := range o.inputs {
		for _, out := range o.outputs {
			g := newGene(nodePair{in, out},
				withActivationFunction(o.activate))
			o.add(g)
		}
	}
}

func (o *organism) connectFlow() {
	m := max(len(o.inputs), len(o.outputs))

	for i := 0; i < m; i++ {
		input := o.inputs[i%len(o.inputs)]
		output := o.outputs[i%len(o.outputs)]

		g := newGene(nodePair{input, output},
			withActivationFunction(o.activate))
		o.add(g)
	}
}

func (o *organism) connected(input, output nodeID) bool {
	if o.connections[input] == nil {
		return false
	}

	return o.connections[input][output]
}

func (o *organism) add(g *gene) {
	if _, ok := o.nodes[g.p.input]; !ok {
		panic(fmt.Sprintf("node not found %d", g.p.input))
	}

	if _, ok := o.nodes[g.p.output]; !ok {
		panic(fmt.Sprintf("node not found %d", g.p.output))
	}

	o.nodes[g.p.input] = 0
	o.nodes[g.p.output] = 0

	// Make note that the nodes are connected
	if o.connections[g.p.input] == nil {
		o.connections[g.p.input] = make(map[nodeID]bool)
	}
	o.connections[g.p.input][g.p.output] = true

	// Add gene at end of innovation order
	o.oinnov = append(o.oinnov, g)

	// The rest of the function deals with insering the node at an approriate
	// place in the evalutaion order

	// If ´g´ is the first gene then just insert it and where're done.
	if len(o.oinnov) == 1 {
		o.oeval = append(o.oeval, g)
		return
	}

	var i int
	var x *gene

	// Store the position of the first gene in the evaluation order for which
	// the input node is the output node of ´g´ and use it later to test if ´g´
	// introduces recurrence
	var firstDep int
	var firstDepFound bool
	for i, x = range o.oeval {
		if !firstDepFound && x.p.input == g.p.output {
			firstDep = i
			firstDepFound = true
			break
		}
	}

	var lastDep int

	// Store the position of the last gene in the evaluation order for which
	// the output node is the input node of ´g´.
	for i, x = range o.oeval[i+1:] {
		if x.p.output == g.p.input {
			lastDep = i
		}
	}

	// If a node that depends on the output of ´g´ exists prior to 'g' in the
	// evaluation order then we have recurrence.
	if firstDepFound && firstDep < lastDep && !o.recurrence {
		panic("recurrence not configured")
	}

	// Calculate the insert position
	var p int
	if lastDep != 0 {
		p = lastDep - 1
	}

	// Insert 'g'
	o.oeval = append(o.oeval[:p], append([]*gene{g}, o.oeval[p:]...)...)
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

func withRecurrence(r bool) organismOpt {
	return func(o *organism) {
		o.recurrence = r
	}
}

func (o *organism) randomNode() nodeID {
	x := rand.Intn(len(o.nodes))
	for id := range o.nodes {
		if x == 0 {
			return id
		}

		x--
	}

	// This will never happen but the compiler can't figure that out
	return nodeID(0)
}

func (o *organism) clear() {
	for id := range o.nodes {
		o.nodes[id] = 0
	}
}

func (o *organism) Eval(input []float64) []float64 {
	if len(input) != len(o.inputs) {
		panic("Length of input vector must equal number of input nodes")
	}

	o.clear()

	for i, id := range o.inputs {
		o.nodes[id] = input[i]
	}

	for _, g := range o.oeval {
		if g.disabled {
			// Skip disabled genes
			continue
		}

		input := o.nodes[g.p.input]

		v := g.activate(input) * g.weight

		o.nodes[g.p.output] += v
	}

	output := make([]float64, len(o.outputs))
	for i, id := range o.outputs {
		output[i] = o.nodes[id]
	}

	return output
}
