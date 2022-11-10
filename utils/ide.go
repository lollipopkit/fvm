package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/term"
	"github.com/tidwall/gjson"
)

func vscode() error {
	var now gjson.Result
	if Exists(consts.VscodeSettingPath) {
		data, err := ioutil.ReadFile(consts.VscodeSettingPath)
		if err != nil {
			return err
		}

		now = gjson.ParseBytes(data)

	}
	add := gjson.Parse(consts.VSC_CONFIG)
	combined := map[string]any{}
	for k, v := range now.Map() {
		combined[k] = v.Value()
	}
	for k, v := range add.Map() {
		combined[k] = v.Value()
	}

	bytes, err := json.Marshal(combined)
	if err != nil {
		return err
	}

	println()
	term.Info("Content write to " + consts.VscodeSettingPath)
	print(string(bytes))
	write := term.Confirm("\nConfirm?", true)
	if write {
		if !Exists(consts.VscodeDirName) {
			err = os.Mkdir(consts.VscodeDirName, 0755)
			if err != nil {
				return err
			}
		}
		term.Success("Config success!")
		return ioutil.WriteFile(consts.VscodeSettingPath, bytes, 0644)
	}
	return nil
}

func idea() error {
	term.Warn("IDEA is not supported yet")
	return nil
}

func ConfigIde() error {
	options := []string{"VSCode", "IDEA", "skip"}
	idx := term.Option("Which IDE do you want to auto config?", options, len(options)-1)
	switch idx {
	case 0:
		return vscode()
	case 1:
		return idea()
	}
	return nil
}
