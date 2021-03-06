package neater

import (
	"sync/atomic"
)

type (
	nodeID uint64

	node struct {
		id        nodeID
		sum       float64
		val       float64
		weight    float64
		evaluated bool
		visited   bool

		activate activationFunction
	}
)

const (
	biasID     nodeID  = 0
	biasOutput float64 = 1.0
)

var (
	nodeIDCount     = uint64(0)
	nodeIDGenerator = nextNodeID
)

func resetNodeID() {
	atomic.StoreUint64(&nodeIDCount, 0)
}

func nextNodeID() nodeID {
	return nodeID(atomic.AddUint64(&nodeIDCount, 1))
}
