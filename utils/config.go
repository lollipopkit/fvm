package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path"

	"github.com/lollipopkit/gofvm/model"
)

var (
	config      model.Config
	ErrNoConfig = errors.New("No config file")
)

func init() {
	err := GetConfig()
	if err == nil {
		IsInChina = &config.InChina
	}
}

func GetConfig() error {
	configPath := path.Join(Path(), "config.json")
	if Exists(configPath) {
		data, err := ioutil.ReadFile(configPath)
		if err == nil {
			err = json.Unmarshal(data, &config)
		}
		return err
	}
	return ErrNoConfig
}

func SaveConfig() error {
	configPath := path.Join(Path(), "config.json")
	data, err := json.Marshal(config)
	if err == nil {
		err = ioutil.WriteFile(configPath, data, 0644)
	}
	return err
}
