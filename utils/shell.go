package utils

import (
	"os"
	"path"
)

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
	
	shell = path.Base(shell)
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