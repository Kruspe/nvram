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
	"strings"
	"time"
)

type times struct {
	Delete []time.Duration
	Get    []time.Duration
	Put    []time.Duration
}

const (
	key     = "test"
	repeats = 200
)

func NewBenchmark(nvram *nvram.Nvram) error {
	for _, entries := range []int{1, 100, 5000, 10000} {
		m := prepareMap(entries)
		for _, v := range []int{100, 10240, 30720, 51200} {
			fmt.Printf("\n--- Benchmark with %d entries and value size %dKB ---\n", entries, v/1024)

			//JSON
			jsonTimes, err := jsonFileBenchmark(m, v)
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
			gobTimes, err := gobBenchmark(m, v)
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
			rocksDbTimes, err := rocksDbBenchmark(m, v)
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
			nvramTimes, err := nvramBenchmark(nvram, m, v)
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
		}

		//JSON

		//output[entries] = map[string]averageTimes{
		//	"\\glsentryshort{json}": {
		//		get:    jsonGet / repeats,
		//		put:    jsonPut / repeats,
		//		delete: jsonDelete / repeats,
		//	},
		//	"gob": {
		//		get:    gobGet / repeats,
		//		put:    gobPut / repeats,
		//		delete: gobDelete / repeats,
		//	},
		//	"RocksDB": {
		//		get:    rocksDbGet / repeats,
		//		put:    rocksDbPut / repeats,
		//		delete: rocksDbDelete / repeats,
		//	},
		//	"\\glsentryshort{nvram}": {
		//		get:    nvramGet / repeats,
		//		put:    nvramPut / repeats,
		//		delete: nvramDelete / repeats,
		//	},
		//}
	}
	return nil
}

func prepareCustomMap(entries int, bytes int) map[string]string {
	v := "a"
	for i := 0; i < (100-49)/4; i++ {
		v += "aaa"
	}
	for i := 0; i < (bytes-100)/100; i++ {
		v += strings.Repeat("a", 100)
	}

	m := make(map[string]string)
	for i := 0; i < entries; i++ {
		key := fmt.Sprintf("static%d", i)
		m[key] = v
	}
	return m
}

func prepareMap(entries int) map[string]string {
	m := make(map[string]string)
	for i := 0; i < entries; i++ {
		key := fmt.Sprintf("static%d", i)
		m[key] = "1234"
	}
	return m
}

func nvramBenchmark(nvram *nvram.Nvram, startMap map[string]string, valueSize int) (*times, error) {
	times := times{
		Delete: make([]time.Duration, 0),
		Get:    make([]time.Duration, 0),
		Put:    make([]time.Duration, 0),
	}
	value := strings.Repeat("a", valueSize)
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
		duration := time.Since(startTime)
		if err != nil {
			return nil, err
		}
		times.Put = append(times.Put, duration)

		startTime = time.Now()
		_, err = nvram.Get(currentKey)
		duration = time.Since(startTime)
		if err != nil {
			return nil, err
		}
		times.Get = append(times.Get, duration)

		startTime = time.Now()
		err = nvram.Delete(currentKey)
		duration = time.Since(startTime)
		if err != nil {
			return nil, err
		}
		times.Delete = append(times.Delete, duration)
	}

	return &times, nil
}

func rocksDbBenchmark(startMap map[string]string, valueSize int) (*times, error) {
	times := times{
		Delete: make([]time.Duration, 0),
		Get:    make([]time.Duration, 0),
		Put:    make([]time.Duration, 0),
	}
	value := strings.Repeat("a", valueSize)
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
		duration := time.Since(startTime)
		if err != nil {
			return nil, err
		}
		times.Put = append(times.Put, duration)

		startTime = time.Now()
		_, err = syncDb.Get(currentKey)
		duration = time.Since(startTime)
		if err != nil {
			return nil, err
		}
		times.Get = append(times.Get, duration)

		startTime = time.Now()
		err = syncDb.Delete(currentKey)
		duration = time.Since(startTime)
		if err != nil {
			return nil, err
		}
		times.Delete = append(times.Delete, duration)
	}

	return &times, nil
}

func jsonFileBenchmark(startMap map[string]string, valueSize int) (*times, error) {
	fileName := "out/data.json"
	times := times{
		Delete: make([]time.Duration, 0),
		Get:    make([]time.Duration, 0),
		Put:    make([]time.Duration, 0),
	}
	value := strings.Repeat("a", valueSize)

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

func gobBenchmark(startMap map[string]string, valueSize int) (*times, error) {
	fileName := "out/data.gob"
	times := times{
		Delete: make([]time.Duration, 0),
		Get:    make([]time.Duration, 0),
		Put:    make([]time.Duration, 0),
	}
	value := strings.Repeat("a", valueSize)

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
