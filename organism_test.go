package neater

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func createInputsOuputs(c *Configuration) ([]nodeID, []nodeID) {
	resetNodeID()
	inputs := make([]nodeID, c.Inputs)
	for i := range inputs {
		inputs[i] = nextNodeID()
	}
	outputs := make([]nodeID, c.Outputs)
	for i := range outputs {
		outputs[i] = nextNodeID()
	}

	return inputs, outputs
}

func TestEval(t *testing.T) {
	tests := []struct {
		name string
		conf *Configuration

		input  []float64
		expect []float64
	}{
		{
			name: "single unit",
			conf: &Configuration{
				Inputs:   1,
				Outputs:  1,
				activate: unit,
			},

			input:  []float64{1},
			expect: []float64{defaultWeight * unit(1)},
		},
		{
			name: "single sigmoid",
			conf: &Configuration{
				Inputs:   1,
				Outputs:  1,
				activate: sigmoid,
			},

			input:  []float64{1},
			expect: []float64{defaultWeight * sigmoid(1)},
		},
		{
			name: "double unit",
			conf: &Configuration{
				Inputs:   2,
				Outputs:  2,
				activate: unit,
			},

			input:  []float64{1, 2},
			expect: []float64{defaultWeight * unit(1), defaultWeight * unit(2)},
		},
		{
			name: "double sigmoid",
			conf: &Configuration{
				Inputs:   2,
				Outputs:  2,
				activate: sigmoid,
			},

			input:  []float64{1, 2},
			expect: []float64{defaultWeight * sigmoid(1), defaultWeight * sigmoid(2)},
		},
		{
			name: "single split",
			conf: &Configuration{
				Inputs:   1,
				Outputs:  2,
				activate: sigmoid,
			},

			input:  []float64{1},
			expect: []float64{defaultWeight * sigmoid(1), defaultWeight * sigmoid(1)},
		},
		{
			name: "single join",
			conf: &Configuration{
				Inputs:   2,
				Outputs:  1,
				activate: sigmoid,
			},

			input:  []float64{1, 2},
			expect: []float64{defaultWeight*sigmoid(1) + defaultWeight*sigmoid(2)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputs, outputs := createInputsOuputs(test.conf)
			o := newOrganism(test.conf, inputs, outputs,
				withConnectStrategy(connectFlow))

			output := o.Eval(test.input)
			require.Equal(t, test.expect, output)
		})
	}
}

func TestAddNotRecurrent(t *testing.T) {
	var nCount uint64
	nodeIDGenerator = func() nodeID {
		nCount++
		return nodeID(nCount)
	}

	tests := []struct {
		name   string
		conf   *Configuration
		pairs  []nodePair
		input  []float64
		output []float64
		expect []float64
	}{
		{
			//
			//      +--> 3 --+
			//		|        |
			//		|        v
			//	1 --+--------+---> 2

			name: "One additional node connecting input and output",
			conf: &Configuration{
				Inputs:   1,
				Outputs:  1,
				activate: unit,
			},
			pairs: []nodePair{
				nodePair{1, 2},
				nodePair{1, 3},
				nodePair{3, 2},
			},
			input:  []float64{1},
			expect: []float64{2},
		},
		{
			//	1 ---+        +---> 3
			//		 |        |
			//		 +--> 5 --+
			//		 |        |
			//	2 ---+        +---> 4

			name: "Two inputs join and then split to two outputs",
			conf: &Configuration{
				Inputs:   2,
				Outputs:  2,
				activate: unit,
			},
			pairs: []nodePair{
				nodePair{1, 5},
				nodePair{5, 3},
				nodePair{2, 5},
				nodePair{5, 4},
			},
			input:  []float64{1, 2},
			expect: []float64{3, 3},
		},
		{
			//	1 ---+         +---> 8 ---+
			//		 |         |          |
			//		 |         |          v
			//		 +--> 7 ---+----------+---> 4
			//		 |         |
			//	2 ---+         +--------------> 5
			//	3 ----------------------------> 6

			name: "Complex topology 1",
			conf: &Configuration{
				Inputs:   3,
				Outputs:  3,
				activate: unit,
			},
			pairs: []nodePair{
				nodePair{1, 7},
				nodePair{7, 8},
				nodePair{8, 4},
				nodePair{3, 6},
				nodePair{7, 4},
				nodePair{7, 5},
				nodePair{2, 7},
			},
			input:  []float64{2, 1, 7},
			expect: []float64{6, 3, 7},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputs, outputs := createInputsOuputs(test.conf)

			o := newOrganism(test.conf, inputs, outputs,
				withConnectStrategy(connectNone))

			// Reset node ID counter
			nCount = 0
			for _, p := range test.pairs {

				o.nodes[p.input] = 0
				o.nodes[p.output] = 0

				g := newGene(p, defaultWeight, unit)
				o.addGene(g)
			}

			output := o.Eval(test.input)
			require.Equal(t, test.expect, output)
		})
	}
}

