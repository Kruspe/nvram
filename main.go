package main

import (
	"fmt"
	"github.com/kruspe/nvram/nvram"
)

func main() {
	n := nvram.New()
	//n.Set("test", "test")
	fmt.Println(n.Get("test"))
}
