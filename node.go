package neat

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

var (
	nodeIDCount     = uint64(0)
	nodeIDGenerator = nextNodeID
)

func nextNodeID() nodeID {
	return nodeID(atomic.AddUint64(&nodeIDCount, 1))
}
