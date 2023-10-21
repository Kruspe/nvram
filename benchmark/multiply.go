package benchmark

import (
	"fmt"
	"github.com/kruspe/nvram/nvram"
	"github.com/kruspe/nvram/problem"
	"math/rand"
	"time"
)

func Multiply(nvram *nvram.Nvram) error {
	a := make([][]int, 1024)
	b := make([][]int, 1024)
	for i := 0; i < 1024; i++ {
		a[i] = make([]int, 1024)
		b[i] = make([]int, 1024)
		for j := 0; j < 1024; j++ {
			a[i][j] = rand.Intn(100)
			b[i][j] = rand.Intn(100)
		}

	}

	p, err := problem.NewProblem(nvram)
	if err != nil {
		return err
	}

	fmt.Println("Starting Matrix Multiplication")
	startTime := time.Now()
	_, err = p.Multiply(a, b)
	if err != nil {
		return err
	}
	duration := time.Since(startTime)
	fmt.Printf("Regular execution took %d μs\n", duration.Microseconds())

	startTime = time.Now()
	_, err = p.MultiplyRegularCheckpoints(a, b)
	if err != nil {
		return err
	}
	duration = time.Since(startTime)
	fmt.Printf("With Checkpointing took %d μs\n", duration.Microseconds())

	startTime = time.Now()
	_, err = p.MultiplyGobCheckpoints(a, b)
	if err != nil {
		return err
	}
	duration = time.Since(startTime)
	fmt.Printf("With gob checkpointing took %d μs\n", duration.Microseconds())

	startTime = time.Now()
	_, err = p.MultiplyNvramCheckpoints(a, b)
	if err != nil {
		return err
	}
	duration = time.Since(startTime)
	fmt.Printf("With NVRAM Checkpointing took %d μs\n", duration.Microseconds())

	return nil
}
