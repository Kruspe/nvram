package problem

import (
	"github.com/kruspe/nvram/checkpoint"
	"github.com/kruspe/nvram/nvram"
)

type Problem struct {
	checkpoint *checkpoint.Checkpoint
}

func NewProblem(nvram *nvram.Nvram) (*Problem, error) {
	c, err := checkpoint.NewCheckpoint(nvram)
	if err != nil {
		return nil, err
	}
	return &Problem{checkpoint: c}, nil
}
