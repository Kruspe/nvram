package benchmark

import (
	"fmt"
	"github.com/kruspe/nvram/nvram"
	"github.com/kruspe/nvram/problem"
	"math/rand"
	"time"
)

const (
	matrixColumns = 200
	matrixRows
)

func Multiply(nvram *nvram.Nvram) error {
	a := make([][]int, matrixRows)
	b := make([][]int, matrixRows)
	for i := 0; i < matrixRows; i++ {
		a[i] = make([]int, matrixColumns)
		b[i] = make([]int, matrixColumns)
		for j := 0; j < matrixColumns; j++ {
			a[i][j] = rand.Intn(100)
			b[i][j] = rand.Intn(100)
		}
	}

	p, err := problem.NewProblem(nvram)
	defer nvram.Delete("currentResult")
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
