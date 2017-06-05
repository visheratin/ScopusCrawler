package config

import (
	"bufio"
	"math/rand"
	"os"
)

func (config *Configuration) InitKeys(path string) error {
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
	config.keys = lines
	return scanner.Err()
}

func (config *Configuration) GetKey() string {
	index := rand.Intn(len(config.keys))
	return config.keys[index]
}

func (config *Configuration) RemoveKey(key string) {
	for i := 0; i < len(config.keys); i++ {
		if config.keys[i] == key {
			config.keys = append(config.keys[:i], config.keys[i+1:]...)
			return
		}
	}
}
