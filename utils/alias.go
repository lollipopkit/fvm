package utils

import (
	"github.com/lollipopkit/gommon/log"
)

var (
	aliasLines2Add = []string{"alias dart='fvm dart'", "alias flutter='fvm flutter'"}
)

type aliasConfiger interface {
	SetAlias() error
}

type fishAliasConfiger struct{}

func (fishAliasConfiger) SetAlias() error {
	if Exists(RcPath) {
		err := AppendIfNotContains(RcPath, aliasLines2Add)
		if err != nil {
			return err
		}
		log.Suc("Configured %s", RcPath)
		return nil
	}

	return ErrShellConfigNotFound
}

func SetAlias() error {
	var c aliasConfiger
	switch Shell {
	case ShellTypeBash, ShellTypeZsh, ShellTypeFish:
		c = fishAliasConfiger{}
	default:
		return ErrUnsupportedShell
	}

	err := c.SetAlias()
	if err != nil {
		return err
	}

	log.Warn("\nPlease run following command to reload shell config file:")
	println("source " + RcPath)
	return nil
}
