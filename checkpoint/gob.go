package checkpoint

import (
	"encoding/gob"
	"fmt"
	"os"
)

type GobCheckpointing struct {
}

const filePath = "out/checkpoint.gob"

func NewGobCheckpointing() (*GobCheckpointing, error) {
	_, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	return &GobCheckpointing{}, nil
}

func (g *GobCheckpointing) Rewrite(data []Data) error {
	f, err := os.OpenFile(filePath, os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	err = f.Truncate(0)
	if err != nil {
		return err
	}

	m := make(map[string]string)
	for _, d := range data {
		m[d.Key] = d.Value
	}
	return gob.NewEncoder(f).Encode(m)
}

func (g *GobCheckpointing) New(data []Data) error {
	f, err := os.Create(fmt.Sprintf("out/%s.gob", data[0].Key))
	if err != nil {
		return err
	}
	m := make(map[string]string)
	for _, d := range data {
		m[d.Key] = d.Value
	}
	return gob.NewEncoder(f).Encode(m)
}
