package main

import "testing"

func TestMating(t *testing.T) {
	type insert struct {
		target int
		from   nodeID
		to     nodeID
	}

	tests := []struct {
		name    string
		inputs  int
		outputs int
		inserts []insert
	}{
		{
			name:    "sune",
			inputs:  1,
			outputs: 1,
			inserts: []insert{
				{
					target: 0,
					from:   nodeID(1),
					to:     nodeID(2),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
		})
	}
}
