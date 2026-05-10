package cache

import (
	"fmt"
	"os"
	"path/filepath"
)

func CacheLocation() (path string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		err = fmt.Errorf("unable to determine user's home directory: %w", err)
		return
	}

	path = filepath.Join(home, ".cache", "civitaicli")

	return filepath.Abs(path)
}

func EnsureCacheLocation() (path string, err error) {
	path, err = CacheLocation()
	if err != nil {
		err = fmt.Errorf("unable to determine cache location: %w", err)
		return
	}

	err = os.MkdirAll(path, 0750)
	if err != nil {
		err = fmt.Errorf("unable to create cache directory: %w", err)
		return
	}

	return
}
