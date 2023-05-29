package main

import (
	"fmt"
	"github.com/kruspe/nvram/nvram"
	"strconv"
)

type helper struct {
	prefix string
	value  string
	nvram  *nvram.Nvram

	writeFailure  []int
	deleteFailure []int
}

func main() {
	n := nvram.New()
	defer n.Teardown()
	value := ""
	// this will take up 1002 bytes
	for i := 0; i < 664; i++ {
		value += "a"
	}
	h := &helper{
		prefix: "test",
		value:  value,
		nvram:  n,
	}
	defer func() {
		if len(h.writeFailure) > 0 {
			fmt.Printf("write failures: %v\n", h.writeFailure)
		}
		if len(h.deleteFailure) > 0 {
			fmt.Printf("delete failures: %v\n", h.deleteFailure)
		}
	}()
}

func (h *helper) write(times int) {
	for i := 0; i < times; i++ {
		err := h.nvram.Set(fmt.Sprintf("%s%d", h.prefix, i), h.value)
		if err != nil {
			h.writeFailure = append(h.writeFailure, i)
		}
	}
}

func (h *helper) get(times int) {
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

func (h *helper) delete(times int) {
	for i := 0; i < times; i++ {
		err := h.nvram.Delete(fmt.Sprintf("%s%d", h.prefix, i))
		if err != nil {
			h.deleteFailure = append(h.deleteFailure, i)
		}
	}
}
