package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path"

	"github.com/lollipopkit/fvm/model"
)

var (
	Config      model.Config
	ErrNoConfig = errors.New("No config file")
)

func init() {
	GetConfig()
}

func GetConfig() error {
	configPath := path.Join(FvmHome, "config.json")
	if Exists(configPath) {
		data, err := ioutil.ReadFile(configPath)
		if err == nil {
			err = json.Unmarshal(data, &Config)
		}
		return err
	}
	return ErrNoConfig
}

func SaveConfig() error {
	configPath := path.Join(FvmHome, "config.json")
	data, err := json.Marshal(Config)
	if err == nil {
		err = ioutil.WriteFile(configPath, data, 0644)
	}
	return err
}
