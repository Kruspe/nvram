package main

import (
	"fmt"
	"github.com/kruspe/nvram/nvram"
)

func main() {
	n := nvram.New()
	defer n.Teardown()

	err := n.Set("test", "test")
	if err != nil {
		panic(err)
	}

	value, err := n.Get("test")
	if err != nil {
		panic(err)
	}
	fmt.Println(value)

	err = n.Delete("test")
	if err != nil {
		panic(err)
	}

	value, err = n.Get("test")
	if err != nil {
		fmt.Println(err)
	}
}
