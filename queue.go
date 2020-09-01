package main

type (
	queue interface {
		put(*gene)
		get() *gene
		len() int
		clear()
	}

	sliceQueue []*gene
)

func newSliceQueue(initialSize int) *sliceQueue {
	if initialSize <= 0 {
		panic("Queue size must be greater than 0")
	}

	q := sliceQueue(make([]*gene, 0, initialSize))

	return &q
}

func (q *sliceQueue) put(g *gene) {
	*q = append(*q, g)
}

func (q *sliceQueue) get() *gene {
	if len(*q) == 0 {
		panic("Queue empty")
	}

	g := (*q)[0]
	*q = (*q)[1:]

	return g
}

func (q *sliceQueue) len() int {
	return len(*q)
}

func (q *sliceQueue) clear() {
	*q = (*q)[:0]
}
