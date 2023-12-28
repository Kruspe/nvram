package benchmark

import (
	"bufio"
	"fmt"
	"github.com/kruspe/nvram/nvram"
	"github.com/kruspe/nvram/problem"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	matrixColumns = 2048
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

	var command string
	for {
		wait := true
		for wait {
			fmt.Println("Type 'c' to add a value or 'q' to quit.")
			input, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				return err
			}
			if input == "q\n" {
				return nil
			} else {
				command = input
				wait = false
			}
		}

		if command == "n\n" {
			fmt.Println("Starting Matrix Multiplication")
			startTime := time.Now()
			_, err = p.Multiply(a, b)
			if err != nil {
				return err
			}
			duration := time.Since(startTime)
			fmt.Printf("Regular execution took %d μs\n", duration.Microseconds())
		} else if command == "gob\n" {
			fmt.Println("Starting Matrix Multiplication")
			startTime := time.Now()
			_, err = p.MultiplyGobCheckpoints(a, b)
			if err != nil {
				return err
			}
			duration := time.Since(startTime)
			fmt.Printf("With gob checkpointing took %d μs\n", duration.Microseconds())
		} else if command == "nvram\n" {
			fmt.Println("Starting Matrix Multiplication")
			startTime := time.Now()
			_, lastCheckpointNumber, err := p.MultiplyNvramCheckpoints(a, b)
			duration := time.Since(startTime)
			fmt.Printf("With NVRAM Checkpointing took %d μs\n", duration.Microseconds())
			for i := 0; i < lastCheckpointNumber; i++ {
				err := nvram.Delete("result" + fmt.Sprintf("000000%d", i)[len(strconv.Itoa(i)):])
				if err != nil {
					return err
				}
			}
			if err != nil {
				return err
			}
		}
	}
}
