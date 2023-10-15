package nvram

import (
	"fmt"
	"strconv"
)

type Helper struct {
	prefix string
	value  string
	nvram  *Nvram

	writeFailure  []int
	deleteFailure []int
}

func NewNvramHelper(prefix, value string, nvram *Nvram) *Helper {
	return &Helper{
		prefix: prefix,
		value:  value,
		nvram:  nvram,
	}
}

func (h *Helper) write(times int) {
	for i := 0; i < times; i++ {
		err := h.nvram.Set(fmt.Sprintf("%s%d", h.prefix, i), h.value)
		if err != nil {
			h.writeFailure = append(h.writeFailure, i)
		}
	}
}

func (h *Helper) get(times int) {
	for i := 0; i < times; i++ {
		key := h.prefix + strconv.Itoa(i)
		value, err := h.nvram.Get(key)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%s: %s\n", key, value)
	}
}

func (h *Helper) delete(times int) {
	for i := 0; i < times; i++ {
		err := h.nvram.Delete(fmt.Sprintf("%s%d", h.prefix, i))
		if err != nil {
			h.deleteFailure = append(h.deleteFailure, i)
		}
	}
}
