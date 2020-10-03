package neater

import (
	"fmt"
	"io"
	"strings"
)

func Graph(o *organism, w io.Writer) error {
	isIn := func(n nodeID, ns []nodeID) bool {
		for _, x := range ns {
			if n == x {
				return true
			}
		}

		return false
	}

	b := new(strings.Builder)

	b.Write([]byte("digraph G {\n"))
	b.Write([]byte("  concatenate=False;\n"))
	b.Write([]byte("  rankdir=LR;\n"))

	for n := range o.nodes {
		if isIn(n, o.inputs) {
			b.Write([]byte(fmt.Sprintf("    node%d [shape=circle, style=filled, color=red];\n", n)))
		}
	}

	for n := range o.nodes {
		if isIn(n, o.outputs) {
			b.Write([]byte(fmt.Sprintf("    node%d [shape=circle, style=filled, color=blue];\n", n)))
		}
	}

	for n := range o.nodes {
		if !isIn(n, o.inputs) && !isIn(n, o.outputs) {
			b.Write([]byte(fmt.Sprintf("    node%d [shape=circle];\n", n)))
		}
	}

	for _, g := range o.oinnov {
		if g.disabled {
			b.Write([]byte(fmt.Sprintf("   gene%d [shape=record, label=\"w: %.2f|disabled\"];\n", g.innov, g.weight)))
			continue
		}
		b.Write([]byte(fmt.Sprintf("   gene%d [shape=record, label=\"w: %.2f|enabled\"];\n", g.innov, g.weight)))

	}

	for _, g := range o.oinnov {
		b.Write([]byte(fmt.Sprintf("   node%d -> gene%d;\n", g.p.input, g.innov)))
		b.Write([]byte(fmt.Sprintf("   gene%d -> node%d;\n", g.innov, g.p.output)))
	}

	b.Write([]byte("  }\n"))

	_, err := w.Write([]byte(b.String()))
	return err
}
