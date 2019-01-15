package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	BackendUrl string `json:"backend_url"`
}

var (
	config *Config
)

func init() {
	//confFile, err := os.Open("conf.json")
	confRaw, err := ioutil.ReadFile("~/.drive-cli-conf.json")
	if err != nil {
		config = defaultConfig()
	} else {
		err := json.Unmarshal(confRaw, &config)
		if err != nil {
			config = defaultConfig()
		}
	}
}
func defaultConfig() *Config {
	return &Config {
		BackendUrl: "http://localhost:8889/api",
	}
}

func GetConfig() *Config {
	return config
}