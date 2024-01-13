package problem

import (
	"errors"
	"fmt"
	"github.com/kruspe/nvram/checkpoint"
	"math/rand"
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

func (p *Problem) MultiplyGobCheckpoints(a [][]int, b [][]int) ([][]int, error) {
	if len(a[0]) != len(b) {
		return nil, errors.New("rows of a must be equal to columns of b")
	}

	checkPointCounter := 0
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
		if len(data) > 1024*50 {
			err := p.checkpoint.Gob.New([]checkpoint.Data{
				{
					Key:   "result" + fmt.Sprintf("000000%d", checkPointCounter)[len(strconv.Itoa(checkPointCounter)):],
					Value: data[:len(data)-1],
				},
			})
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			data = ""
			checkPointCounter++
		}
	}

	return result, nil
}

func (p *Problem) MultiplyNvramCheckpoints(a [][]int, b [][]int) ([][]int, int, error) {
	if len(a[0]) != len(b) {
		return nil, 0, errors.New("rows of a must be equal to columns of b")
	}

	checkPointCounter := 0
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
		if len(data) > 1024*50 {
			err := p.checkpoint.Nvram.Write([]checkpoint.Data{
				{
					Key:   "result" + fmt.Sprintf("000000%d", checkPointCounter)[len(strconv.Itoa(checkPointCounter)):],
					Value: data[:len(data)-1],
				},
			})
			if err != nil {
				fmt.Println(err.Error())
				return nil, checkPointCounter, err
			}
			data = ""
			checkPointCounter++
		}
	}
	return result, checkPointCounter, nil
}

func (p *Problem) MultiplyMeasureNvramSet(a [][]int, b [][]int) ([][]int, int, error) {
	if len(a[0]) != len(b) {
		return nil, 0, errors.New("rows of a must be equal to columns of b")
	}

	var nvramSetDuration int64
	checkPointCounter := 0
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
		if len(data) > 1024*50 {
			d, err := p.checkpoint.Nvram.WriteWithMeasurement([]checkpoint.Data{
				{
					Key:   "result" + fmt.Sprintf("000000%d", checkPointCounter)[len(strconv.Itoa(checkPointCounter)):],
					Value: data[:len(data)-1],
				},
			})
			if err != nil {
				fmt.Println(err.Error())
				return nil, checkPointCounter, err
			}
			nvramSetDuration += d
			data = ""
			checkPointCounter++
		}
	}
	fmt.Println("nvram set duration:", nvramSetDuration/int64(checkPointCounter))
	return result, checkPointCounter, nil
}

func CreateMatrices(rows, columns int) ([][]int, [][]int) {
	a := make([][]int, rows)
	b := make([][]int, rows)
	for i := 0; i < rows; i++ {
		a[i] = make([]int, columns)
		b[i] = make([]int, columns)
		for j := 0; j < columns; j++ {
			a[i][j] = rand.Intn(100)
			b[i][j] = rand.Intn(100)
		}
	}
	return a, b
}
