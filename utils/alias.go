package utils

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/lollipopkit/fvm/term"
)

var (
	lines2Add = []string{"alias dart='fvm dart'", "alias flutter='fvm flutter'"}
)

type aliasConfiger interface {
	SetAlias() error
}

type zshAliasConfiger struct{}

func (zshAliasConfiger) SetAlias() error {
	if Exists(RcPath) {
		f, err := os.OpenFile(RcPath, os.O_APPEND|os.O_RDWR, 0600)
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
		term.Success("Configured %s", RcPath)
		return nil
	}

	return ErrShellConfigNotFound
}

type bashAliasConfiger struct{}

func (bashAliasConfiger) SetAlias() error {
	if Exists(RcPath) {
		f, err := os.OpenFile(RcPath, os.O_APPEND|os.O_RDWR, 0600)
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
		term.Success("Configured %s", RcPath)
		return nil
	}

	return ErrShellConfigNotFound
}

type fishAliasConfiger struct{}

func (fishAliasConfiger) SetAlias() error {
	if Exists(RcPath) {
		f, err := os.OpenFile(RcPath, os.O_APPEND|os.O_RDWR, 0600)
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
		term.Success("Configured %s", RcPath)
		return nil
	}

	return ErrShellConfigNotFound
}

func SetAlias() error {
	var c aliasConfiger
	switch Shell {
	case ShellTypeZsh:
		c = zshAliasConfiger{}
	case ShellTypeBash:
		c = bashAliasConfiger{}
	case ShellTypeFish:
		c = fishAliasConfiger{}
	default:
		return ErrUnsupportedShell
	}
	return c.SetAlias()
}
