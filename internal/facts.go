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
	StaticStorageDirectory = filepath.Join(must(os.UserHomeDir()), twDirName)
	TempStorageDirectory   = filepath.Join(os.TempDir(), twDirName)
)
