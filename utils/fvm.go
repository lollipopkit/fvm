package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/model"
	"github.com/lollipopkit/fvm/term"
)

var (
	ErrVersionNotInstalled    = errors.New("Version not installed. \nPlease install it before using.")
	ErrUnsupportedShellPrefix = "Unsupported shell: "

	envNames4JudgeInChina = map[string][]string{
		"TZ":     {"Asia/Shanghai", "Asia/Chongqing"},
		"LC_ALL": {"zh_CN.UTF-8", "zh_CN.GB18030", "zh_CN.GBK"},
		"LANG":   {"zh_CN.UTF-8", "zh_CN.GB18030", "zh_CN.GBK"},
	}

	IsInChina *bool
)

func InChina(notify bool) bool {
	china := false
	for envName, envValues := range envNames4JudgeInChina {
		envValue := os.Getenv(envName)
		for _, v := range envValues {
			if envValue == v {
				china = true
				break
			}
		}
	}

	if IsInChina == nil {
		result := term.Confirm("Do you want to use mirror site in China?", china)
		IsInChina = &result
		err := SaveConfig()
		if err != nil {
			term.Red("Save config failed: "+err.Error(), true)
		}
	}

	if china && notify {
		term.Yellow("Using mirror site " + consts.ReleaseChinaUrlPrefix)
	}
	return china
}

func GetReleases() (releases []model.Release, err error) {
	goos := GetOS()
	inChina := InChina(true)
	url := func() string {
		if inChina {
			return fmt.Sprintf(consts.ReleaseChinaJsonUrlFmt, goos)
		}
		return fmt.Sprintf(consts.ReleaseJsonUrlFmt, goos)
	}()

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var allReleases model.AllReleases
	err = json.Unmarshal(data, &allReleases)
	if err != nil {
		return
	}

	releases = allReleases.Releases

	return
}

func GetReleaseByVersion(releases []model.Release, version string) (r model.Release, err error) {
	arch := GetArch()
	AGAIN:
	for _, v := range releases {
		if v.Version == version && arch == v.DartSdkArch {
			r = v
			return
		}
	}

	if arch == "arm64" {
		term.Yellow("No arm64 version found, will use x64 version.")
		arch = "x64"
		goto AGAIN
	}

	err = fmt.Errorf("Version %s not found", version)
	return
}

func Install(r model.Release) error {
	tmp := strings.Split(r.Archive, "/")
	if len(tmp) < 3 {
		return fmt.Errorf("Invalid archive name: %s", r.Archive)
	}
	fileName := tmp[2]
	if IsVersionInstalled(r.Version) {
		term.Yellow("\nVersion " + r.Version + " already installed")
		redownload := term.Confirm("Do you want to redownload it?", false)
		if !redownload {
			return nil
		}
	}

	zipPath := path.Join(FvmHome, fileName)

	download := true
	if Exists(zipPath) {
		hash, err := GetFileHash(zipPath)
		if err != nil {
			return err
		}
		if hash == r.Sha256 {
			term.Yellow("Archive already exists, skip downloading.")
			download = false
		} else {
			term.Yellow("Archive already exists, but hash not match, will download again.")
		}
	}

	if download {
		url := func() string {
			if InChina(false) {
				return consts.ReleaseChinaUrlPrefix + consts.ReleasePath + r.Archive
			}
			return consts.ReleaseUrlPrefix + consts.ReleasePath + r.Archive
		}()
		term.Cyan("Downloading " + url)

		err := Execute("wget", "-O", zipPath, url)
		if err != nil {
			return err
		}
	}

	term.Cyan("Uncompressing " + fileName)
	err := os.Mkdir(path.Join(FvmHome, r.Version), 0755)
	if err != nil {
		return err
	}
	err = Uncompress(zipPath, path.Join(FvmHome, r.Version))
	if err != nil {
		return err
	}

	term.Cyan("Removing " + fileName)
	err = os.Remove(zipPath)
	if err != nil {
		return err
	}

	return nil
}

func Global(version string) error {
	installPath := path.Join(FvmHome, version, "flutter")
	if !Exists(installPath) {
		return ErrVersionNotInstalled
	}

	dst := path.Join(FvmHome, "global")
	term.Cyan("Using Flutter " + version)

	err := Symlink(installPath, dst)
	if err != nil {
		return err
	}

	err = Test()
	if err != nil {
		term.Yellow("\nIt seems like that you have to config PATH.")
		unsupport := false
		confirm := term.Confirm("Do you want to automatically config PATH?", true)
		if confirm {
			err = ConfigPath()
			if err != nil {
				if strings.Contains(err.Error(), ErrUnsupportedShellPrefix) {
					unsupport = true
					term.Yellow("Sorry, your shell is not supported.")
				} else {
					return err
				}
			}
		}
		if unsupport || !confirm {
			term.Yellow("Please add the following line to your shell config file:\n\nexport PATH=$PATH:" + path.Join(FvmHome, "global", "bin"))
		}
	}

	return nil
}

func Use(v string) error {
	installPath := path.Join(FvmHome, v, "flutter")
	if !Exists(installPath) {
		return ErrVersionNotInstalled
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	dst := path.Join(wd, consts.FVM_DIR_NAME)
	term.Cyan("Using Flutter " + v)

	err = Execute("ln", "-sf", installPath, dst)
	if err != nil {
		return err
	}

	return ConfigIde()
}

func Test() error {
	cmd := exec.Command("flutter")
	return cmd.Run()
}

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
	

func IsVersionInstalled(version string) bool {
	installPath := path.Join(FvmHome, version, "flutter")
	return Exists(installPath)
}