func TestPanicCases(t *testing.T) {
	var nCount uint64
	nodeIDGenerator = func() nodeID {
		nCount++
		return nodeID(nCount)
	}

	tests := []struct {
		name  string
		conf  *Configuration
		pairs []nodePair
		input []float64
	}{
		{
			name: "case 1",
			conf: &Configuration{
				Inputs:  2,
				Outputs: 1,
			},
			pairs: []nodePair{
				nodePair{1, 3},
				nodePair{2, 3},
				nodePair{1, 6},
				nodePair{6, 3},
			},
			input: []float64{1, 1},
		},
		{
			name: "case 2",
			conf: &Configuration{
				Inputs:  2,
				Outputs: 1,
			},
			pairs: []nodePair{
				nodePair{1, 3},    // 1
				nodePair{2, 3},    // 2
				nodePair{1, 5},    // 5
				nodePair{5, 3},    // 6
				nodePair{1, 10},   //15
				nodePair{10, 3},   //16
				nodePair{5, 18},   // 31
				nodePair{18, 3},   // 32
				nodePair{1, 2319}, // 4621
				nodePair{2319, 5}, // 4622
			},
			input: []float64{1, 1},
		},
		{
			name: "case 3",
			conf: &Configuration{
				Inputs:  2,
				Outputs: 1,
			},
			pairs: []nodePair{
				nodePair{1, 3},    // 1
				nodePair{2, 3},    // 2
				nodePair{1, 5},    // 5
				nodePair{5, 3},    // 6
				nodePair{5, 18},   //31
				nodePair{18, 3},   //32
				nodePair{1, 743},  // 1451
				nodePair{743, 5},  // 1452
				nodePair{5, 735},  // 1453
				nodePair{735, 18}, // 1454
			},
			input: []float64{1, 1},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputs, outputs := createInputsOuputs(test.conf)
			o := newOrganism(test.conf, inputs, outputs,
				withConnectStrategy(connectNone))

			// Reset node ID counter
			nCount = 0
			for _, p := range test.pairs {

				o.nodes[p.input] = 0
				o.nodes[p.output] = 0

				//fmt.Printf("Add: %s\n", p)
				g := newGene(p, defaultWeight, unit)
				o.addGene(g)
				//fmt.Printf("---Iteration %d----\n%s\n\n", i+1, o)
			}

			o.Eval(test.input)
		})
	}
}

