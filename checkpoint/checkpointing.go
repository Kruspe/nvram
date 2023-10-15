package checkpoint

import (
	"encoding/json"
	"github.com/kruspe/nvram/nvram"
	"os"
)

type Checkpoint struct {
	nvram *nvram.Nvram
	file  *os.File
}

type Data struct {
	Key   string
	Value string
}

func NewCheckpoint(nvram *nvram.Nvram) (*Checkpoint, error) {
	file, err := os.Create("out/checkpoint.json")
	if err != nil {
		return nil, err
	}
	return &Checkpoint{
		nvram: nvram,
		file:  file,
	}, nil
}

func (c *Checkpoint) Regular(data []Data) error {
	s := make(map[string]string)
	for _, d := range data {
		s[d.Key] = d.Value
	}
	marshal, err := json.Marshal(s)
	if err != nil {
		return err
	}
	_, err = c.file.WriteAt(marshal, 0)
	if err != nil {
		return err
	}
	return nil
}

func (c *Checkpoint) Nvram(data []Data) error {
	for _, d := range data {
		err := c.nvram.Set(d.Key, d.Value)
		if err != nil {
			return err
		}
	}
	return nil
}
