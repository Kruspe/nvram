package main

import (
	"encoding/json"
	"fmt"
	"github.com/kruspe/nvram/nvram"
	"os"
	"strconv"
	"time"
)

type nvramHelper struct {
	prefix string
	value  string
	nvram  *nvram.Nvram

	writeFailure  []int
	deleteFailure []int
}

func NewNvramHelper(prefix, value string, nvram *nvram.Nvram) *nvramHelper {
	return &nvramHelper{
		prefix: prefix,
		value:  value,
		nvram:  nvram,
	}
}

func (h *nvramHelper) write(times int) {
	for i := 0; i < times; i++ {
		err := h.nvram.Set(fmt.Sprintf("%s%d", h.prefix, i), h.value)
		if err != nil {
			h.writeFailure = append(h.writeFailure, i)
		}
	}
}

func (h *nvramHelper) get(times int) {
	for i := 0; i < times; i++ {
		key := h.prefix + strconv.Itoa(i)
		_, err := h.nvram.Get(key)
		if err != nil {
			fmt.Println(err)
			return
		}
		//fmt.Printf("%s: %s\n", key, value)
	}
}

func (h *nvramHelper) delete(times int) {
	for i := 0; i < times; i++ {
		err := h.nvram.Delete(fmt.Sprintf("%s%d", h.prefix, i))
		if err != nil {
			h.deleteFailure = append(h.deleteFailure, i)
		}
	}
}

func main() {
	n := nvram.New()
	defer n.Teardown()
	value := ""
	// this will take up 1002 bytes
	for i := 0; i < 664; i++ {
		value += "a"
	}
	h := NewNvramHelper("test", value, n)

	startTime := time.Now()
	h.write(1)
	duration := time.Since(startTime)
	fmt.Printf("writing to NVRAM took %d μs\n", duration.Microseconds())

	startTime = time.Now()
	h.get(1)
	duration = time.Since(startTime)
	fmt.Printf("reading from NVRAM took %d μs\n", duration.Microseconds())
	h.delete(1)

	startTime = time.Now()
	data := map[string]string{"test0": value}
	d, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("out/data.json", d, 0666)
	if err != nil {
		panic(err)
	}
	duration = time.Since(startTime)
	fmt.Printf("writing to storage took %d μs\n", duration.Microseconds())

	startTime = time.Now()
	f, err := os.ReadFile("out/data.json")
	if err != nil {
		panic(err)
	}
	var test map[string]string
	err = json.Unmarshal(f, &test)
	if err != nil {
		panic(err)
	}
	duration = time.Since(startTime)
	fmt.Printf("reading from storage took %d μs\n", duration.Microseconds())

	defer func() {
		if len(h.writeFailure) > 0 {
			fmt.Printf("write failures: %v\n", h.writeFailure)
		}
		if len(h.deleteFailure) > 0 {
			fmt.Printf("delete failures: %v\n", h.deleteFailure)
		}
	}()
}
