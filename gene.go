package main

import "sync/atomic"

const (
	terminal = nodeID(0)

	defaultWeight   = float64(0.5)
	defaultDisabled = false
)

var (
	defaultActivationFunction = sigmoid
)

type (
	geneID int64

	geneOpt func(*gene)

	activationFunction func(float64) float64

	gene struct {
		innov geneID

		input    nodeID
		output   nodeID
		weight   float64
		disabled bool

		activate activationFunction
	}
)

var (
	innovCount = uint64(0)
)

func nextInnov() geneID {
	return geneID(atomic.AddUint64(&innovCount, 1))
}

func currInnov() geneID {
	return geneID(atomic.LoadUint64(&innovCount))
}

func newGene(input, output nodeID,
	opts ...geneOpt) *gene {
	g := &gene{
		innov:    nextInnov(),
		input:    input,
		output:   output,
		weight:   defaultWeight,
		disabled: defaultDisabled,
		activate: defaultActivationFunction,
	}

	for _, o := range opts {
		o(g)
	}

	return g
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
