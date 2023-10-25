package benchmark

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kruspe/nvram/nvram"
	"github.com/kruspe/nvram/rocksdb"
	"io"
	"os"
	"strconv"
	"time"
)

type times struct {
	Delete []time.Duration
	Get    []time.Duration
	Put    []time.Duration
}

const (
	key     = "test"
	value   = "1234"
	repeats = 200
)

func NewBenchmark(nvram *nvram.Nvram) error {
	output := make(map[int]map[string]averageTimes)
	for _, entries := range []int{1, 100, 5000, 10000} {
		fmt.Printf("\n--- Benchmark with %d entries ---\n", entries)
		m := prepareMap(entries)

		//JSON
		jsonTimes, err := jsonFileBenchmark(m)
		if err != nil {
			return err
		}
		var jsonGet int64
		for _, duration := range jsonTimes.Get {
			jsonGet += duration.Microseconds()
		}
		fmt.Printf("JSON avg GET: %d μs\n", jsonGet/repeats)
		var jsonPut int64
		for _, duration := range jsonTimes.Put {
			jsonPut += duration.Microseconds()
		}
		fmt.Printf("JSON avg PUT: %d μs\n", jsonPut/repeats)
		var jsonDelete int64
		for _, duration := range jsonTimes.Delete {
			jsonDelete += duration.Microseconds()
		}
		fmt.Printf("JSON avg DELETE: %d μs\n\n", jsonDelete/repeats)

		//gob
		gobTimes, err := gobBenchmark(m)
		if err != nil {
			return err
		}
		var gobGet int64
		for _, duration := range gobTimes.Get {
			gobGet += duration.Microseconds()
		}
		fmt.Printf("gob avg GET: %d μs\n", gobGet/repeats)
		var gobPut int64
		for _, duration := range gobTimes.Put {
			gobPut += duration.Microseconds()
		}
		fmt.Printf("gob avg PUT: %d μs\n", gobPut/repeats)
		var gobDelete int64
		for _, duration := range gobTimes.Delete {
			gobDelete += duration.Microseconds()
		}
		fmt.Printf("gob avg DELETE: %d μs\n\n", gobDelete/repeats)

		//RocksDB
		rocksDbTimes, err := rocksDbBenchmark(m)
		if err != nil {
			return err
		}
		var rocksDbGet int64
		for _, duration := range rocksDbTimes.Get {
			rocksDbGet += duration.Microseconds()
		}
		fmt.Printf("RocksDB avg GET: %d μs\n", rocksDbGet/repeats)
		var rocksDbPut int64
		for _, duration := range rocksDbTimes.Put {
			rocksDbPut += duration.Microseconds()
		}
		fmt.Printf("RocksDB avg PUT: %d μs\n", rocksDbPut/repeats)
		var rocksDbDelete int64
		for _, duration := range rocksDbTimes.Delete {
			rocksDbDelete += duration.Microseconds()
		}
		fmt.Printf("RocksDB avg DELETE: %d μs\n\n", rocksDbDelete/repeats)

		//NVRAM
		nvramTimes, err := nvramBenchmark(nvram, m)
		if err != nil {
			return err
		}
		var nvramGet int64
		for _, duration := range nvramTimes.Get {
			nvramGet += duration.Microseconds()
		}
		fmt.Printf("NVRAM avg GET: %d μs\n", nvramGet/repeats)
		var nvramPut int64
		for _, duration := range nvramTimes.Put {
			nvramPut += duration.Microseconds()
		}
		fmt.Printf("NVRAM avg PUT: %d μs\n", nvramPut/repeats)
		var nvramDelete int64
		for _, duration := range nvramTimes.Delete {
			nvramDelete += duration.Microseconds()
		}
		fmt.Printf("NVRAM avg DELETE: %d μs\n\n", nvramDelete/repeats)

		output[entries] = map[string]averageTimes{
			"\\glsentryshort{json}": {
				get:    jsonGet / repeats,
				put:    jsonPut / repeats,
				delete: jsonDelete / repeats,
			},
			"gob": {
				get:    gobGet / repeats,
				put:    gobPut / repeats,
				delete: gobDelete / repeats,
			},
			"RocksDB": {
				get:    rocksDbGet / repeats,
				put:    rocksDbPut / repeats,
				delete: rocksDbDelete / repeats,
			},
			"\\glsentryshort{nvram}": {
				get:    nvramGet / repeats,
				put:    nvramPut / repeats,
				delete: nvramDelete / repeats,
			},
		}
	}

	outputForLatex(output)

	return nil
}

type averageTimes struct {
	get    int64
	put    int64
	delete int64
}

