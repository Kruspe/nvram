package benchmark

import (
	"fmt"
	"github.com/kruspe/nvram/nvram"
	"github.com/kruspe/nvram/problem"
	"strconv"
	"strings"
	"time"
)

func CombinedBenchmark(n *nvram.Nvram) error {
	value := strings.Repeat("a", 50000)
	p, err := problem.NewProblem(n)
	if err != nil {
		return err
	}
	for _, i := range []int{1024, 2048, 3072} {
		a, b := problem.CreateMatrices(i, i)

		//err := benchmarkNvram(n, value)
		//if err != nil {
		//	return err
		//}

		_, lastCheckpointNumber, err := p.MultiplyMeasureNvramSet(a, b)
		if err != nil {
			return err
		}
		err = deleteNvramEntries(n, lastCheckpointNumber)
		if err != nil {
			return err
		}

		err = benchmarkNvram(n, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteNvramEntries(n *nvram.Nvram, lastCheckpointNumber int) error {
	for i := 0; i < lastCheckpointNumber; i++ {
		err := n.Delete("result" + fmt.Sprintf("000000%d", i)[len(strconv.Itoa(i)):])
		if err != nil {
			return err
		}
	}
	return nil
}

func benchmarkNvram(n *nvram.Nvram, value string) error {
	read, set, d := int64(0), int64(0), int64(0)
	for i := 0; i < 200; i++ {
		now := time.Now()
		err := n.Set("test", value)
		set += time.Since(now).Microseconds()
		if err != nil {
			return err
		}

		now = time.Now()
		_, err = n.Get("test")
		read += time.Since(now).Microseconds()
		if err != nil {
			return err
		}

		now = time.Now()
		err = n.Delete("test")
		d += time.Since(now).Microseconds()
		if err != nil {
			return err
		}
	}
	fmt.Printf("Read: %dμs, Set: %dμs, Delete: %dμs\n", read/200, set/200, d/200)
	return nil
}
