package main

import (
	"github.com/kruspe/nvram/benchmark"
	"github.com/kruspe/nvram/nvram"
)

func main() {
	n := nvram.NewNvram()
	defer n.Teardown()

	//keys, err := tester.CheckSize(n, 1000)
	//if err != nil {
	//	panic(err)
	//}
	//for _, key := range keys {
	//	err := n.Delete(key)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}

	//err := benchmark.NewBenchmark(n)
	//if err != nil {
	//	panic(err)
	//}

	err := benchmark.Multiply(n)
	if err != nil {
		panic(err)
	}
}
