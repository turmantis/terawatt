package internal

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
)

var ErrArchiveFileNotFound = errors.New("archive file not found")

// Unzip extracts a file from an archive and writes it to a destination.
func Unzip(archive string, file string, dest io.Writer) error {
	r, err := zip.OpenReader(archive)
	if err != nil {
		return fmt.Errorf("internal: %w", err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()
	for _, f := range r.File {
		if f.Name == file {
			var rc io.ReadCloser
			rc, err = f.Open()
			if err != nil {
				return fmt.Errorf("internal: %w", err)
			}
			if _, err = io.Copy(dest, rc); err != nil {
				return fmt.Errorf("internal: %w", err)
			}
			if err = rc.Close(); err != nil {
				return fmt.Errorf("internal: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("internal: %w: %s %s", ErrArchiveFileNotFound, file, archive)
}
