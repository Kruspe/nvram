package tester

import (
	"bufio"
	"fmt"
	"github.com/kruspe/nvram/nvram"
	"os"
	"strconv"
)

func CheckSize(nvram *nvram.Nvram, bytes int) ([]string, error) {
	if bytes%100 != 0 {
		return nil, fmt.Errorf("bytes must be a multiple of 100")
	}

	value := "a"
	for i := 0; i < (100-49)/4; i++ {
		value += "aaa"
	}

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
		for i := 0; i < bytes/100; i++ {
			key := fmt.Sprintf("000000%d", counter)[len(strconv.Itoa(counter)):]
			err := nvram.Set(key, value)
			if err != nil {
				return nil, err
			}
			keys = append(keys, key)
			counter++
		}
		fmt.Printf("Added a total of %d bytes\n", len(keys)*100)
	}
}
