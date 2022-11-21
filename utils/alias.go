package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/term"
)

var (
	ErrShellConfigFileNotFound = fmt.Errorf("shell config file not found")

	lines2Add = []string{"alias dart='fvm dart'", "alias flutter='fvm flutter'"}
)

func SetAlias() error {
	shell := GetShell()

	switch shell {
	case ShellTypeZsh:
		return setZshAlias()
	case ShellTypeBash:
		return setBashAlias()
	case ShellTypeFish:
		return setFishAlias()
	}
	return fmt.Errorf(ErrUnsupportedShellPrefix + shell.String())
}

func setZshAlias() error {
	aliasFile := path.Join(os.Getenv("HOME"), consts.ZshRcName)
	if Exists(aliasFile) {
		term.Info("\nConfiguring [%s]...", consts.ZshRcName)
		f, err := os.OpenFile(aliasFile, os.O_APPEND|os.O_RDWR, 0600)
		if err != nil {
			return err
		}

		defer f.Close()

		data, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		lines := strings.Split(string(data), "\n")

		for _, line2Add := range lines2Add {
			if !Contains(lines, line2Add) {
				if _, err = f.WriteString("\n" + line2Add); err != nil {
					return err
				}
			}
		}
		term.Success("Configured [%s]", consts.ZshRcName)
		return nil
	}

	return ErrShellConfigFileNotFound
}

func setBashAlias() error {
	lines2Add := []string{"alias dart='fvm dart'", "alias flutter='fvm flutter'"}
	aliasFile := path.Join(os.Getenv("HOME"), consts.BashRcName)
	if Exists(aliasFile) {
		term.Info("\nConfiguring [%s]...", consts.BashRcName)
		f, err := os.OpenFile(aliasFile, os.O_APPEND|os.O_RDWR, 0600)
		if err != nil {
			return err
		}

		defer f.Close()

		data, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		lines := strings.Split(string(data), "\n")

		for _, line2Add := range lines2Add {
			if !Contains(lines, line2Add) {
				if _, err = f.WriteString("\n" + line2Add); err != nil {
					return err
				}
			}
		}
		term.Success("Configured [%s]", consts.BashRcName)
		return nil
	}

	return ErrShellConfigFileNotFound
}

func setFishAlias() error {
	aliasFile := path.Join(os.Getenv("HOME"), consts.FishConfigPath)
	if Exists(aliasFile) {
		term.Info("\nConfiguring [config.fish]...")
		f, err := os.OpenFile(aliasFile, os.O_APPEND|os.O_RDWR, 0600)
		if err != nil {
			return err
		}

		defer f.Close()

		data, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		lines := strings.Split(string(data), "\n")

		for _, line2Add := range lines2Add {
			if !Contains(lines, line2Add) {
				if _, err = f.WriteString("\n" + line2Add); err != nil {
					return err
				}
			}
		}
		term.Success("Configured [config.fish]")
		return nil
	}

	return ErrShellConfigFileNotFound
}
