package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/fvm/model"
	"github.com/lollipopkit/gommon/log"
	"github.com/lollipopkit/gommon/term"
	"github.com/lollipopkit/gommon/util"
)

var (
	ErrVersionNotInstalled = errors.New("Version not installed. Please install it before using.")
	ErrGlobalNotSet        = errors.New("Global version not set. Please set it before using.")

	envNames4JudgeInChina = map[string][]string{
		"TZ":                       {"Asia/Shanghai", "Asia/Chongqing"},
		"LC_ALL":                   {"zh_CN.UTF-8", "zh_CN.GB18030", "zh_CN.GBK"},
		"LANG":                     {"zh_CN.UTF-8", "zh_CN.GB18030", "zh_CN.GBK"},
		"FLUTTER_STORAGE_BASE_URL": {"https://storage.flutter-io.cn"},
		"PUB_HOSTED_URL":           {"https://pub.flutter-io.cn"},
	}
)

func JudgeUseMirror() bool {
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
			log.Err("Save config failed: " + err.Error())
		}
	}
	return *Config.UseMirror
}

func GetReleases() (releases []model.Release, err error) {
	spinner := term.NewSpinner()
	defer spinner.Stop(false)
	goos := GetOS()
	inChina := JudgeUseMirror()
	url := func() string {
		if inChina {
			return fmt.Sprintf(consts.ReleaseChinaJsonUrlFmt, goos)
		}
		return fmt.Sprintf(consts.ReleaseJsonUrlFmt, goos)
	}()
	if inChina {
		spinner.SetString("Using mirror: " + consts.ReleaseChinaUrlPrefix)
	} else {
		spinner.SetString("Using official: " + consts.ReleaseUrlPrefix)
	}

	data, code, err := util.HttpDo("GET", url, nil, nil)
	if code != 200 || err != nil {
		if err == nil {
			err = fmt.Errorf("Get releases failed: code %d", code)
		}
		return
	}

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
		log.Warn("No arm64 version found, will use x64 version.")
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
		log.Warn("Version " + r.Version + " already installed.")
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
			log.Warn("Archive already exists, skip downloading.")
			download = false
		} else {
			log.Warn("Archive already exists, but hash not match, will download again.")
		}
	}

	if download {
		url := func() string {
			if JudgeUseMirror() {
				return consts.ReleaseChinaUrlPrefix + consts.ReleasePath + r.Archive
			}
			return consts.ReleaseUrlPrefix + consts.ReleasePath + r.Archive
		}()

		err := DownloadFile(url, archieve)
		if err != nil {
			return err
		}
	}

	spinner := term.NewSpinner()
	defer spinner.Stop(false)
	spinner.SetString("Checking SHA256...")
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

	spinner.SetString("Version " + r.Version + " installed successfully")

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
		log.Warn("It seems like that you have to config PATH.")
		confirm := term.Confirm("Do you want to automatically config PATH?", true)
		if confirm {
			err = ConfigPath()
			if err != nil {
				return err
			}
		}
		if !confirm {
			log.Warn("Please add the following line to your shell config file:")
			println("export PATH=$PATH:" + filepath.Join(FvmHome, "global", "bin") + "\n")
		}
	}

	log.Suc("Global version -> " + version)
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
		log.Suc("Removed old version: " + dst)
	}

	err = Symlink(installPath, dst)
	if err != nil {
		return err
	}
	log.Suc("Added symlink: " + installPath + " -> " + dst)

	if err = ConfigIde(); err != nil {
		return err
	}

	if err = ConfigGitIgnore(); err != nil {
		return err
	}

	log.Suc("Project Flutter -> " + v)
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

func GetGlobalVersion() (string, error) {
	globalPath := filepath.Join(FvmHome, "global")
	if !Exists(globalPath) {
		return "", ErrGlobalNotSet
	}

	target, err := os.Readlink(globalPath)
	if err != nil {
		return "", err
	}

	return filepath.Base(filepath.Dir(target)), nil
}
