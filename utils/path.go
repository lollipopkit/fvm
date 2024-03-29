package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/lollipopkit/gommon/log"
)

func ConfigPath() error {
	switch Shell {
	case ShellTypeZsh, ShellTypeBash:
		return configPath4Bash()
	case ShellTypeFish:
		return configPath4Fish()
	default:
		return ErrUnsupportedShell
	}
}

func configPath4Bash() error {
	if !Exists(RcPath) {
		return ErrShellConfigNotFound
	}

	content, err := ioutil.ReadFile(RcPath)
	if err != nil {
		return err
	}

	line2Add := "export PATH=$PATH:" + filepath.Join(FvmHome, "global", "bin")
	lines := strings.Split(string(content), "\n")
	if Contains(lines, line2Add) {
		log.Info("\nPATH already configured. Skip.")
	} else {
		f, err := os.OpenFile(RcPath, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err = f.WriteString("\n" + line2Add); err != nil {
			return err
		}
	}

	log.Warn("\nPlease run following command to reload shell config file:")
	println("source " + RcPath + "\n")
	return nil
}

func configPath4Fish() error {
	if !Exists(RcPath) {
		return ErrShellConfigNotFound
	}

	content, err := ioutil.ReadFile(RcPath)
	if err != nil {
		return err
	}

	line2Add := "set PATH " + filepath.Join(FvmHome, "global", "bin") + "$PATH"
	lines := strings.Split(string(content), "\n")
	if Contains(lines, line2Add) {
		log.Info("\nPATH already configured. Skip.")
	} else {
		f, err := os.OpenFile(RcPath, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}

		defer f.Close()

		if _, err = f.WriteString("\n" + line2Add); err != nil {
			return err
		}
	}

	log.Warn("\nPlease run following command to reload shell config file:")
	println("source " + RcPath + "\n")
	return nil
}
