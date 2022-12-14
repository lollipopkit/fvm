package utils

import (
	"os"
	"path/filepath"

	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/term"
)

func ConfigGitIgnore() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	gitIgnoreFile := filepath.Join(wd, ".gitignore")
	if Exists(gitIgnoreFile) {
		err := AppendIfNotContains(gitIgnoreFile, []string{consts.FvmDirName})
		if err != nil {
			return err
		}
		term.Success(".gitignore already configured. Skip.")
		return nil
	} else {
		if err := os.WriteFile(gitIgnoreFile, []byte(consts.FvmDirName), 0644); err != nil {
			return err
		}
	}
	term.Success("Configured .gitignore")
	return nil
}
