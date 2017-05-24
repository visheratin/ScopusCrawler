package config

import (
	"bufio"
	"math/rand"
	"os"
)

var (
	keys []string
)

func InitKeys(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	keys = lines
	return scanner.Err()
}

func GetKey() string {
	index := rand.Intn(len(keys))
	return keys[index]
}

func RemoveKey(key string) {
	for i := 0; i < len(keys); i++ {
		if keys[i] == key {
			keys = append(keys[:i], keys[i+1:]...)
			return
		}
	}
}
