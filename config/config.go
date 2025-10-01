package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Wlczak/lylink-jellyfin/logs"
)

type Config struct {
	Port              int
	JellyfinServerUrl string
}

func GetConfig() Config {
	var fileContent []byte
	var config Config
	zap := logs.GetLogger()

	configFile, err := os.OpenFile("config.json", os.O_CREATE, 0644)

	if err != nil {
		panic(err)
	}

	_, err = configFile.Read(fileContent)
	if err != nil {
		zap.Error(err.Error())
		panic(err)
	}

	if len(fileContent) == 0 {
		fmt.Println("default")
		defaultConfig, err := json.Marshal(Config{
			Port:              0,
			JellyfinServerUrl: "",
		})
		if err != nil {
			zap.Error(err.Error())
			panic(err)
		}

		fmt.Println(string(defaultConfig))

		_, err = configFile.WriteString(string(defaultConfig))
		if err != nil {
			zap.Error(err.Error())
			panic(err)
		}

		err = configFile.Close()
		if err != nil {
			zap.Error(err.Error())
			panic(err)
		}

		fileContent = defaultConfig
	}

	err = json.Unmarshal(fileContent, &config)

	if err != nil {
		zap.Error(err.Error())
		panic(err)
	}

	return config
}
