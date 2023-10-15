package benchmark

import (
	"encoding/json"
	"fmt"
	"github.com/kruspe/nvram/nvram"
	"github.com/kruspe/nvram/rocksdb"
	"os"
	"strconv"
	"time"
)

type times struct {
	Delete []time.Duration
	Get    []time.Duration
	Put    []time.Duration
}

func NewBenchmark(nvram *nvram.Nvram) {
	db := rocksdb.OpenDb()

	key := "test"
	value := "1234"

	var nvramTimes times
	var rocksDbTimes times
	var jsonTimes times
	for i := 0; i < 100; i++ {
		currentKey := fmt.Sprintf("%s%s", key, strconv.Itoa(i))

		startTime := time.Now()
		err := nvram.Set(currentKey, value)
		duration := time.Since(startTime)
		if err != nil {
			fmt.Println("nvram PUT", err)
		}
		nvramTimes.Put = append(nvramTimes.Put, duration)

		startTime = time.Now()
		err = db.Put(currentKey, value)
		duration = time.Since(startTime)
		if err != nil {
			fmt.Println("db PUT", err)
		}
		rocksDbTimes.Put = append(rocksDbTimes.Put, duration)

		startTime = time.Now()
		f, _ := os.ReadFile("out/data.json")
		var test map[string]string
		_ = json.Unmarshal(f, &test)
		test[currentKey] = value
		marshal, _ := json.Marshal(test)
		err = os.WriteFile("out/data.json", marshal, 0666)
		duration = time.Since(startTime)
		if err != nil {
			fmt.Println("json PUT", err)
		}
		jsonTimes.Put = append(jsonTimes.Put, duration)
	}
	for i := 0; i < 100; i++ {
		currentKey := fmt.Sprintf("%s%s", key, strconv.Itoa(i))

		startTime := time.Now()
		result, err := nvram.Get(currentKey)
		duration := time.Since(startTime)
		if err != nil || result != value {
			fmt.Println("nvram GET", result, err)
		}
		nvramTimes.Get = append(nvramTimes.Get, duration)

		startTime = time.Now()
		result, err = db.Get(currentKey)
		duration = time.Since(startTime)
		if err != nil || result != value {
			fmt.Println("db GET", result, err)
		}
		rocksDbTimes.Get = append(rocksDbTimes.Get, duration)

		startTime = time.Now()
		f, _ := os.ReadFile("out/data.json")
		var test map[string]string
		err = json.Unmarshal(f, &test)
		duration = time.Since(startTime)
		if err != nil || test[currentKey] != value {
			fmt.Println("json GET", test[currentKey], err)
		}
		jsonTimes.Get = append(jsonTimes.Get, duration)
	}
	for i := 0; i < 100; i++ {
		currentKey := fmt.Sprintf("%s%s", key, strconv.Itoa(i))

		startTime := time.Now()
		err := nvram.Delete(currentKey)
		duration := time.Since(startTime)
		if err != nil {
			fmt.Println("nvram DELETE", err)
		}
		nvramTimes.Delete = append(nvramTimes.Delete, duration)

		startTime = time.Now()
		err = db.Delete(currentKey)
		duration = time.Since(startTime)
		if err != nil {
			fmt.Println("db DELETE", err)
		}
		rocksDbTimes.Delete = append(rocksDbTimes.Delete, duration)

		startTime = time.Now()
		f, _ := os.ReadFile("out/data.json")
		var test map[string]string
		_ = json.Unmarshal(f, &test)
		delete(test, currentKey)
		marshal, _ := json.Marshal(test)
		err = os.WriteFile("out/data.json", marshal, 0666)
		duration = time.Since(startTime)
		if err != nil {
			fmt.Println("json DELETE", err)
		}
		jsonTimes.Delete = append(jsonTimes.Delete, duration)
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
}
