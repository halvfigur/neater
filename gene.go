package main

import "sync/atomic"

const (
	terminal = geneID(0)

	defaultOutput   = terminal
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

		output   geneID
		weight   float64
		disabled bool
		sum      float64

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

func newGene(opts ...geneOpt) *gene {
	g := &gene{
		innov:    nextInnov(),
		output:   defaultOutput,
		weight:   defaultWeight,
		disabled: defaultDisabled,
		activate: defaultActivationFunction,
	}

	for _, o := range opts {
		o(g)
	}

	return g
}

func withOutput(output geneID) geneOpt {
	return func(g *gene) {
		g.output = output
	}
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

func (g *gene) add(v float64) {
	g.sum += v
}

func (g *gene) val() float64 {
	return g.activate(g.sum)
}

func (g *gene) clear() {
	g.sum = 0
}

func (g *gene) terminal() bool {
	return g.output == terminal
}
