package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/term"
	"github.com/tidwall/gjson"
)

func vscode() error {
	data, err := os.ReadFile(consts.VscodeSettingPath)
	if err != nil {
		if !Exists(consts.VscodeSettingPath) {
			data = []byte("{}")
		} else {
			return err
		}
	}

	now := gjson.ParseBytes(data)
	add := gjson.Parse(consts.VSC_CONFIG)
	combined := map[string]any{}
	for k, v := range now.Map() {
		combined[k] = v.Value()
	}
	for k, v := range add.Map() {
		combined[k] = v.Value()
	}

	_bytes, err := json.MarshalIndent(combined, "", "  ")
	if err != nil {
		return err
	}

	if bytes.Equal(_bytes, data) {
		return nil
	}

	println()
	print(string(_bytes))

	write := term.Confirm(fmt.Sprintf("\nWrite above content into %s?", consts.VscodeSettingPath), true)
	if write {
		if !Exists(consts.VscodeDirName) {
			err = os.Mkdir(consts.VscodeDirName, 0755)
			if err != nil {
				return err
			}
		}
		if err = os.WriteFile(consts.VscodeSettingPath, _bytes, 0644); err != nil {
			return err
		}
		term.Success("Configured VSCode.")
	}
	return nil
}

func idea() error {
	term.Warn("IDEA is not supported yet.")
	return nil
}

func ConfigIde() error {
	if os.Getenv("VSCODE_INJECTION") != "" {
		return vscode()
	}
	err := vscode()
	if err != nil {
		return err
	}

	return idea()
}
