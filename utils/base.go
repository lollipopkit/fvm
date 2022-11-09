package utils

import (
	"archive/zip"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/lollipopkit/gofvm/consts"
	"github.com/lollipopkit/gofvm/term"
)

var (
	ErrUnsuppotedCompressFormat = fmt.Errorf("unsupported compress format")
)

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

func Path() string {
	return os.Getenv(consts.PATH_NAME)
}

func GetVersionDir(v string) string {
	return Path() + "/" + v
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
	if err := os.Symlink(src, dst); err != nil {
		return err
	}
	return nil
}

func IsSymlink(name string) (bool, error) {
	info, err := os.Lstat(name)
	if os.IsNotExist(err) {
		return false, err
	} else if err != nil {
		term.Yellow("Error when check symlink: " + err.Error())
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
