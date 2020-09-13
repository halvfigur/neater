package main

import "sync/atomic"

const (
	terminal = nodeID(0)

	defaultWeight   = float64(1)
	defaultDisabled = false
)

type (
	geneID int64

	geneOpt func(*gene)

	activationFunction func(float64) float64

	nodePair struct {
		input  nodeID
		output nodeID
	}

	gene struct {
		innov geneID

		p        nodePair
		weight   float64
		disabled bool

		activate activationFunction
	}
)

var (
	innovCount                = uint64(0)
	defaultActivationFunction = sigmoid
	innovIDGenerator          = nextInnov
)

func nextInnov() geneID {
	return geneID(atomic.AddUint64(&innovCount, 1))
}

func newGene(p nodePair,
	opts ...geneOpt) *gene {
	g := &gene{
		innov:    innovIDGenerator(),
		p:        p,
		weight:   defaultWeight,
		disabled: defaultDisabled,
		activate: defaultActivationFunction,
	}

	for _, o := range opts {
		o(g)
	}

	return g
}

func (g *gene) copy() *gene {
	c := *g
	return &c
}

func (g *gene) equalTo(x *gene) bool {
	return g.innov == x.innov &&
		g.p.input == x.p.input &&
		g.p.output == x.p.output &&
		g.weight == x.weight &&
		g.disabled == x.disabled

}

func withWeight(weight float64) geneOpt {
	return func(g *gene) {
		g.weight = weight
	}
}

func withDisabled(disabled bool) geneOpt {
	return func(g *gene) {
		g.disabled = disabled
	}
}

func withActivationFunction(f activationFunction) geneOpt {
	return func(g *gene) {
		g.activate = f
	}
}
