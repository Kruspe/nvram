package benchmark

import (
	"fmt"
	"github.com/kruspe/nvram/nvram"
	"github.com/kruspe/nvram/problem"
	"math/rand"
	"strconv"
	"time"
)

const (
	matrixColumns = 3072
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
	_, lastCheckpointNumber, err := p.MultiplyNvramCheckpoints(a, b)
	duration = time.Since(startTime)
	fmt.Printf("With NVRAM Checkpointing took %d μs\n", duration.Microseconds())
	fmt.Println("Last checkpoint number: ", lastCheckpointNumber)
	for i := 0; i < lastCheckpointNumber; i++ {
		err := nvram.Delete("result" + fmt.Sprintf("000000%d", i)[len(strconv.Itoa(i)):])
		if err != nil {
			return err
		}
	}
	defer nvram.Delete("lastField")
	if err != nil {
		return err
	}

	return nil
}
