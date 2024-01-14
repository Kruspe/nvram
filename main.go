package main

import (
	"bufio"
	"fmt"
	"github.com/kruspe/nvram/benchmark"
	"github.com/kruspe/nvram/nvram"
	"github.com/kruspe/nvram/tester"
	"os"
	"strconv"
	"time"
)

func main() {
	n := nvram.NewNvram()
	defer n.Teardown()
	for {
		var command string
		wait := true
		for wait {
			fmt.Println("Type 'b' for benchmark, 'm' for matmult, 'c' for matmult with nvram benchmark, 'l' to add large value, 'v' to add several values or 'q' to quit.")
			input, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				panic(err)
			}
			command = input
			wait = false
		}
		if command == "q\n" {
			break
		} else if command == "b\n" {
			err := benchmark.NewBenchmark(n)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
		} else if command == "m\n" {
			err := benchmark.Multiply(n)
			if err != nil {
				panic(err)
			}
		} else if command == "c\n" {
			err := benchmark.CombinedBenchmark(n)
			if err != nil {
				panic(err)
			}
		} else if command == "l\n" {
			addOneLargeValue(n)
		} else if command == "v\n" {
			addLargeAmountOfValues(n)
		}
	}
}

func addOneLargeValue(n *nvram.Nvram) {
	err := tester.StoreLargeValue(n)
	if err != nil {
		panic(err)
	}
}

func addLargeAmountOfValues(n *nvram.Nvram) {
	keys, err := tester.CheckSize(n, 1000*1000*10)
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
	fmt.Println(oldValueGet)
	fmt.Println(lastValueGet)

	fmt.Printf("First value get took %d μs\n", oldSum/len(oldValueGet))
	fmt.Printf("Last value get took %d μs\n", lastSum/len(lastValueGet))
	for _, key := range keys {
		err := n.Delete(key)
		if err != nil {
			fmt.Println(err)
		}
	}
}
