package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/LollipopKit/gofvm/consts"
	"github.com/LollipopKit/gofvm/model"
	"github.com/LollipopKit/gofvm/term"
)

var (
	ErrPathNotSet = errors.New("FVM_PATH not set. \nPlease set it in ENV before using gofvm.")

	envNames4JudgeInChina = map[string][]string{
		"TZ":     {"Asia/Shanghai", "Asia/Chongqing"},
		"LC_ALL": {"zh_CN.UTF-8", "zh_CN.GB18030", "zh_CN.GBK"},
		"LANG":   {"zh_CN.UTF-8", "zh_CN.GB18030", "zh_CN.GBK"},
	}
)

func Precheck() error {
	if os.Getenv(consts.PATH_NAME) == "" {
		return ErrPathNotSet
	}
	return nil
}

func InChina() bool {
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
	return china
}

func GetReleases() (releases []model.Release, err error) {
	goarch := GetArch()
	goos := GetOS()
	inChina := InChina()
	url := func() string {
		if inChina {
			return fmt.Sprintf(consts.ReleaseChinaJsonUrlFmt, goos)
		}
		return fmt.Sprintf(consts.ReleaseJsonUrlFmt, goos)
	}()

	if inChina {
		term.Yellow("Using mirror site " + consts.ReleaseChinaUrlPrefix)
	}

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
	url := consts.ReleaseDownloadUrlPrefix + r.Archive
	term.Cyan("Downloading " + url)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Download failed: %s", resp.Status)
	}

	tmp := strings.Split(r.Archive, "/")
	if len(tmp) < 3 {
		return fmt.Errorf("Invalid archive name: %s", r.Archive)
	}
	fileName := tmp[2]

	zipPath := path.Join(Path(), fileName)
	out, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	term.Cyan("Unzipping " + fileName)
	err = Unzip(zipPath, path.Join(Path(), r.Version))
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
	if !Exists(path.Join(Path(), version)) {
		return fmt.Errorf("Version %s not found", version)
	}
	term.Cyan("Using Flutter " + version)

	return os.Symlink(path.Join(Path(), version), path.Join(Path(), "current"))
}
