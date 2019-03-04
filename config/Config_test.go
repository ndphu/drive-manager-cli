package config

import (
	"os/user"
	"testing"
)

func TestConfig_Load(t *testing.T) {

}

func TestConfig_Save(t *testing.T) {

}

func TestGetConfigFile(t *testing.T) {
	usr, _ := user.Current()
	if GetConfigFile() != usr.HomeDir + "/" + configFileName {
		t.Errorf("Config File Mismatched")
	}
}

func TestLoadConfig(t *testing.T) {

}

func TestSaveConfig(t *testing.T) {

}