package benchmark

import (
	"encoding/json"
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
	key   = "test"
	value = "1234"
)

func NewBenchmark(nvram *nvram.Nvram) error {
	jsonTimes, err := jsonFile()
	if err != nil {
		return err
	}

	db := rocksdb.OpenDb()

	key := "test"
	value := "1234"

	var nvramTimes times
	var rocksDbTimes times
	for i := 0; i < 100; i++ {
		currentKey := fmt.Sprintf("%s%s", key, strconv.Itoa(i))

		startTime := time.Now()
		err := nvram.Set(currentKey, value)
		if err != nil {
			return err
		}
		duration := time.Since(startTime)
		nvramTimes.Put = append(nvramTimes.Put, duration)

		startTime = time.Now()
		err = db.Put(currentKey, value)
		if err != nil {
			return err
		}
		duration = time.Since(startTime)
		rocksDbTimes.Put = append(rocksDbTimes.Put, duration)
	}
	for i := 0; i < 100; i++ {
		currentKey := fmt.Sprintf("%s%s", key, strconv.Itoa(i))

		startTime := time.Now()
		result, err := nvram.Get(currentKey)
		if err != nil || result != value {
			return err
		}
		duration := time.Since(startTime)
		nvramTimes.Get = append(nvramTimes.Get, duration)

		startTime = time.Now()
		result, err = db.Get(currentKey)
		if err != nil || result != value {
			return err
		}
		duration = time.Since(startTime)
		rocksDbTimes.Get = append(rocksDbTimes.Get, duration)
	}
	for i := 0; i < 100; i++ {
		currentKey := fmt.Sprintf("%s%s", key, strconv.Itoa(i))

		startTime := time.Now()
		err := nvram.Delete(currentKey)
		if err != nil {
			return err
		}
		duration := time.Since(startTime)
		nvramTimes.Delete = append(nvramTimes.Delete, duration)

		startTime = time.Now()
		err = db.Delete(currentKey)
		if err != nil {
			return err
		}
		duration = time.Since(startTime)
		rocksDbTimes.Delete = append(rocksDbTimes.Delete, duration)
	}

	// NVRAM
	var nvramGet int64
	for _, duration := range nvramTimes.Get {
		nvramGet += duration.Microseconds()
	}
	fmt.Printf("NVRAM avg GET: %d μs\n", nvramGet/100)
	var nvramPut int64
	for _, duration := range nvramTimes.Put {
		nvramPut += duration.Microseconds()
	}
	fmt.Printf("NVRAM avg PUT: %d μs\n", nvramPut/100)
	var nvramDelete int64
	for _, duration := range nvramTimes.Delete {
		nvramDelete += duration.Microseconds()
	}
	fmt.Printf("NVRAM avg DELETE: %d μs\n", nvramDelete/100)

	// RocksDB
	var rocksDbGet int64
	for _, duration := range rocksDbTimes.Get {
		rocksDbGet += duration.Microseconds()
	}
	fmt.Printf("RocksDB avg GET: %d μs\n", rocksDbGet/100)
	var rocksDbPut int64
	for _, duration := range rocksDbTimes.Put {
		rocksDbPut += duration.Microseconds()
	}
	fmt.Printf("RocksDB avg PUT: %d μs\n", rocksDbPut/100)
	var rocksDbDelete int64
	for _, duration := range rocksDbTimes.Delete {
		rocksDbDelete += duration.Microseconds()
	}
	fmt.Printf("RocksDB avg DELETE: %d μs\n", rocksDbDelete/100)

	// JSON
	var jsonGet int64
	for _, duration := range jsonTimes.Get {
		jsonGet += duration.Microseconds()
	}
	fmt.Printf("JSON avg GET: %d μs\n", jsonGet/100)
	var jsonPut int64
	for _, duration := range jsonTimes.Put {
		jsonPut += duration.Microseconds()
	}
	fmt.Printf("JSON avg PUT: %d μs\n", jsonPut/100)
	var jsonDelete int64
	for _, duration := range jsonTimes.Delete {
		jsonDelete += duration.Microseconds()
	}
	fmt.Printf("JSON avg DELETE: %d μs\n", jsonDelete/100)

	return nil
}

func jsonFile() (*times, error) {
	times := times{
		Delete: make([]time.Duration, 0),
		Get:    make([]time.Duration, 0),
		Put:    make([]time.Duration, 0),
	}

	f, err := os.Create("out/data.json")
	if err != nil {
		return nil, err
	}
	_, err = f.Write([]byte("{}"))
	if err != nil {
		return nil, err
	}

	for i := 0; i < 100; i++ {
		currentKey := fmt.Sprintf("%s%s", key, strconv.Itoa(i))

		startTime := time.Now()
		_, err := f.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}
		var test map[string]string
		content, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(content, &test)
		if err != nil {
			return nil, err
		}
		test[currentKey] = value
		marshal, err := json.Marshal(test)
		if err != nil {
			return nil, err
		}
		_, err = f.WriteAt(marshal, 0)
		if err != nil {
			return nil, err
		}
		duration := time.Since(startTime)
		times.Put = append(times.Put, duration)
	}

	for i := 0; i < 100; i++ {
		currentKey := fmt.Sprintf("%s%s", key, strconv.Itoa(i))

		startTime := time.Now()
		_, err := f.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}
		content, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}
		var test map[string]string
		err = json.Unmarshal(content, &test)
		if err != nil || test[currentKey] != value {
			return nil, err
		}
		duration := time.Since(startTime)
		times.Get = append(times.Get, duration)
	}

	for i := 0; i < 100; i++ {
		currentKey := fmt.Sprintf("%s%s", key, strconv.Itoa(i))

		startTime := time.Now()
		_, err := f.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}
		content, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}
		var test map[string]string
		err = json.Unmarshal(content, &test)
		if err != nil {
			return nil, err
		}
		delete(test, currentKey)
		marshal, err := json.Marshal(test)
		if err != nil {
			return nil, err
		}
		err = f.Truncate(0)
		if err != nil {
			return nil, err
		}
		_, err = f.WriteAt(marshal, 0)
		if err != nil {
			return nil, err
		}
		duration := time.Since(startTime)
		times.Delete = append(times.Delete, duration)
	}
	return &times, nil
}
