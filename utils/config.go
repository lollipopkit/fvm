package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path"

	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/model"
)

var (
	Config      model.Config
	ErrNoConfig = errors.New("No config file")
	configPath  = path.Join(FvmHome, consts.ConfigFileName)
)

func init() {
	GetConfig()
}

func GetConfig() error {
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
	data, err := json.Marshal(Config)
	if err == nil {
		err = ioutil.WriteFile(configPath, data, 0644)
	}
	return err
}
