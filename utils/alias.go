package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/lollipopkit/fvm/term"
)

var (
	ErrShellConfigFileNotFound = fmt.Errorf("shell config file not found")

	lines2Add = []string{"alias dart='fvm dart'", "alias flutter='fvm flutter'"}
)

func SetAlias() error {
	shell := os.Getenv("SHELL")
	shell = path.Base(shell)

	switch shell {
	case "zsh":
		return setZshAlias()
	case "bash":
		return setBashAlias()
	case "fish":
		return setFishAlias()
	}
	return fmt.Errorf(ErrUnsupportedShellPrefix + shell)
}

func setZshAlias() error {
	aliasFile := path.Join(os.Getenv("HOME"), ".zshrc")
	if Exists(aliasFile) {
		term.Info("\nConfiguring \".zshrc\"...")
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
		term.Success("Configured \".zshrc\"")
		return nil
	}

	return ErrShellConfigFileNotFound
}

func setBashAlias() error {
	lines2Add := []string{"alias dart='fvm dart'", "alias flutter='fvm flutter'"}
	aliasFile := path.Join(os.Getenv("HOME"), ".bashrc")
	if Exists(aliasFile) {
		term.Info("\nConfiguring \".bashrc\"...")
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
		term.Success("Configured \".bashrc\"")
		return nil
	}

	return ErrShellConfigFileNotFound
}

func setFishAlias() error {
	aliasFile := path.Join(os.Getenv("HOME"), ".config/fish/config.fish")
	if Exists(aliasFile) {
		term.Info("\nConfiguring \"config.fish\"...")
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
		term.Success("Configured \"config.fish\"")
		return nil
	}

	return ErrShellConfigFileNotFound
}
