package problem

import (
	"errors"
	"fmt"
	"github.com/kruspe/nvram/checkpoint"
	"strconv"
	"time"
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

	var durations []int64
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
			now := time.Now()
			err := p.checkpoint.Gob.New([]checkpoint.Data{
				{
					Key:   "result" + fmt.Sprintf("000000%d", checkPointCounter)[len(strconv.Itoa(checkPointCounter)):],
					Value: data[:len(data)-1],
				},
			})
			durations = append(durations, time.Since(now).Microseconds())

			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			data = ""
			checkPointCounter++
		}
	}
	fmt.Println(durations)

	return result, nil
}

func (p *Problem) MultiplyNvramCheckpoints(a [][]int, b [][]int) ([][]int, int, error) {
	if len(a[0]) != len(b) {
		return nil, 0, errors.New("rows of a must be equal to columns of b")
	}

	var durations []int64
	checkPointCounter := 0
	data := ""
	result := make([][]int, len(a))
	for i := range result {
		result[i] = make([]int, len(b[0]))
	}
	//var lastField string
	for i := range result {
		for j := range b[0] {
			var r = 0
			for x := range b {
				r += a[i][x] * b[x][j]
			}
			result[i][j] = r
			//lastField = fmt.Sprintf("%d,%d", i, j)
			data += fmt.Sprintf("%d,", r)
		}
		if len(data) > 1024*50 {
			duration, err := p.checkpoint.Nvram.Write([]checkpoint.Data{
				{
					Key:   "result" + fmt.Sprintf("000000%d", checkPointCounter)[len(strconv.Itoa(checkPointCounter)):],
					Value: data[:len(data)-1],
				},
				//{
				//	Key:   "lastField",
				//	Value: lastField,
				//},
			})
			durations = append(durations, duration)
			if err != nil {
				fmt.Println(err.Error())
				return nil, checkPointCounter, err
			}
			data = ""
			checkPointCounter++
		}
	}
	fmt.Println(durations)
	return result, checkPointCounter, nil
}
