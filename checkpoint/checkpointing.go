package checkpoint

import (
	"encoding/gob"
	"encoding/json"
	"github.com/kruspe/nvram/nvram"
	"io"
	"os"
)

type Checkpoint struct {
	nvram    *nvram.Nvram
	jsonFile *os.File
	gobFile  *os.File
}

type Data struct {
	Key   string
	Value string
}

func NewCheckpoint(nvram *nvram.Nvram) (*Checkpoint, error) {
	jsonFile, err := os.Create("out/checkpoint.json")
	if err != nil {
		return nil, err
	}
	gobFile, err := os.Create("out/checkpoint.gob")
	if err != nil {
		return nil, err
	}
	return &Checkpoint{
		nvram:    nvram,
		jsonFile: jsonFile,
		gobFile:  gobFile,
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
	_, err = c.jsonFile.WriteAt(marshal, 0)
	if err != nil {
		return err
	}
	return nil
}

func (c *Checkpoint) Gob(data []Data) error {
	s := make(map[string]string)
	for _, d := range data {
		s[d.Key] = d.Value
	}
	err := gob.NewEncoder(c.gobFile).Encode(s)
	if err != nil {
		return err
	}
	_, err = c.gobFile.Seek(0, io.SeekStart)
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
