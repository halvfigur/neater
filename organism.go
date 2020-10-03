package neater

import (
	"fmt"
	"math/rand"
	"strings"
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
		//connections map[nodeID]map[nodeID]bool

		// strategy determines how to connect the nodes during the initial
		// setup
		strategy connectStrategy

		// fitness is the organisms fitness
		fitness float64
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
		conf:    conf,
		inputs:  make([]nodeID, conf.Inputs),
		outputs: make([]nodeID, conf.Outputs),
		oeval:   make([]*gene, 0, nNodes),
		oinnov:  make([]*gene, 0, nNodes),
		nodes:   make(map[nodeID]float64, nNodes),
		//connections: make(map[nodeID]map[nodeID]bool),
		strategy: defaultConnectStrategy,
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
	x := newCleanOrganism(o.conf)

	copy(x.inputs, o.inputs)
	copy(x.outputs, o.outputs)

	for _, g := range o.oinnov {
		x.oinnov = append(x.oinnov, g.copy())
	}

	for _, g := range o.oeval {
		x.oeval = append(x.oeval, g.copy())
	}

	x.nodes = make(map[nodeID]float64, len(o.nodes))
	for k, v := range o.nodes {
		x.nodes[k] = v
	}

	/*
		x.connections = make(map[nodeID]map[nodeID]bool)
		for i, m := range o.connections {
			outputSet := make(map[nodeID]bool)
			for o := range m {
				outputSet[o] = true
			}
			x.connections[i] = outputSet
		}
	*/

	x.strategy = o.strategy

	return x
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
			g := newGene(nodePair{in, out}, defaultWeight, o.conf.activate)
			o.add(g)
		}
	}
}

func (o *organism) connectFlow() {
	m := max(len(o.inputs), len(o.outputs))

	for i := 0; i < m; i++ {
		input := o.inputs[i%len(o.inputs)]
		output := o.outputs[i%len(o.outputs)]

		g := newGene(nodePair{input, output}, defaultWeight, o.conf.activate)
		o.add(g)
	}
}

/*
func (o *organism) connect(p nodePair) {
	if o.connections[p.input] == nil {
		o.connections[p.input] = make(map[nodeID]bool)
	}
	o.connections[p.input][p.output] = true

	if o.connections[p.output] == nil {
		o.connections[p.output] = make(map[nodeID]bool)
	}
	o.connections[p.output][p.input] = true
}

func (o *organism) connected(p nodePair) bool {
	if o.connections[p.input] == nil {
		return false
	}

	return o.connections[p.input][p.output]
}
*/

func (o *organism) add(g *gene) {
	if _, ok := o.nodes[g.p.input]; !ok {
		panic(fmt.Sprintf("node not found %d", g.p.input))
	}

	if _, ok := o.nodes[g.p.output]; !ok {
		panic(fmt.Sprintf("node not found %d", g.p.output))
	}

	g = g.copy()

	// Make note that the nodes are connected
	//o.connect(g.p)

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

	// Store the index of the last gene that share the same input node (if any)
	// so that we can inser our gene immediately after it if there's no output
	// dependency.
	lastCommonInput := -1
	_ = lastCommonInput

	lastCommonOutput := -1
	_ = lastCommonOutput

	for i, x := range o.oeval {
		if x.p.output == g.p.input {
			inputDep = i
		}

		if outputDep == -1 && x.p.input == g.p.output {
			outputDep = i
		}

		if x.p.input == g.p.input {
			lastCommonInput = i
		}

		if x.p.output == g.p.output {
			lastCommonOutput = i
		}
	}

	//fmt.Printf("inputDep: %-4d outputDep: %-4d\n", inputDep, outputDep)

	if inputDep == -1 && outputDep == -1 {
		o.oeval = append([]*gene{g}, o.oeval...)
		return
	}

	if inputDep != -1 && outputDep != -1 {
		if !o.conf.Recurrent {
			// If recurrency is not permitted then the output dependencies must
			// occur before the input dependencies.
			if outputDep <= inputDep {
				//fmt.Println("InputDep ", inputDep, " OutputDep ", outputDep)
				//fmt.Println("Add: ", g)
				//fmt.Println(o)
				panic("recurrence not configured")
			}
		}
	}

	if inputDep != -1 {
		// Insert the new gene immediately after its last dependency
		o.oeval = append(o.oeval[:inputDep+1], append([]*gene{g}, o.oeval[inputDep+1:]...)...)
		return
	}

	if outputDep != -1 {
		// Insert the new gene immediately before its first depedant
		o.oeval = append(o.oeval[:outputDep], append([]*gene{g}, o.oeval[outputDep:]...)...)
		return
	}

	if lastCommonInput != -1 {
		// 'g' doesn't depend on any other gene, append at the end of the
		// evaluation order.
		o.oeval = append(o.oeval[:lastCommonInput+1], append([]*gene{g}, o.oeval[lastCommonInput+1:]...)...)
		return
	}

	// 'g' doesn't depend on any other gene, append at the end of the
	// evaluation order.
	//o.oeval = append(o.oeval, g)
	//o.oeval = append([]*gene{g}, o.oeval...)

}

