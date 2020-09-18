package neat

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

		// strategy determines how to connect the nodes during the initial
		// setup
		strategy connectStrategy

		// fitness is the organisms fitness
		fitness float64

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
	case connectFull:
		o.connectFull()
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

	// Find the last gene that ´g´ accepts input from
	inputDep := -1
	// Find the first gene that 'g' outputs to
	outputDep := -1

	for i, x := range o.oeval {
		if x.p.output == g.p.input {
			inputDep = i
		}

		if outputDep == -1 && x.p.input == g.p.output {
			outputDep = i
		}
	}

	if inputDep != -1 && outputDep != -1 {
		if !o.conf.Recurrent {
			// If recurrency is not permitted then the output dependencies must
			// occur before the input dependencies.
			if outputDep <= inputDep {
				panic("recurrence not configured")
			}
		}
	}

	if inputDep != -1 {
		// Append after the last gene that outputs to 'g'
		if inputDep == len(o.oeval)-1 {
			o.oeval = append(o.oeval, g)
			return
		}

		o.oeval = append(o.oeval[:inputDep+1], append([]*gene{g}, o.oeval[inputDep+1:]...)...)
		return
	}

	if outputDep != -1 {
		// Append before the first gene that accepts output from 'g'
		if outputDep == 0 {
			o.oeval = append([]*gene{g}, o.oeval...)
			return
		}

		o.oeval = append(o.oeval[:outputDep], append([]*gene{g}, o.oeval[outputDep:]...)...)
		return
	}

	// 'g' doesn't depend on any other gene, append at the end of the
	// evaluation order.
	o.oeval = append(o.oeval, g)
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

func (o *organism) randomNode() nodeID {
	x := randIntn(len(o.nodes))
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

// Mutation things

func (o *organism) getRandUnconnectedNodePair() nodePair {

	var p nodePair

	// Assume the nodes are already connected and keep going until we find
	// a pair that aren't connected. This may get us stuck in an infinite loop.
	alreadyConnected := true

	for alreadyConnected {
		p.input = o.randomNode()
		p.output = o.randomNode()

		// Make sure the input and output are different
		if p.input == p.output {
			continue
		}

		if !o.conf.Recurrent {
			// If reccurent connections aren't allowed then the first gene that
			// takes ´p.output´ as input must appear after the last gene that
			// outputs to 'p.input' in the evaluation order
			firstIdx := -1
			lastIdx := -1
			for i, g := range o.oeval {
				if firstIdx == -1 {
					if g.p.input == p.output {
						firstIdx = i
					}
				}

				if g.p.output == p.input {
					lastIdx = i
				}
			}

			// If there is a patch from ´p.input' to ´p.output' make sure that
			// firstIdx > lastIdx
			if firstIdx > lastIdx {
				continue
			}
		}

		// Make sure the nodes aren't already connected
		alreadyConnected = o.connected(p.input, p.output)
	}

	return p
}

func (o *organism) connectNodes(p nodePair, innovCache map[nodePair]*gene) *gene {
	if g, ok := innovCache[p]; ok {
		// This innovation has already been made
		x := g.copy()
		o.add(x)
		return x
	}

	g := newGene(p)
	innovCache[p] = g
	o.add(g)

	return g
}

func (o *organism) mutate(innovCache map[nodePair]*gene) {
	for _, g := range o.oinnov {
		if rand.Intn(o.conf.WeightMutationProb) == 0 {
			g.weight *= rand.Float64() * o.conf.WeightMutationPower
		}

		if rand.Intn(o.conf.AddNodeMutationProb) == 0 {
			p := o.getRandUnconnectedNodePair()

			o.connectNodes(p, innovCache)
		}

		if rand.Intn(o.conf.ConnectNodesMutationProb) == 0 {
			g.disabled = true

			id := nodeIDGenerator()
			o.nodes[id] = 0

			o.connectNodes(nodePair{g.p.input, id}, innovCache)
			x := o.connectNodes(nodePair{id, g.p.output}, innovCache)
			x.weight = g.weight
		}
	}
}
