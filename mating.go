package main

func mate(a, b *organism) *organism {

	// Switch if necessary so that `a` has the best performance
	if a.score < b.score {
		a, b = b, a
	}

	o := newOrganism(len(a.inputs), len(a.outputs))

	// Copy input nodes
	for _, id := range a.inputs {
		o.inputs = append(o.inputs, id)
	}

	// Copy output nodes
	for _, id := range a.outputs {
		o.outputs = append(o.outputs, id)
	}

	i, j := 0, 0

	// Copy genes and hidden nodes
	for i < len(a.oinnov) && j < len(b.oinnov) {
		var g gene

		if a.oinnov[i].innov == b.oinnov[j].innov {
			// ´a´ has the better performance so copy the gene from from `a`
			g = *a.oinnov[i]
			i = min(i+1, len(a.oinnov))
			j = min(j+1, len(b.oinnov))
		} else if a.oinnov[i].innov < b.oinnov[j].innov {
			// `a` has a gene not present in ´b´
			g = *a.oinnov[i]
			i = min(i+1, len(a.oinnov))
		} else {
			// `b` has a gene not present in ´a´
			g = *b.oinnov[i]
			j = min(j+1, len(b.oinnov))
		}

		o.add(&g)
	}

	// Handle trailing genes (if any)
	for ; i < len(a.oinnov); i++ {
		g := *a.oinnov[i]
		o.add(&g)
	}

	// Handle trailing genes (if any)
	for ; j < len(b.oinnov); j++ {
		g := *b.oinnov[i]
		o.add(&g)
	}

	return o
}
