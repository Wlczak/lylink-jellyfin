package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Port              int
	JellyfinServerUrl string
}

func GetConfig() Config {
	var fileContent []byte
	var config Config

	configFile, err := os.OpenFile("config.json", os.O_CREATE, 0644)

	if err != nil {
		panic(err)
	}

	_, err = configFile.Read(fileContent)
	if err != nil {
		panic(err)
	}

	if len(fileContent) == 0 {
		fmt.Println("default")
		defaultConfig, err := json.Marshal(Config{
			Port:              0,
			JellyfinServerUrl: "",
		})
		if err != nil {
			panic(err)
		}

		fmt.Println(string(defaultConfig))

		configFile.WriteString(string(defaultConfig))
		configFile.Close()

		fileContent = defaultConfig
	}

	err = json.Unmarshal(fileContent, &config)

	if err != nil {
		fmt.Println(fileContent)
		panic(err)
	}

	return config
}
