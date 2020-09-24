package neat

import (
	"fmt"
	"sync/atomic"
)

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

	genePair struct {
		alpha *gene
		beta  *gene
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

func newGene(p nodePair, w float64, f activationFunction) *gene {
	g := &gene{
		innov:    innovIDGenerator(),
		p:        p,
		weight:   w,
		disabled: defaultDisabled,
		activate: f,
	}

	return g
}

func (g *gene) copy() *gene {
	return &gene{
		innov:    g.innov,
		p:        g.p,
		weight:   g.weight,
		disabled: g.disabled,
		activate: g.activate,
	}
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

func (g *gene) String() string {
	return fmt.Sprintf("I: %-2d W: %2.2f P: %s", g.innov, g.weight, g.p)
}

func (p nodePair) String() string {
	return fmt.Sprintf("In: %-2d Out: %-2d", p.input, p.output)
}
