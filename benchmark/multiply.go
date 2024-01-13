package benchmark

import (
	"bufio"
	"fmt"
	"github.com/kruspe/nvram/nvram"
	"github.com/kruspe/nvram/problem"
	"os"
	"strconv"
	"time"
)

const (
	matrixColumns = 2048
	matrixRows
)

func Multiply(nvram *nvram.Nvram) error {
	a, b := problem.CreateMatrices(matrixRows, matrixColumns)

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
