package tester

import (
	"bufio"
	"fmt"
	"github.com/kruspe/nvram/nvram"
	"os"
	"strconv"
	"strings"
)

func CheckSize(nvram *nvram.Nvram, bytes int) ([]string, error) {
	if bytes%100 != 0 {
		return nil, fmt.Errorf("bytes must be a multiple of 100")
	}

	value := "a"
	for i := 0; i < (100-49)/4; i++ {
		value += "aaa"
	}
	value += strings.Repeat("a", 49900)

	counter := 0
	var keys []string
	for {
		wait := true
		for wait {
			fmt.Println("Type 'c' to add a value or 'q' to quit.")
			input, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				return nil, err
			}
			if input == "q\n" {
				return keys, nil
			}
			if input == "c\n" {
				wait = false
			}
		}

		fmt.Printf("Adding %d bytes to NVRAM\n", bytes)
		for i := 0; i < bytes/50000; i++ {
			key := fmt.Sprintf("000000%d", counter)[len(strconv.Itoa(counter)):]
			err := nvram.Set(key, value)
			if err != nil {
				return nil, err
			}
			keys = append(keys, key)
			counter++
		}
		fmt.Printf("Added a total of %d bytes\n", len(keys)*50000)
	}
}

func StoreLargeValue(nvram *nvram.Nvram) error {
	// 60645 size
	value := "a" + strings.Repeat("aaa", 1374987)
	bytesToAdd := strings.Repeat("a", 100)

	var key string
	counter := 0
	for {
		wait := true
		for wait {
			fmt.Println("Type 'c' to add next size or 'q' to quit.")
			input, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				return err
			}
			if input == "q\n" {
				return nvram.Delete(key)
			}
			if input == "c\n" {
				wait = false
				err := nvram.Delete(key)
				if err != nil {
					return err
				}
			}
		}

		key = fmt.Sprintf("000000%d", counter)[len(strconv.Itoa(counter)):]
		fmt.Println(len(value))
		err := nvram.Set(key, value)
		if err != nil {
			return err
		}
		counter++
		value += bytesToAdd
	}
}
