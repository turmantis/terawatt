package internal

import (
	"os"
	"path/filepath"
)

const (
	version   = "0.0.0"
	userAgent = "Terawatt/" + version
	twDirName = ".terawatt"
)

var (
	StaticStorageDirectory = filepath.Join(mustGetHomeDir(), twDirName)
	TempStorageDirectory   = filepath.Join(os.TempDir(), twDirName)
)

func mustGetHomeDir() string {
	pwd, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return pwd
}
