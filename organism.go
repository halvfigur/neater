package neater

import (
	"fmt"
	"math/rand"
	"strings"
	"sync/atomic"
)

type (
	organismID      uint64
	connectStrategy int

	organism struct {
		id organismID

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
		// obias holds all bias connections
		obias []*gene
		// nodes holds all the nodes values
		nodes map[nodeID]float64

		// terminalNodes the set of input and output nodeIDs
		terminalNodes map[nodeID]bool

		// strategy determines how to connect the nodes during the initial
		// setup
		strategy connectStrategy

		// fitness is the organism's fitness
		fitness float64
	}

	organismOpt func(*organism)
)

var (
	organismCount = uint64(0)
)

const (
	connectNone = connectStrategy(iota)
	connectFull
	connectFlow
	connectRandom

	defaultConnectStrategy = connectFull
)

func nextOrganismID() organismID {
	return organismID(atomic.AddUint64(&organismCount, 1))
}

func newCleanOrganism(conf *Configuration) *organism {
	if conf.Inputs <= 0 {
		panic("Number of inputs must be greater than 0")
	}

	if conf.Outputs <= 0 {
		panic("Number of outputs must be greater than 0")
	}

	// Total number of initial nodes is #inputs + #outputs + bias
	nNodes := conf.Inputs*conf.Outputs + 1
	return &organism{
		id:            nextOrganismID(),
		conf:          conf,
		inputs:        make([]nodeID, conf.Inputs),
		outputs:       make([]nodeID, conf.Outputs),
		oinnov:        make([]*gene, 0, nNodes),
		oeval:         make([]*gene, 0, nNodes),
		obias:         make([]*gene, 0, conf.Outputs),
		nodes:         make(map[nodeID]float64, nNodes),
		terminalNodes: make(map[nodeID]bool),
		strategy:      defaultConnectStrategy,
	}
}