func TestMutateAddNode(t *testing.T) {
	newGene := func(p nodePair, w float64, f activationFunction, innov geneID, disabled bool) *gene {
		return &gene{
			innov:    innov,
			p:        p,
			weight:   w,
			disabled: disabled,
			activate: f,
		}
	}

	tests := []struct {
		name      string
		conf      *Configuration
		randVal   int
		nCount    uint64
		gCount    uint64
		nodeCache map[nodePair]genePair
		genes     []*gene
		expect    []*gene
	}{
		{
			name: "One inuput one output no cache",
			conf: &Configuration{
				Inputs:            1,
				Outputs:           1,
				InitialBiasWeight: 0,
				activate:          sigmoid,
			},
			randVal:   0,
			nCount:    2,
			gCount:    1,
			nodeCache: make(map[nodePair]genePair),
			genes: []*gene{
				newGene(nodePair{1, 2}, 1, sigmoid, geneID(1), false),
			},
			expect: []*gene{
				newGene(nodePair{1, 2}, 1, sigmoid, geneID(1), true),
				// Skip geneID 2 which would be assigned when connecting bias
				newGene(nodePair{1, 3}, 1, sigmoid, geneID(3), false),
				newGene(nodePair{3, 2}, 1, sigmoid, geneID(4), false),
			},
		},
		{
			name: "One inuput one output with cache",
			conf: &Configuration{
				Inputs:   1,
				Outputs:  1,
				activate: sigmoid,
			},
			randVal: 0,
			nCount:  2,
			gCount:  1,
			nodeCache: map[nodePair]genePair{
				nodePair{1, 2}: genePair{
					// Skip geneID 2 which would be assigned when connecting bias
					newGene(nodePair{1, 3}, 1, sigmoid, geneID(3), false),
					newGene(nodePair{3, 2}, 1, sigmoid, geneID(4), false),
				},
			},
			genes: []*gene{
				newGene(nodePair{1, 2}, 1, sigmoid, geneID(1), false),
			},
			expect: []*gene{
				newGene(nodePair{1, 2}, 1, sigmoid, geneID(1), true),
				// Skip geneID 2 which would be assigned when connecting bias
				newGene(nodePair{1, 3}, 1, sigmoid, geneID(3), false),
				newGene(nodePair{3, 2}, 1, sigmoid, geneID(4), false),
			},
		},
		{
			name: "Two inputs two outputs no cache",
			conf: &Configuration{
				Inputs:   2,
				Outputs:  2,
				activate: sigmoid,
			},
			randVal:   0,
			nCount:    4,
			gCount:    4,
			nodeCache: make(map[nodePair]genePair),
			genes: []*gene{
				newGene(nodePair{1, 3}, 1, sigmoid, geneID(1), false),
				newGene(nodePair{1, 4}, 1, sigmoid, geneID(2), false),
				newGene(nodePair{2, 3}, 1, sigmoid, geneID(3), false),
				newGene(nodePair{2, 4}, 1, sigmoid, geneID(4), false),
			},
			expect: []*gene{
				newGene(nodePair{1, 3}, 1, sigmoid, geneID(1), true),
				newGene(nodePair{1, 4}, 1, sigmoid, geneID(2), false),
				newGene(nodePair{2, 3}, 1, sigmoid, geneID(3), false),
				newGene(nodePair{2, 4}, 1, sigmoid, geneID(4), false),
				// Skip geneID 5 which would be assigned when connecting bias
				newGene(nodePair{1, 5}, 1, sigmoid, geneID(6), false),
				newGene(nodePair{5, 3}, 1, sigmoid, geneID(7), false),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// Reset node ID counter
			atomic.StoreUint64(&nodeIDCount, test.nCount)
			// Reset gene ID counter
			atomic.StoreUint64(&innovCount, test.gCount)

			o := newCleanOrganism(test.conf)

			for _, g := range test.genes {
				o.nodes[g.p.input] = 0
				o.nodes[g.p.output] = 0
				o.addGene(g)
			}

			randIntn = func(int) int {
				return test.randVal
			}

			o.mutateAddNode(test.nodeCache)

			require.Equal(t, len(test.expect), len(o.oinnov))
			for i, x := range test.expect {
				y := o.oinnov[i]
				require.True(t, x.equalTo(y), "Have: %#v Want: %#v", y, x)
			}
		})
	}
}
