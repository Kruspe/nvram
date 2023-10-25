package main

import (
	"github.com/kruspe/nvram/benchmark"
	"github.com/kruspe/nvram/nvram"
)

func main() {
	n := nvram.NewNvram()
	defer n.Teardown()

	err := benchmark.NewBenchmark(n)
	if err != nil {
		panic(err)
	}

	//err := benchmark.Multiply(n)
	//if err != nil {
	//	panic(err)
	//}
}
