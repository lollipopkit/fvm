package utils

import (
	"archive/zip"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/lollipopkit/fvm/consts"
	"github.com/lollipopkit/gommon/log"
)

var (
	ErrUnsuppotedCompressFormat = fmt.Errorf("unsupported compress format")
	errWinDevModeOff            = errors.New("You may need to turn on 'Developer mode'")

	FvmHome = os.Getenv(consts.FVM_HOME)
)

func init() {
	if FvmHome == "" {
		FvmHome = filepath.Join(consts.HOME, ".fvm")
		log.Warn("FVM_HOME is not set, using default path: " + FvmHome)
	}
	if !Exists(FvmHome) {
		err := os.MkdirAll(FvmHome, 0755)
		if err != nil {
			log.Err("Failed to create FVM_HOME: " + FvmHome)
			os.Exit(1)
		}
	}
}

func Contains[T string | int | int64 | float64](list []T, item T) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func GetArch() string {
	arch := runtime.GOARCH
	switch arch {
	case "amd64":
		return "x64"
	}
	return arch
}

func GetOS() string {
	goos := runtime.GOOS
	switch goos {
	case "darwin":
		return "macos"
	}
	return goos
}

func GetVersionDir(v string) string {
	return FvmHome + "/" + v
}

// Copy from https://stackoverflow.com/questions/20357223/easy-way-to-unzip-file-with-golang
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func Symlink(src, dst string) error {
	switch GetOS() {
	case "windows":
		err := Execute("mklink", "/d", dst, src)
		if err != nil {
			return errors.Join(errWinDevModeOff, err)
		}
		return err
	default:
		return Execute("ln", "-s", src, dst)
	}
}

func IsSymlink(name string) (bool, error) {
	info, err := os.Lstat(name)
	if os.IsNotExist(err) {
		return false, err
	} else if err != nil {
		log.Warn("Error when check symlink: " + err.Error())
		return false, err
	}
	return (info.Mode() & os.ModeSymlink) != 0, nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func Uncompress(src, dest string) error {
	switch filepath.Ext(src) {
	case ".zip":
		if err := Unzip(src, dest); err != nil {
			return err
		}
		return nil
	case ".gz":
		return Execute("tar", "-xzf", src, "-C", dest)
	case ".xz":
		return Execute("tar", "-xJf", src, "-C", dest)
	}
	return ErrUnsuppotedCompressFormat
}

func Execute(bin string, args ...string) error {
	cmd := exec.Command(bin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func GetFileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func AppendIfNotContains(path string, lines2Add []string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")

	for _, line2Add := range lines2Add {
		if !Contains(lines, line2Add) {
			if _, err = f.WriteString("\n" + line2Add); err != nil {
				return err
			}
		}
	}
	return nil
}

func DownloadFile(url string, dest string) error {
	switch GetOS() {
	case "windows":
		return Execute("powershell", "Invoke-WebRequest", "-Uri", url, "-OutFile", dest)
	default:
		return Execute("wget", "-O", dest, url)
	}
}