func newOrganism(conf *Configuration, inputs, outputs []nodeID, opts ...organismOpt) *organism {
	o := newCleanOrganism(conf)

	for _, opt := range opts {
		opt(o)
	}

	o.inputs = make([]nodeID, len(inputs))
	copy(o.inputs, inputs)
	for _, id := range o.inputs {
		// The input nodes will not be biased so don't call addNode
		o.nodes[id] = 0
		o.terminalNodes[id] = true
	}

	o.outputs = make([]nodeID, len(outputs))
	copy(o.outputs, outputs)
	for _, id := range o.outputs {
		o.nodes[id] = 0
		o.terminalNodes[id] = true
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

	for _, g := range o.obias {
		x.obias = append(x.obias, g.copy())
	}

	x.nodes = make(map[nodeID]float64, len(o.nodes))
	for k, v := range o.nodes {
		x.nodes[k] = v
	}

	x.terminalNodes = make(map[nodeID]bool, len(o.terminalNodes))
	for k, v := range o.terminalNodes {
		x.terminalNodes[k] = v
	}

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
	// Connect input to putput
	for _, in := range o.inputs {
		for _, out := range o.outputs {
			g := newGene(nodePair{in, out}, defaultWeight, o.conf.activate)
			o.addGene(g)
		}
	}
}

func (o *organism) connectFlow() {
	m := max(len(o.inputs), len(o.outputs))

	for i := 0; i < m; i++ {
		input := o.inputs[i%len(o.inputs)]
		output := o.outputs[i%len(o.outputs)]

		g := newGene(nodePair{input, output}, defaultWeight, o.conf.activate)
		o.addGene(g)
	}
}

func (o *organism) addBias(id nodeID) {
	// Disallow adding bias to terminal nodes
	if o.terminalNodes[id] {
		return
	}

	p := nodePair{biasID, id}

	// Check that the node isn't already biased
	for _, g := range o.obias {
		if g.p == p {
			return
		}
	}

	//g := newGene(p, defaultWeight, o.conf.activate)
	g := newGene(p, o.conf.InitialBiasWeight, o.conf.activate)
	o.obias = append(o.obias, g)
}

func (o *organism) addNode(id nodeID) {
	o.nodes[id] = 0
	o.addBias(id)
}

func (o *organism) addGene(g *gene) {
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

func (o *organism) Eval(input []float64) []float64 {
	if len(input) != len(o.inputs) {
		panic("Length of input vector must equal number of input nodes")
	}

	for i, id := range o.inputs {
		o.nodes[id] = input[i]
	}

	// Initialize each node with the corresponding weighted bias value
	for _, g := range o.obias {
		o.nodes[g.p.output] = g.activate(biasOutput) * g.weight
	}

	// Clear the output nodes
	for _, id := range o.outputs {
		o.nodes[id] = 0
	}

	// Iterate over the gene evaluation order and update the nodes accordingly
	for _, g := range o.oeval {
		if g.disabled {
			// Skip disabled genes
			continue
		}

		input := o.nodes[g.p.input]

		v := g.activate(input) * g.weight

		o.nodes[g.p.output] += v
	}

	// Copy the output nodes to the output slice
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

	// Randomly select two genes, ´first´ and ´last´, from the gene evaluation
	// order such that ´first´ occurs before ´last´. Choose the input node from
	// ´first´ and the output node from ´last´. This ensures that the
	// connection is not recurrent.

	// If length is N then firstIdx should be in the range [0, (N - 1))
	n := len(o.oeval)
	firstIdx := randIntn(n - 1)
	// The last index should be in the range (firstIdx, N-1]
	lastIdx := firstIdx + 1 + randIntn(n-(firstIdx+1))

	input := o.oeval[firstIdx].p.input
	output := o.oeval[lastIdx].p.output

	return nodePair{input, output}
}

func (o *organism) mutateWeight() {
	for _, g := range o.oinnov {
		if g.disabled {
			continue
		}

		if randFloat64() > o.conf.WeightMutationProb {
			continue
		}

		w := rand.NormFloat64() * o.conf.WeightMutationStandardDeviation
		// Clamp the weight modification so that it doesn't exceed the weight
		// mutation power
		if w < -o.conf.WeightMutationPower {
			w = -o.conf.WeightMutationPower
		} else if w > o.conf.WeightMutationPower {
			w = o.conf.WeightMutationPower
		}

		g.weight += w
	}
}

func (o *organism) mutateConnectNodes(connCache map[nodePair]*gene) {
	p := o.getNodePair()

	for _, g := range o.oinnov {
		if g.p == p {
			// These nodes are already connected, try again next time
			return
		}
	}

	if g, ok := connCache[p]; ok {
		// This innovation has already been made somewhere else
		o.addGene(g)
		return
	}

	// This is a new innovation
	g := newGene(p, defaultWeight, o.conf.activate)
	connCache[p] = g
	o.addGene(g)
}

func (o *organism) mutateAddNode(nodeCache map[nodePair]genePair) {
	// When adding a new node don't consider genes involving the bias node
	i := randIntn(len(o.oinnov))
	g := o.oinnov[i]

	if g.disabled {
		// Try again next time
		return
	}

	p, ok := nodeCache[g.p]
	if ok {
		// This innovation has already been made somewhere else
		o.nodes[p.alpha.p.output] = 0
		o.addGene(p.alpha)
		o.addGene(p.beta)
		g.disabled = true

		o.addBias(p.alpha.p.output)
		return
	}

	id := nodeIDGenerator()
	o.addNode(id)

	alpha := newGene(nodePair{g.p.input, id}, defaultWeight, o.conf.activate)
	beta := newGene(nodePair{id, g.p.output}, g.weight, o.conf.activate)
	nodeCache[g.p] = genePair{alpha, beta}
	o.addGene(alpha)
	o.addGene(beta)
	g.disabled = true
}

func (o *organism) mutate(connCache map[nodePair]*gene, nodeCache map[nodePair]genePair) {
	o.mutateWeight()

	if randFloat64() < o.conf.ConnectNodesMutationProb {
		o.mutateConnectNodes(connCache)
	}

	if randFloat64() < o.conf.AddNodeMutationProb {
		o.mutateAddNode(nodeCache)
	}
}

func (o *organism) isDisjoint() bool {
	isIn := func(n nodeID, ns []nodeID) bool {
		for _, x := range ns {
			if n == x {
				return true
			}
		}
		return false
	}

	nodes := make(map[nodeID]bool)
	for id := range o.nodes {
		nodes[id] = false
	}

	for _, g := range o.oinnov {
		nodes[g.p.output] = true
	}

	for n, v := range nodes {
		if !v && !isIn(n, o.inputs) {
			return true
		}
	}

	return false
}

func (o *organism) String() string {
	l := make([]string, 0, 16)

	for _, g := range o.oeval {
		l = append(l, g.String())
	}
	return strings.Join(l, "\n")
}
