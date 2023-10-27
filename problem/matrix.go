package problem

import (
	"errors"
	"fmt"
	"github.com/kruspe/nvram/checkpoint"
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

func (p *Problem) MultiplyGobCheckpoints(a [][]int, b [][]int) ([][]int, error) {
	if len(a[0]) != len(b) {
		return nil, errors.New("rows of a must be equal to columns of b")
	}

	data := ""
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
			data += fmt.Sprintf("%d,", r)
		}

		err := p.checkpoint.Gob.Rewrite([]checkpoint.Data{
			{
				Key:   "currentResult",
				Value: data[:len(data)-1],
			},
		})
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
	}
	return result, nil
}

func (p *Problem) MultiplyNvramCheckpoints(a [][]int, b [][]int) ([][]int, error) {
	if len(a[0]) != len(b) {
		return nil, errors.New("rows of a must be equal to columns of b")
	}

	data := ""
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
			data += fmt.Sprintf("%d,", r)
		}
		err := p.checkpoint.Nvram.Write([]checkpoint.Data{
			{
				Key:   "currentResult",
				Value: data[:len(data)-1],
			},
		})
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
	}
	return result, nil
}
