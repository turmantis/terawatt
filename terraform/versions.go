package terraform

import (
	"context"
	"errors"
	"fmt"
	"github.com/turmantis/terawatt/internal"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/scanner"
)

const versionFileName = ".terraform-version"

var (
	// ErrVersionFileNotFound returned when .terraform-version is not found.
	ErrVersionFileNotFound = errors.New("couldn't find " + versionFileName)
	// ErrRequiredVersionNotFound returned when none of the `.tf` files in the current directory
	// have defined `required_version`.
	ErrRequiredVersionNotFound = errors.New("required_version not found")
)

// DesiredVersion get the desired version of terraform. If the version isn't defined in
// `.terraform-version`, all of the hcl files in the current directory will be searched for a
// `required_version` statement. If no version is found the latest will be used.
func DesiredVersion(ctx context.Context) (string, error) {
	if tv, err := tfEnvVersion(); err == nil {
		return tv, nil
	}
	if hv, err := hclVersion(); err == nil {
		return hv, nil
	}
	return latestVersion(ctx)
}

// latestVersion fetches the latest version of terraform
func latestVersion(ctx context.Context) (string, error) {
	vf, err := getAvailableVersions(ctx)
	if err != nil {
		return "", fmt.Errorf("terraform: %w", err)
	}
	defer func() {
		if err := vf.Close(); err != nil {
			panic(err)
		}
	}()
	versions, err := parseVersionFile(vf)
	if err != nil {
		return "", err
	}
	for _, v := range versions {
		// prerelease versions are in the form 0.0.0-rc0
		if !strings.ContainsRune(v, '-') {
			return v, nil
		}
	}
	return versions[len(versions)-1], nil
}

// hclVersion searches all hcl files in the current directory for `required_version`.
func hclVersion() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("terraform: %w", err)
	}
	entries, err := os.ReadDir(pwd)
	if err != nil {
		return "", fmt.Errorf("terraform: %w", err)
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		path := filepath.Join(pwd, e.Name())
		if filepath.Ext(path) == ".tf" {
			ver, err := parseHclVersion(path)
			if err == nil {
				return ver, nil
			}
		}
	}
	return "", ErrRequiredVersionNotFound
}

// tfEnvVersion get the version from `.terraform-version` if it exists.
func tfEnvVersion() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("terraform: %w", err)
	}
	versionAbsPath := filepath.Join(pwd, versionFileName)
	if _, err = os.Stat(versionAbsPath); os.IsNotExist(err) {
		return "", fmt.Errorf("terraform: %w", ErrVersionFileNotFound)
	}
	bs, err := os.ReadFile(versionAbsPath)
	if err != nil {
		return "", fmt.Errorf("terraform: %w", err)
	}
	return strings.TrimSpace(string(bs)), nil
}

// getAvailableVersions returns a file ptr to a file containing each terraform version on a new
// line. If this file doesn't exist it will be created.
func getAvailableVersions(ctx context.Context) (*os.File, error) {
	allVersionsFile := filepath.Join(internal.TempStorageDirectory, "terraform_versions.html")
	if _, err := os.Stat(allVersionsFile); os.IsNotExist(err) {
		var vf *os.File
		vf, err = internal.CreateOrTruncate(allVersionsFile)
		if err != nil {
			return nil, fmt.Errorf("terraform: %w", err)
		}
		if err = internal.DownloadTo(ctx, repository, vf); err != nil {
			return nil, fmt.Errorf("terraform: %w", err)
		}
		if err = vf.Close(); err != nil {
			panic(err)
		}
	}
	f, err := os.Open(allVersionsFile)
	if err != nil {
		return nil, fmt.Errorf("terraform: %w", err)
	}
	return f, nil
}

// parseVersionFile extract available versions from html.
func parseVersionFile(r io.Reader) ([]string, error) {
	var previousToken string
	var versions []string
	var s scanner.Scanner
	s.Init(r)
	s.Mode = scanner.ScanIdents | scanner.ScanStrings | scanner.ScanRawStrings
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		// Multi-rune tokens start with -2. All of the tokens that matter for this are multi-rune.
		if tok > 0 {
			continue
		}
		t := s.TokenText()
		// version path is always preceded by href.
		if previousToken == "href" {
			// The first link is a relative link. This skips it.
			if t != "\"../\"" {
				v, err := parseVersion(t)
				if err != nil {
					return nil, fmt.Errorf("terraform: %w", err)
				}
				versions = append(versions, v)
			}
		}
		previousToken = t
	}
	return versions, nil
}

// parseVersion "/foo/0.0.0-bar/" -> 0.0.0-bar
func parseVersion(quotedExpr string) (string, error) {
	// v = /foo/0.0.0-bar/
	v, err := strconv.Unquote(quotedExpr)
	if err != nil {
		return "", fmt.Errorf("terraform: %w", err)
	}
	start := 0
	if v[0] == '/' {
		start = 1
	}
	end := len(v)
	if v[len(v)-1] == '/' {
		end -= 1
	}
	// v = foo/0.0.0-bar
	v = v[start:end]
	idx := strings.Index(v, "/")
	idx += 1 // skip slash
	// v = 0.0.0-bar
	return v[idx:], nil
}

func parseHclVersion(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("terraform: %w", err)
	}
	version, err := internal.HclRequiredVersion(f)
	if err != nil {
		return "", fmt.Errorf("terraform: %w", err)
	}
	return version, nil
}
