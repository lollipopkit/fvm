package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/LollipopKit/gofvm/consts"
	"github.com/LollipopKit/gofvm/model"
	"github.com/LollipopKit/gofvm/term"
)

var (
	ErrPathNotSet             = errors.New("FVM_PATH not set. \nPlease set it in ENV before using gofvm.")
	ErrVersionNotInstalled    = errors.New("Version not installed. \nPlease install it before using.")
	ErrUnsupportedShellPrefix = "Unsupported shell: "

	envNames4JudgeInChina = map[string][]string{
		"TZ":     {"Asia/Shanghai", "Asia/Chongqing"},
		"LC_ALL": {"zh_CN.UTF-8", "zh_CN.GB18030", "zh_CN.GBK"},
		"LANG":   {"zh_CN.UTF-8", "zh_CN.GB18030", "zh_CN.GBK"},
	}

	IsInChina *bool
)

func Precheck() error {
	if os.Getenv(consts.PATH_NAME) == "" {
		return ErrPathNotSet
	}
	return nil
}

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
		result := Confirm("Do you want to use the mirror in China?", china)
		IsInChina = &result
	}

	if china && notify {
		term.Yellow("Using mirror site " + consts.ReleaseChinaUrlPrefix)
	}
	return china
}

func GetReleases() (releases []model.Release, err error) {
	goarch := GetArch()
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

	archs := []string{}
	for idx := range allReleases.Releases {
		a := allReleases.Releases[idx].DartSdkArch
		if !Contains(archs, a) {
			archs = append(archs, a)
		}
	}

	if !Contains(archs, goarch) {
		err = fmt.Errorf("Supported archs: %v, but your arch: %s", archs, goarch)
		return
	}

	for idx := range allReleases.Releases {
		if allReleases.Releases[idx].DartSdkArch == goarch {
			releases = append(releases, allReleases.Releases[idx])
		}
	}

	return
}

func GetReleaseByVersion(releases []model.Release, version string) (r model.Release, err error) {
	for _, v := range releases {
		if v.Version == version {
			r = v
			return
		}
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

	zipPath := path.Join(Path(), fileName)

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

		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Download failed: %s", resp.Status)
		}

		out, err := os.Create(zipPath)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
	}

	term.Cyan("Uncompressing " + fileName)
	err := os.Mkdir(path.Join(Path(), r.Version), 0755)
	if err != nil {
		return err
	}
	err = Uncompress(zipPath, path.Join(Path(), r.Version))
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

func Use(version string) error {
	installPath := path.Join(Path(), version, "flutter")
	if !Exists(installPath) {
		return ErrVersionNotInstalled
	}

	dst := path.Join(Path(), "global")
	term.Cyan("Using Flutter " + version)

	err := Execute("ln", "-sf", installPath, dst)
	if err != nil {
		return err
	}

	err = Test()
	if err != nil {
		term.Yellow("\nIt seems like that you have to config PATH.")
		unsupport := false
		confirm := Confirm("Do you want to automatically config PATH?", true)
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
			term.Yellow("Please add the following line to your shell config file:\n\nexport PATH=$PATH:" + path.Join(Path(), "global", "bin"))
		}
	}

	return nil
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

	shellConfigFile := ""
	switch shell {
	case "/bin/bash":
		shellConfigFile = path.Join(os.Getenv("HOME"), ".bashrc")
	case "/bin/zsh", "/usr/bin/zsh":
		shellConfigFile = path.Join(os.Getenv("HOME"), ".zshrc")
	default:
		return fmt.Errorf(ErrUnsupportedShellPrefix+"%s", shell)
	}

	if !Exists(shellConfigFile) {
		return fmt.Errorf("Shell config file not found: %s", shellConfigFile)
	}

	content, err := ioutil.ReadFile(shellConfigFile)
	if err != nil {
		return err
	}

	line2Add := "export PATH=$PATH:" + path.Join(Path(), "global", "bin")
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
