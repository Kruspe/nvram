package problem

import (
	"errors"
	"fmt"
	"github.com/kruspe/nvram/checkpoint"
	"strconv"
)

func (p *Problem) Multiply(a [][]int, b [][]int) ([][]int, error) {
	if len(a[0]) != len(b) {
		return nil, errors.New("rows of a must be equal to columns of b")
	}

	result := make([][]int, len(a))
	for i := range result {
		result[i] = make([]int, len(b[0]))
	}
	for i := range result {
		for j := range b[0] {
			var r = 0
			for x := range b {
				r += a[i][x] * b[x][j]
			}
			result[i][j] = r
		}

	}
	return result, nil
}

func (p *Problem) MultiplyRegularCheckpoints(a [][]int, b [][]int) ([][]int, error) {
	if len(a[0]) != len(b) {
		return nil, errors.New("rows of a must be equal to columns of b")
	}

	result := make([][]int, len(a))
	for i := range result {
		result[i] = make([]int, len(b[0]))
	}
	for i := range result {
		for j := range b[0] {
			var r = 0
			for x := range b {
				r += a[i][x] * b[x][j]
			}
			result[i][j] = r
			currentStep := i*len(b) + j
			data := createCurrentData(result, currentStep)
			err := p.checkpoint.Regular([]checkpoint.Data{
				{
					Key:   "currentResult",
					Value: data,
				},
			})
			if err != nil {
				fmt.Println(err.Error())
			}
		}

	}
	return result, nil
}

func (p *Problem) MultiplyGobCheckpoints(a [][]int, b [][]int) ([][]int, error) {
	if len(a[0]) != len(b) {
		return nil, errors.New("rows of a must be equal to columns of b")
	}

	result := make([][]int, len(a))
	for i := range result {
		result[i] = make([]int, len(b[0]))
	}
	for i := range result {
		for j := range b[0] {
			var r = 0
			for x := range b {
				r += a[i][x] * b[x][j]
			}
			result[i][j] = r
			currentStep := i*len(b) + j
			data := createCurrentData(result, currentStep)
			err := p.checkpoint.Gob([]checkpoint.Data{
				{
					Key:   "currentResult",
					Value: data,
				},
			})
			if err != nil {
				fmt.Println(err.Error())
			}
		}

	}
	return result, nil
}

func (p *Problem) MultiplyNvramCheckpoints(a [][]int, b [][]int) ([][]int, error) {
	if len(a[0]) != len(b) {
		return nil, errors.New("rows of a must be equal to columns of b")
	}

	result := make([][]int, len(a))
	for i := range result {
		result[i] = make([]int, len(b[0]))
	}
	for i := range result {
		for j := range b[0] {
			var r = 0
			for x := range b {
				r += a[i][x] * b[x][j]
			}
			result[i][j] = r
			currentStep := i*len(b) + j
			data := createCurrentData(result, currentStep)
			err := p.checkpoint.Nvram([]checkpoint.Data{
				{
					Key:   "currentResult",
					Value: data,
				},
			})
			if err != nil {
				fmt.Println(err.Error())
			}
		}

	}
	return result, nil
}

func createCurrentData(result [][]int, currentStep int) string {
	var data string
	for x, row := range result {
		for y, field := range row {
			if x*len(result)+y > currentStep {
				return data
			}
			if x == 0 && y == 0 {
				data += strconv.Itoa(field)
			} else {
				data += fmt.Sprintf(",%d", field)
			}
		}
	}
	return data
}
