package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os/user"
)

type Config struct {
	BaseUrl string `json:"base_url"`
	Token   string `json:"token"`
}

var configFile = ""
var configFileName = ".drive-manager.conf"

func init() {
	usr, err := user.Current()
	configFile = usr.HomeDir + "/" + configFileName
	_, err = ioutil.ReadFile(configFile)
	if err != nil {
		writeDefaultConfig(configFile)
	}
}

func GetConfigFile() string {
	return configFile
}

func writeDefaultConfig(configFile string) ([]byte, error) {
	config := Config{}
	config.BaseUrl = "http://localhost:8889/api"
	config.Token = ""

	confRaw, err := json.Marshal(&config)
	if err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(configFile, confRaw, 0755); err != nil {
		return nil, err
	}

	return confRaw, nil
}

func (c *Config) Save(file string) error {
	confRaw, err := json.Marshal(&c)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(file, confRaw, 0755); err != nil {
		return err
	}
	return nil
}

func (c *Config) Load(file string) error {
	confRaw, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(confRaw, &c); err != nil {
		log.Println("fail to load config")
		return err
	}
	return nil
}

func LoadConfig() (*Config, error) {
	c := &Config{}
	if err := c.Load(GetConfigFile()); err != nil {
		return nil, err
	}
	return c, nil
}

func SaveConfig(c *Config) error {
	return c.Save(GetConfigFile())
}