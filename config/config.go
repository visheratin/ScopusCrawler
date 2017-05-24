package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	ListenPort string
	LogPath    string
}

var (
	defaultConfig = Configuration{
		ListenPort: "12345"}
	defaultConfigPath = "config.json"
)

func ReadConfig(configPath string) Configuration {
	if configPath == "" {
		configPath = defaultConfigPath
	}
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Println("Unable to open configuration file.")
		return defaultConfig
	}
	decoder := json.NewDecoder(file)
	var config Configuration
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println(err)
		return defaultConfig
	}
	return config
}
