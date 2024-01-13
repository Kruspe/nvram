package checkpoint

import (
	"github.com/kruspe/nvram/nvram"
	"time"
)

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

func (n *NvramCheckpointing) WriteWithMeasurement(data []Data) (int64, error) {
	var duration int64
	for _, d := range data {
		now := time.Now()
		err := n.nvram.Set(d.Key, d.Value)
		duration = time.Since(now).Microseconds()
		if err != nil {
			return 0, err
		}
	}
	return duration, nil
}
