package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/model"
	"github.com/lollipopkit/gommon/term"
)

var (
	ErrVersionNotInstalled = errors.New("Version not installed. Please install it before using.")

	envNames4JudgeInChina = map[string][]string{
		"TZ":                       {"Asia/Shanghai", "Asia/Chongqing"},
		"LC_ALL":                   {"zh_CN.UTF-8", "zh_CN.GB18030", "zh_CN.GBK"},
		"LANG":                     {"zh_CN.UTF-8", "zh_CN.GB18030", "zh_CN.GBK"},
		"FLUTTER_STORAGE_BASE_URL": {"https://storage.flutter-io.cn"},
		"PUB_HOSTED_URL":           {"https://pub.flutter-io.cn"},
	}
	spinner = term.NewSpinner()
)

func JudgeUseMirror(notify bool) bool {
	if Config.UseMirror == nil {
		china := false
		for envName, envValues := range envNames4JudgeInChina {
			envValue := os.Getenv(envName)
			for _, v := range envValues {
				if envValue == v {
					china = true
					break
				}
			}
			if china {
				break
			}
		}

		result := term.Confirm("Do you want to use mirror site in China?", china)
		Config.UseMirror = &result
		err := SaveConfig()
		if err != nil {
			term.Err("Save config failed: " + err.Error())
		}
	}

	if *Config.UseMirror && notify {
		term.Info("Using mirror site " + consts.ReleaseChinaUrlPrefix)
	}
	return *Config.UseMirror
}

func GetReleases() (releases []model.Release, err error) {
	goos := GetOS()
	inChina := JudgeUseMirror(true)
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
	data, err := io.ReadAll(resp.Body)
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
		term.Warn("No arm64 version found, will use x64 version.")
		arch = "x64"
		goto AGAIN
	}

	err = fmt.Errorf("Version %s not found", version)
	return
}

func Install(r model.Release, force bool) error {
	tmp := strings.Split(r.Archive, "/")
	if len(tmp) < 3 {
		return fmt.Errorf("Invalid archive name: %s", r.Archive)
	}
	fileName := tmp[2]
	if IsVersionInstalled(r.Version) && !force {
		term.Warn("Version " + r.Version + " already installed.")
		redownload := term.Confirm("Do you want to redownload it?", false)
		if !redownload {
			return nil
		}
	}

	archieve := filepath.Join(FvmHome, fileName)
	
	download := true
	if Exists(archieve) {
		hash, err := GetFileHash(archieve)
		if err != nil {
			return err
		}
		if hash == r.Sha256 {
			term.Warn("Archive already exists, skip downloading.")
			download = false
		} else {
			term.Warn("Archive already exists, but hash not match, will download again.")
		}
	}

	if download {
		url := func() string {
			if JudgeUseMirror(false) {
				return consts.ReleaseChinaUrlPrefix + consts.ReleasePath + r.Archive
			}
			return consts.ReleaseUrlPrefix + consts.ReleasePath + r.Archive
		}()

		err := DownloadFile(url, archieve)
		if err != nil {
			return err
		}
	}

	spinner.SetString("Checking SHA256...")
	spinner.Start(100*time.Millisecond)
	hash, err := GetFileHash(archieve)
	if err != nil {
		return err
	}
	if hash != r.Sha256 {
		return fmt.Errorf("SHA256 not match, expect %s, got %s", r.Sha256, hash)
	}

	spinner.SetString("Uncompressing " + fileName)
	err = os.Mkdir(filepath.Join(FvmHome, r.Version), 0755)
	if err != nil {
		return err
	}
	err = Uncompress(archieve, filepath.Join(FvmHome, r.Version))
	if err != nil {
		return err
	}

	spinner.SetString("Removing " + fileName)
	err = os.Remove(archieve)
	if err != nil {
		return err
	}

	spinner.Stop()
	term.Info("Version " + r.Version + " installed successfully")

	return nil
}

func Global(version string) error {
	installPath := filepath.Join(FvmHome, version, "flutter")
	if !Exists(installPath) {
		return ErrVersionNotInstalled
	}

	dst := filepath.Join(FvmHome, "global")

	err := Execute("rm", "-rf", dst)
	if err != nil {
		return err
	}

	err = Symlink(installPath, dst)
	if err != nil {
		return err
	}

	err = TestFlutter()
	if err != nil {
		term.Warn("It seems like that you have to config PATH.")
		confirm := term.Confirm("Do you want to automatically config PATH?", true)
		if confirm {
			err = ConfigPath()
			if err != nil {
				return err
			}
		}
		if !confirm {
			term.Warn("Please add the following line to your shell config file:")
			println("export PATH=$PATH:" + filepath.Join(FvmHome, "global", "bin") + "\n")
		}
	}

	term.Suc("Global version -> " + version)
	return nil
}

func Use(v string) error {
	installPath := filepath.Join(FvmHome, v, "flutter")
	if !Exists(installPath) {
		return ErrVersionNotInstalled
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	dst := filepath.Join(wd, consts.FvmDirName)
	err = os.RemoveAll(dst)
	if err == nil {
		term.Suc("Removed old version: " + dst)
	}

	err = Symlink(installPath, dst)
	if err != nil {
		return err
	}
	term.Suc("Added symlink: " + installPath + " -> " + dst)

	if err = ConfigIde(); err != nil {
		return err
	}

	if err = ConfigGitIgnore(); err != nil {
		return err
	}

	term.Suc("Project Flutter -> " + v)
	return nil
}

func TestFvm() error {
	cmd := exec.Command("fvm")
	return cmd.Run()
}

func TestFlutter() error {
	cmd := exec.Command("flutter")
	return cmd.Run()
}

func IsVersionInstalled(version string) bool {
	installPath := filepath.Join(FvmHome, version, "flutter")
	return Exists(installPath)
}

func Delete(version string) error {
	installPath := filepath.Join(FvmHome, version)
	if !Exists(installPath) {
		return ErrVersionNotInstalled
	}

	err := os.RemoveAll(installPath)
	if err != nil {
		return err
	}
	return nil
}