func outputForLatex(input map[int]map[string]averageTimes) {
	fmt.Println("\\begin{table}[H]")
	for index, t := range []string{"\\glsentryshort{json}", "gob", "RocksDB", "\\glsentryshort{nvram}"} {
		fmt.Println("    \\begin{subtable}{.9\\linewidth}")
		fmt.Println("        \\centering")
		fmt.Println("        \\begin{tabular}{| c || c | c | c |}")
		fmt.Println("        \\hline")
		fmt.Printf("        \\textbf{%s} & get & write & delete \\\\\n", t)
		fmt.Println("        \\hline")
		for _, i := range []int{1, 100, 5000, 10000} {
			fmt.Printf("        %d entries & %d$\\mu$s & %d$\\mu$s & %d$\\mu$s \\\\\n", i, input[i][t].get, input[i][t].put, input[i][t].delete)
			fmt.Println("        \\hline")
		}
		fmt.Println("        \\end{tabular}")
		var caption string
		if index == 1 || index == 2 {
			caption = t
		}
		if index == 0 {
			caption = "\\gls{json}"
		}
		if index == 3 {
			caption = "\\gls{nvram}"
		}
		fmt.Printf("        \\caption{%s operation speed}\n", caption)

		var label string
		if index == 0 {
			label = "json"
		}
		if index == 1 {
			label = "gob"
		}
		if index == 2 {
			label = "rocksDB"
		}
		if index == 3 {
			label = "nvram"
		}
		fmt.Printf("        \\label{tbl:speed-%s}\n", label)
		fmt.Println("    \\end{subtable}")
		if index == 3 {
			break
		}
		fmt.Println("    \\hfill")
		fmt.Println("    \\newline")
		fmt.Println("    \\hfill")
		fmt.Println("    \\newline")
	}
	fmt.Println("\\end{table}")
}

func prepareMap(entries int) map[string]string {
	m := make(map[string]string)
	for i := 0; i < entries; i++ {
		key := fmt.Sprintf("static%d", i)
		m[key] = value
	}
	return m
}

func nvramBenchmark(nvram *nvram.Nvram, startMap map[string]string) (*times, error) {
	times := times{
		Delete: make([]time.Duration, 0),
		Get:    make([]time.Duration, 0),
		Put:    make([]time.Duration, 0),
	}
	for k, v := range startMap {
		err := nvram.Set(k, v)
		if err != nil {
			return nil, err
		}
	}
	defer func() {
		for k, _ := range startMap {
			err := nvram.Delete(k)
			if err != nil {
				return
			}
		}
	}()

	for i := 0; i < repeats; i++ {
		currentKey := fmt.Sprintf("%s%s", key, strconv.Itoa(i))

		startTime := time.Now()
		err := nvram.Set(currentKey, value)
		if err != nil {
			return nil, err
		}
		duration := time.Since(startTime)
		times.Put = append(times.Put, duration)

		startTime = time.Now()
		result, err := nvram.Get(currentKey)
		if err != nil || result != value {
			return nil, err
		}
		duration = time.Since(startTime)
		times.Get = append(times.Get, duration)

		startTime = time.Now()
		err = nvram.Delete(currentKey)
		if err != nil {
			return nil, err
		}
		duration = time.Since(startTime)
		times.Delete = append(times.Delete, duration)
	}

	return &times, nil
}

func rocksDbBenchmark(startMap map[string]string) (*times, error) {
	times := times{
		Delete: make([]time.Duration, 0),
		Get:    make([]time.Duration, 0),
		Put:    make([]time.Duration, 0),
	}
	db := rocksdb.OpenDb(false)

	for k, v := range startMap {
		err := db.Put(k, v)
		if err != nil {
			return nil, err
		}
	}
	db.Close()

	syncDb := rocksdb.OpenDb(true)
	defer syncDb.DeleteDb()
	for i := 0; i < repeats; i++ {
		currentKey := fmt.Sprintf("%s%s", key, strconv.Itoa(i))

		startTime := time.Now()
		err := syncDb.Put(currentKey, value)
		if err != nil {
			return nil, err
		}
		duration := time.Since(startTime)
		times.Put = append(times.Put, duration)

		startTime = time.Now()
		result, err := syncDb.Get(currentKey)
		if err != nil || result != value {
			return nil, err
		}
		duration = time.Since(startTime)
		times.Get = append(times.Get, duration)

		startTime = time.Now()
		err = syncDb.Delete(currentKey)
		if err != nil {
			return nil, err
		}
		duration = time.Since(startTime)
		times.Delete = append(times.Delete, duration)
	}

	return &times, nil
}

