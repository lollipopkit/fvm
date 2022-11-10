package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/lollipopkit/fvm/term"
)

func ConfigPath() error {
	term.Cyan("\nConfiguring PATH...")
	shell := os.Getenv("SHELL")
	if shell == "" {
		return fmt.Errorf("Can not get SHELL env")
	}
	shell = path.Base(shell)

	shellConfigFile := ""
	switch shell {
	case "bash":
		shellConfigFile = path.Join(os.Getenv("HOME"), ".bashrc")
	case "zsh":
		shellConfigFile = path.Join(os.Getenv("HOME"), ".zshrc")
	case "fish":
		shellConfigFile = path.Join(os.Getenv("HOME"), ".config", "fish", "config.fish")
	default:
		return fmt.Errorf(ErrUnsupportedShellPrefix+"%s", shell)
	}

	switch shell {
	case "bash", "zsh":
		return configPath4Bash(shellConfigFile)
	case "fish":
		return configPath4Fish(shellConfigFile)
	default:
		return fmt.Errorf(ErrUnsupportedShellPrefix+"%s", shell)
	}
}

func configPath4Bash(shellConfigFile string) error {
	if !Exists(shellConfigFile) {
		return fmt.Errorf("Shell config file not found: %s", shellConfigFile)
	}

	content, err := ioutil.ReadFile(shellConfigFile)
	if err != nil {
		return err
	}

	line2Add := "export PATH=$PATH:" + path.Join(FvmHome, "global", "bin")
	lines := strings.Split(string(content), "\n")
	if Contains(lines, line2Add) {
		term.Yellow("PATH already configured. Skip.")
		return nil
	}

	f, err := os.OpenFile(shellConfigFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.WriteString("\n" + line2Add); err != nil {
		return err
	}

	term.Cyan("Please run following command to reload shell config file:\n\nsource " + shellConfigFile)
	return nil
}

func configPath4Fish(shellConfigFile string) error {
	if !Exists(shellConfigFile) {
		return fmt.Errorf("Shell config file not found: %s", shellConfigFile)
	}

	content, err := ioutil.ReadFile(shellConfigFile)
	if err != nil {
		return err
	}

	line2Add := "set PATH " + path.Join(FvmHome, "global", "bin") + "$PATH"
	lines := strings.Split(string(content), "\n")
	if Contains(lines, line2Add) {
		term.Yellow("PATH already configured. Skip.")
		return nil
	}

	f, err := os.OpenFile(shellConfigFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString("\n" + line2Add); err != nil {
		return err
	}

	term.Cyan("Please run following command to reload shell config file:\n\nsource " + shellConfigFile)
	return nil
}
