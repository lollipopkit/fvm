package utils

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/lollipopkit/fvm/consts"
)

var (
	Shell     ShellType
	RcPath    string
	ShellName string

	errNone                = errors.New("")
	ErrShellConfigNotFound = errNone
	ErrUnsupportedShell    = errNone
)

func init() {
	go func() {
		Shell = GetShell()
		RcPath = Shell.RcPath()
		ShellName = Shell.String()
		ErrShellConfigNotFound = errors.New("Shell config file not found: " + RcPath)
		ErrUnsupportedShell = errors.New("Unsupported shell: " + ShellName)
	}()
}

type ShellType uint8

const (
	ShellTypeUnknown ShellType = iota
	ShellTypeBash
	ShellTypeZsh
	ShellTypeFish
)

func GetShell() ShellType {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return ShellTypeUnknown
	}

	shell = filepath.Base(shell)
	switch shell {
	case "zsh":
		return ShellTypeZsh
	case "bash":
		return ShellTypeBash
	case "fish":
		return ShellTypeFish
	default:
		return ShellTypeUnknown
	}
}

func (s ShellType) String() string {
	switch s {
	case ShellTypeZsh:
		return "zsh"
	case ShellTypeBash:
		return "bash"
	case ShellTypeFish:
		return "fish"
	default:
		return "unknown"
	}
}

func (s ShellType) RcPath() string {
	switch s {
	case ShellTypeZsh:
		return filepath.Join(consts.HOME, consts.ZshRcName)
	case ShellTypeBash:
		return filepath.Join(consts.HOME, consts.BashRcName)
	case ShellTypeFish:
		return filepath.Join(consts.HOME, consts.FishConfigPath)
	default:
		return ""
	}
}
