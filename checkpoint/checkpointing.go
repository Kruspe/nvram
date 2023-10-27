package checkpoint

import (
	"github.com/kruspe/nvram/nvram"
)

type Checkpoint struct {
	Nvram *NvramCheckpointing
	Gob   *GobCheckpointing
}

type Data struct {
	Key   string
	Value string
}

func NewCheckpoint(nvram *nvram.Nvram) (*Checkpoint, error) {
	gob, err := NewGobCheckpointing()
	if err != nil {
		return nil, err
	}
	return &Checkpoint{
		Nvram: NewNvramCheckpointing(nvram),
		Gob:   gob,
	}, nil
}