func jsonFileBenchmark(startMap map[string]string) (*times, error) {
	fileName := "out/data.json"
	times := times{
		Delete: make([]time.Duration, 0),
		Get:    make([]time.Duration, 0),
		Put:    make([]time.Duration, 0),
	}

	initFile, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	err = json.NewEncoder(initFile).Encode(startMap)
	if err != nil {
		return nil, err
	}
	err = initFile.Close()
	if err != nil {
		return nil, err
	}

	for i := 0; i < repeats; i++ {
		currentKey := fmt.Sprintf("%s%s", key, strconv.Itoa(i))

		startTime := time.Now()
		putFile, err := os.OpenFile(fileName, os.O_RDWR, 0755)
		if err != nil {
			return nil, err
		}
		var test map[string]string
		err = json.NewDecoder(putFile).Decode(&test)
		if err != nil {
			return nil, err
		}
		test[currentKey] = value
		err = putFile.Truncate(0)
		if err != nil {
			return nil, err
		}
		_, err = putFile.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}
		err = json.NewEncoder(putFile).Encode(test)
		if err != nil {
			return nil, err
		}
		err = putFile.Close()
		if err != nil {
			return nil, err
		}
		duration := time.Since(startTime)
		times.Put = append(times.Put, duration)

		startTime = time.Now()
		getFile, err := os.OpenFile(fileName, os.O_RDWR, 0755)
		if err != nil {
			return nil, err
		}
		err = json.NewDecoder(getFile).Decode(&test)
		if err != nil || test[currentKey] != value {
			return nil, err
		}
		err = getFile.Close()
		if err != nil {
			return nil, err
		}
		duration = time.Since(startTime)
		times.Get = append(times.Get, duration)

		startTime = time.Now()
		deleteFile, err := os.OpenFile(fileName, os.O_RDWR, 0755)
		if err != nil {
			return nil, err
		}
		err = json.NewDecoder(deleteFile).Decode(&test)
		if err != nil {
			return nil, err
		}
		delete(test, currentKey)
		err = deleteFile.Truncate(0)
		if err != nil {
			return nil, err
		}
		_, err = deleteFile.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}
		err = json.NewEncoder(deleteFile).Encode(test)
		if err != nil {
			return nil, err
		}
		err = deleteFile.Close()
		if err != nil {
			return nil, err
		}
		duration = time.Since(startTime)
		times.Delete = append(times.Delete, duration)
	}

	return &times, nil
}

func gobBenchmark(startMap map[string]string) (*times, error) {
	fileName := "out/data.gob"
	times := times{
		Delete: make([]time.Duration, 0),
		Get:    make([]time.Duration, 0),
		Put:    make([]time.Duration, 0),
	}

	initFile, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	err = gob.NewEncoder(initFile).Encode(startMap)
	if err != nil {
		return nil, err
	}
	err = initFile.Close()
	if err != nil {
		return nil, err
	}

	for i := 0; i < repeats; i++ {
		currentKey := fmt.Sprintf("%s%s", key, strconv.Itoa(i))

		startTime := time.Now()
		putFile, err := os.OpenFile(fileName, os.O_RDWR, 0755)
		if err != nil {
			return nil, err
		}
		var test map[string]string
		err = gob.NewDecoder(putFile).Decode(&test)
		if err != nil {
			return nil, err
		}
		test[currentKey] = value
		err = putFile.Truncate(0)
		if err != nil {
			return nil, err
		}
		_, err = putFile.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}
		err = gob.NewEncoder(putFile).Encode(test)
		if err != nil {
			return nil, err
		}
		err = putFile.Close()
		if err != nil {
			return nil, err
		}
		duration := time.Since(startTime)
		times.Put = append(times.Put, duration)

		startTime = time.Now()
		getFile, err := os.OpenFile(fileName, os.O_RDWR, 0755)
		if err != nil {
			return nil, err
		}
		err = gob.NewDecoder(getFile).Decode(&test)
		if test[currentKey] != value {
			return nil, errors.New("item does not exist")
		}
		err = getFile.Close()
		if err != nil {
			return nil, err
		}
		duration = time.Since(startTime)
		times.Get = append(times.Get, duration)

		startTime = time.Now()
		deleteFile, err := os.OpenFile(fileName, os.O_RDWR, 0755)
		if err != nil {
			return nil, err
		}
		err = gob.NewDecoder(deleteFile).Decode(&test)
		if err != nil {
			return nil, err
		}
		delete(test, currentKey)
		err = deleteFile.Truncate(0)
		if err != nil {
			return nil, err
		}
		_, err = deleteFile.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}
		err = gob.NewEncoder(deleteFile).Encode(test)
		if err != nil {
			return nil, err
		}
		err = deleteFile.Close()
		if err != nil {
			return nil, err
		}
		duration = time.Since(startTime)
		times.Delete = append(times.Delete, duration)
	}

	return &times, err
}
