package main

import (
	"fmt"
	"github.com/kruspe/nvram/nvram"
	"github.com/kruspe/nvram/tester"
	"strconv"
	"time"
)

func main() {
	n := nvram.NewNvram()
	defer n.Teardown()

	keys, err := tester.CheckSize(n, 80000)
	if err != nil {
		panic(err)
	}
	var oldValueGet []int64
	var lastValueGet []int64
	counter := 0
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("000000%d", counter)[len(strconv.Itoa(counter)):]
		now := time.Now()
		_, err := n.Get(key)
		if err != nil {
			panic(err)
		}
		oldValueGet = append(oldValueGet, time.Since(now).Microseconds())
		counter++

		now = time.Now()
		_, err = n.Get(keys[len(keys)-1])
		if err != nil {
			panic(err)
		}
		lastValueGet = append(lastValueGet, time.Since(now).Microseconds())
	}

	oldSum := 0
	lastSum := 0
	for _, v := range oldValueGet {
		oldSum += int(v)
	}
	for _, v := range lastValueGet {
		lastSum += int(v)
	}

	fmt.Printf("First value get took %d μs\n", oldSum/len(oldValueGet))
	fmt.Printf("Last value get took %d μs\n", lastSum/len(lastValueGet))
	for _, key := range keys {
		err := n.Delete(key)
		if err != nil {
			fmt.Println(err)
		}
	}

	//err := benchmark.NewBenchmark(n)
	//if err != nil {
	//	panic(err)
	//}

	//err := benchmark.Multiply(n)
	//if err != nil {
	//	panic(err)
	//}
}
