package checkpoint

import "github.com/kruspe/nvram/nvram"

type NvramCheckpointing struct {
	nvram *nvram.Nvram
}

func NewNvramCheckpointing(nvram *nvram.Nvram) *NvramCheckpointing {
	return &NvramCheckpointing{nvram: nvram}
}

func (n *NvramCheckpointing) Write(data []Data) error {
	for _, d := range data {
		err := n.nvram.Set(d.Key, d.Value)
		if err != nil {
			return err
		}
	}
	return nil
}
