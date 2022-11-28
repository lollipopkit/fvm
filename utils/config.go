package utils

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/model"
)

var (
	Config      model.Config
	ErrNoConfig = errors.New("No config file")
	configPath  string
)

func init() {
	go func() {
		configPath = path.Join(FvmHome, consts.ConfigFileName)
		GetConfig()
	}()
}

func GetConfig() error {
	if Exists(configPath) {
		data, err := os.ReadFile(configPath)
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
		err = os.WriteFile(configPath, data, 0644)
	}
	return err
}
