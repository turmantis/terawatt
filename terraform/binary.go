package terraform

import (
	"context"
	"fmt"
	"github.com/turmantis/terawatt/internal"
	"os"
	"path/filepath"
	"runtime"
)

const repository = "https://releases.hashicorp.com/terraform"

// BinaryHostPath returns the path on the host system for a version of terraform. If it doesn't
// exist it will be downloaded.
func BinaryHostPath(ctx context.Context, version string) (string, error) {
	tfHost := hostBinary(version)
	// If the binary already exists on the host there is nothing to do.
	if _, err := os.Stat(tfHost); os.IsNotExist(err) {
		ap := archivePath(version)
		// if the archive has already been downloaded, skip to the unzipping.
		if _, err = os.Stat(ap); os.IsNotExist(err) {
			if err = os.MkdirAll(filepath.Dir(ap), os.ModePerm); err != nil {
				return "", fmt.Errorf("internal: %w", err)
			}
			var archive *os.File
			archive, err = os.Create(ap)
			if err != nil {
				return "", fmt.Errorf("terraform: %w", err)
			}
			if err = internal.DownloadTo(ctx, downloadUrl(version), archive); err != nil {
				return "", fmt.Errorf("terraform: %w", err)
			}
			if err = archive.Close(); err != nil {
				panic(err)
			}
		}
		var bin *os.File
		bin, err = internal.CreateOrTruncate(hostBinary(version))
		if err != nil {
			return "", fmt.Errorf("terraform: %w", err)
		}
		if err = internal.Unzip(ap, "terraform", bin); err != nil {
			return "", fmt.Errorf("terraform: %w", err)
		}
		if err = bin.Close(); err != nil {
			panic(err)
		}
	}
	return tfHost, nil
}

// https://releases.hashicorp.com/terraform/1.6.6/terraform_1.6.6_linux_amd64.zip
func downloadUrl(version string) string {
	return fmt.Sprintf(
		"%s/%s/terraform_%s_%s_%s.zip",
		repository,
		version,
		version,
		runtime.GOOS,
		runtime.GOARCH,
	)
}

// /tmp/terawatt/terraform_0.1.2.zip
func archivePath(version string) string {
	return filepath.Join(internal.TempStorageDirectory, "terraform_"+version+".zip")
}

// /path/to/user/home/.terawatt/terraform_0.1.2
func hostBinary(version string) string {
	return filepath.Join(internal.StaticStorageDirectory, "terraform_"+version)
}
