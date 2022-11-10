package utils

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/lollipopkit/fvm/term"
)

func ConfigGitIgnore() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	gitIgnoreFile := path.Join(wd, ".gitignore")
	term.Info("\nConfiguring \".gitignore\"...")
	if Exists(gitIgnoreFile) {
		f, err := os.OpenFile(gitIgnoreFile, os.O_APPEND|os.O_RDWR, 0600)
		if err != nil {
			return err
		}

		defer f.Close()

		lines, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		line2Add := ".fvm"
		if !Contains(strings.Split(string(lines), "\n"), line2Add) {
			if _, err = f.WriteString("\n" + line2Add); err != nil {
				return err
			}
			term.Success("Configured \".gitignore\"")
		} else {
			term.Success("\".gitignore\" already configured. Skip.")
		}
	} else {
		if err := ioutil.WriteFile(gitIgnoreFile, []byte(".fvm"), 0644); err != nil {
			return err
		}
		term.Success("Configured \".gitignore\"")
	}
	return nil
}
