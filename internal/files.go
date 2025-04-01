package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateOrTruncate(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return nil, fmt.Errorf("internal: %w", err)
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("internal: %w", err)
	}
	return file, nil
}
