package files

import (
	"errors"
	"fmt"
	"os"
)

func FileExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if stat.IsDir() {
		return false, errors.New(fmt.Sprintf("%s is a directory", path))
	}
	return true, nil
}

func DirectoryExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if !stat.IsDir() {
		return false, errors.New(fmt.Sprintf("%s is not a directory", path))
	}
	return true, nil
}

func CreateDirectoryIfNotExists(path string) error {
	crdsDirExists, err := DirectoryExists(path)
	if err != nil {
		return err
	}
	if !crdsDirExists {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
