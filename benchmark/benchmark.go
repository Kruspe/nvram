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
	repeats = 10000
)

func NewBenchmark(nvram *nvram.Nvram) error {
	jsonTimes, err := jsonFileBenchmark()
	if err != nil {
		return err
	}

	gobTimes, err := gobBenchmark()
	if err != nil {
		return err
	}

	rocksDbTimes, err := rocksDbBenchmark()
	if err != nil {
		return err
	}

	nvramTimes, err := nvramBenchmark(nvram)
	if err != nil {
		return err
	}

	// NVRAM
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
	fmt.Printf("NVRAM avg DELETE: %d μs\n", nvramDelete/repeats)

	// RocksDB
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
	fmt.Printf("RocksDB avg DELETE: %d μs\n", rocksDbDelete/repeats)

	// JSON
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
	fmt.Printf("JSON avg DELETE: %d μs\n", jsonDelete/repeats)

	// gob
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
	fmt.Printf("gob avg DELETE: %d μs\n", gobDelete/repeats)

	return nil
}

func nvramBenchmark(nvram *nvram.Nvram) (*times, error) {
	times := times{
		Delete: make([]time.Duration, 0),
		Get:    make([]time.Duration, 0),
		Put:    make([]time.Duration, 0),
	}

	for i := 0; i < repeats; i++ {
		currentKey := fmt.Sprintf("%s%s", key, strconv.Itoa(i))

		startTime := time.Now()
		err := nvram.Set(currentKey, value)
		if err != nil {
			return nil, err
		}
		duration := time.Since(startTime)
		times.Put = append(times.Put, duration)
	}
	for i := 0; i < repeats; i++ {
		currentKey := fmt.Sprintf("%s%s", key, strconv.Itoa(i))

		startTime := time.Now()
		result, err := nvram.Get(currentKey)
		if err != nil || result != value {
			return nil, err
		}
		duration := time.Since(startTime)
		times.Get = append(times.Get, duration)
	}
	for i := 0; i < repeats; i++ {
		currentKey := fmt.Sprintf("%s%s", key, strconv.Itoa(i))

		startTime := time.Now()
		err := nvram.Delete(currentKey)
		if err != nil {
			return nil, err
		}
		duration := time.Since(startTime)
		times.Delete = append(times.Delete, duration)
	}

	return &times, nil
}

func rocksDbBenchmark() (*times, error) {
	times := times{
		Delete: make([]time.Duration, 0),
		Get:    make([]time.Duration, 0),
		Put:    make([]time.Duration, 0),
	}
	db := rocksdb.OpenDb()
	defer db.DeleteDb()

	for i := 0; i < repeats; i++ {
		currentKey := fmt.Sprintf("%s%s", key, strconv.Itoa(i))

		startTime := time.Now()
		err := db.Put(currentKey, value)
		if err != nil {
			return nil, err
		}
		duration := time.Since(startTime)
		times.Put = append(times.Put, duration)
	}
	for i := 0; i < repeats; i++ {
		currentKey := fmt.Sprintf("%s%s", key, strconv.Itoa(i))

		startTime := time.Now()
		result, err := db.Get(currentKey)
		if err != nil || result != value {
			return nil, err
		}
		duration := time.Since(startTime)
		times.Get = append(times.Get, duration)
	}
	for i := 0; i < repeats; i++ {
		currentKey := fmt.Sprintf("%s%s", key, strconv.Itoa(i))

		startTime := time.Now()
		err := db.Delete(currentKey)
		if err != nil {
			return nil, err
		}
		duration := time.Since(startTime)
		times.Delete = append(times.Delete, duration)
	}

	return &times, nil
}

func jsonFileBenchmark() (*times, error) {
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

	for i := 0; i < repeats; i++ {
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
	for i := 0; i < repeats; i++ {
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
	for i := 0; i < repeats; i++ {
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

func gobBenchmark() (*times, error) {
	times := times{
		Delete: make([]time.Duration, 0),
		Get:    make([]time.Duration, 0),
		Put:    make([]time.Duration, 0),
	}

	f, err := os.Create("out/data.gob")
	if err != nil {
		return nil, err
	}
	for i := 0; i < repeats; i++ {
		currentKey := fmt.Sprintf("%s%s", key, strconv.Itoa(i))

		startTime := time.Now()
		var test map[string]string
		_, err := f.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}
		stat, err := f.Stat()
		if err != nil {
			return nil, err
		}
		if stat.Size() == 0 {
			test = make(map[string]string)
		} else {
			err = gob.NewDecoder(f).Decode(&test)
			if err != nil && err != io.EOF {
				return nil, err
			}
		}
		test[currentKey] = value
		_, err = f.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}
		err = gob.NewEncoder(f).Encode(test)
		if err != nil {
			return nil, err
		}
		duration := time.Since(startTime)
		times.Put = append(times.Put, duration)
	}
	for i := 0; i < repeats; i++ {
		currentKey := fmt.Sprintf("%s%s", key, strconv.Itoa(i))

		startTime := time.Now()
		_, err := f.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}
		stat, err := f.Stat()
		if err != nil {
			return nil, err
		}
		var test map[string]string
		if stat.Size() == 0 {
			return nil, errors.New("item does not exist")
		} else {
			err = gob.NewDecoder(f).Decode(&test)
			if err != nil {
				return nil, err
			}
		}
		if test[currentKey] != value {
			return nil, errors.New("item does not exist")
		}
		duration := time.Since(startTime)
		times.Get = append(times.Get, duration)
	}
	for i := 0; i < repeats; i++ {
		currentKey := fmt.Sprintf("%s%s", key, strconv.Itoa(i))

		startTime := time.Now()
		_, err := f.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}
		stat, err := f.Stat()
		if err != nil {
			return nil, err
		}
		var test map[string]string
		if stat.Size() == 0 {
			test = make(map[string]string)
		} else {
			err = gob.NewDecoder(f).Decode(&test)
			if err != nil && err != io.EOF {
				return nil, err
			}
		}
		delete(test, currentKey)
		_, err = f.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}
		err = f.Truncate(0)
		if err != nil {
			return nil, err
		}
		err = gob.NewEncoder(f).Encode(test)
		if err != nil {
			return nil, err
		}
		duration := time.Since(startTime)
		times.Delete = append(times.Delete, duration)
	}

	return &times, err
}
