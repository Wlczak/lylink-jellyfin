package config

import (
	"encoding/json"
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
	var err error

	zap := logs.GetLogger()

	_, err = os.Stat("config.json")

	if os.IsNotExist(err) {
		writeDefaultConfig()
	}

	fileContent, err = os.ReadFile("config.json")

	if err != nil {
		zap.Error(err.Error())
		panic(err)
	}

	if len(fileContent) == 0 {
		fileContent = writeDefaultConfig()
	}

	err = json.Unmarshal(fileContent, &config)

	if err != nil {
		zap.Error(err.Error())
		writeDefaultConfig()
		config = Config{}
	}

	return config
}

func writeDefaultConfig() []byte {
	zap := logs.GetLogger()

	defaultConfig, err := json.Marshal(Config{Port: 8040, JellyfinServerUrl: "http://localhost:8096"})
	if err != nil {
		zap.Error(err.Error())
		panic(err)
	}

	configFile, err := os.Create("config.json")
	if err != nil {
		zap.Error(err.Error())
		panic(err)
	}

	err = configFile.Sync()
	if err != nil {
		zap.Error(err.Error())
		panic(err)
	}

	_, err = configFile.WriteString(string(defaultConfig))
	if err != nil {
		zap.Error(err.Error())
		panic(err)
	}

	err = configFile.Sync()
	if err != nil {
		zap.Error(err.Error())
		panic(err)
	}

	err = configFile.Close()
	if err != nil {
		zap.Error(err.Error())
		panic(err)
	}

	return defaultConfig
}

func (c Config) Save() error {
	zap := logs.GetLogger()

	_, err := os.Stat("config.json")

	if os.IsNotExist(err) {
		writeDefaultConfig()
	}

	configFile, err := os.OpenFile("config.json", os.O_RDWR, 0644)

	if err != nil {
		zap.Error(err.Error())
		return err
	}

	err = configFile.Truncate(0)
	if err != nil {
		zap.Error(err.Error())
		return err
	}

	_, err = configFile.Seek(0, 0)
	if err != nil {
		zap.Error(err.Error())
		return err
	}

	fileContent, err := json.Marshal(c)
	if err != nil {
		zap.Error(err.Error())
		return err
	}

	_, err = configFile.Write(fileContent)
	if err != nil {
		zap.Error(err.Error())
		return err
	}
	return nil
}