func withConnectStrategy(s connectStrategy) organismOpt {
	return func(o *organism) {
		o.strategy = s
	}
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

func (o *organism) getRecurrentNodePair() nodePair {
	firstIdx := 1 + randIntn(len(o.oeval)-1)
	lastIdx := randIntn(firstIdx)

	input := o.oeval[firstIdx].p.input
	output := o.oeval[lastIdx].p.output

	return nodePair{input, output}
}

func (o *organism) getNodePair() nodePair {
	if o.conf.Recurrent && randFloat64() < o.conf.RecurrentConnProb {
		return o.getRecurrentNodePair()
	}

	// Select the nodes from the evaluation order such that the first node
	// appears before the last, this ensures that the resulting connection is
	// not recurrent.

	// If length is N then firstIdx should be in the range [0, (N - 2))
	firstIdx := randIntn(len(o.oeval) - 1)
	// The last index should be in the range (firstIdx, N-1]
	lastIdx := firstIdx + 1 + randIntn(len(o.oeval)-firstIdx-1)

	input := o.oeval[firstIdx].p.input
	output := o.oeval[lastIdx].p.output

	return nodePair{input, output}
}

/*
func (o *organism) _getNodePair() (nodePair, bool) {

	// Create a shallow copy of all the nodes
	nodes := make(map[nodeID]float64)
	for k := range o.nodes {
		nodes[k] = 0
	}

	for len(nodes) > 0 {
		var a nodeID

		// Pick a node at random
		for a = range nodes {
			break
		}

		delete(nodes, a)

		// Iterate over the remaining nodes
		for b := range nodes {

			p := nodePair{a, b}

			if o.connected(p) {
				continue
			}

			if o.conf.Recurrent {
				return p, true
			}

			// If ´a´ should connect to ´b´, then the last gene that outputs to
			// a must appear before the first node that takes input from ´b´ in
			// the evaluation order.

			// Find the last index of the gene that outputs to 'a'
			outputDep := -1

			// Find the first index of the gene that taked input from 'b'
			inputDep := -1

			for i, g := range o.oeval {
				if g.p.output == a {
					outputDep = i
				}

				if inputDep == -1 && g.p.input == b {
					inputDep = i
				}
			}

			//
			if inputDep < outputDep {
				continue
			}

		}
	}

	return nodePair{}, false
}
*/

func (o *organism) mutateWeight() {
	i := randIntn(len(o.oinnov))
	g := o.oinnov[i]

	if g.disabled {
		// Try again next time
		return
	}

	w := rand.NormFloat64() * o.conf.WeightMutationStandardDeviation
	if w < -o.conf.WeightMutationPower {
		w = -o.conf.WeightMutationPower
	} else if w > o.conf.WeightMutationPower {
		w = o.conf.WeightMutationPower
	}

	g.weight += w
}

func (o *organism) mutateConnectedNodes(connCache map[nodePair]*gene) {
	p := o.getNodePair()

	if g, ok := connCache[p]; ok {
		// This innovation has already been made somewhere else
		o.add(g)
		return
	}

	// This is a new innovation
	g := newGene(p, defaultWeight, o.conf.activate)
	connCache[p] = g
	o.add(g)
}

func (o *organism) mutateAddNode(nodeCache map[nodePair]genePair) {
	i := randIntn(len(o.oinnov))
	g := o.oinnov[i]

	if g.disabled {
		// Try again next time
		return
	}

	pair, ok := nodeCache[g.p]
	if ok {
		// This innovation has already been made somewhere else
		o.nodes[pair.alpha.p.output] = 0
		o.add(pair.alpha)
		o.add(pair.beta)
		g.disabled = true

		return
	}

	id := nodeIDGenerator()
	o.nodes[id] = 0

	alpha := newGene(nodePair{g.p.input, id}, defaultWeight, o.conf.activate)
	beta := newGene(nodePair{id, g.p.output}, g.weight, o.conf.activate)
	nodeCache[g.p] = genePair{alpha, beta}

	o.add(alpha)
	o.add(beta)

	g.disabled = true
}

func (o *organism) mutate(connCache map[nodePair]*gene, nodeCache map[nodePair]genePair) {
	if randFloat64() < o.conf.WeightMutationProb {
		o.mutateWeight()
	}

	if randFloat64() < o.conf.ConnectNodesMutationProb {
		o.mutateConnectedNodes(connCache)
	}

	if randFloat64() < o.conf.AddNodeMutationProb {
		o.mutateAddNode(nodeCache)
	}
}

func (o *organism) String() string {
	l := make([]string, 0, 16)

	for _, g := range o.oeval {
		l = append(l, g.String())
	}
	return strings.Join(l, "\n")
}
