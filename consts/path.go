package consts

import "os"

const (
	FvmDirName = ".fvm"

	VscodeDirName     = ".vscode/"
	VscodeSettingPath = VscodeDirName + "settings.json"

	IdeaDirName     = ".idea/"
	IdeaSettingPath = IdeaDirName + "workspace.xml"

	ZshRcName      = ".zshrc"
	BashRcName     = ".bashrc"
	FishConfigPath = ".config/fish/config.fish"
)

var (
	HOME = os.Getenv("HOME")
)
