package internal

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/scanner"
)

var ErrRequiredVersionNotFound = errors.New("required version not found")

// HclRequiredVersion parse HCL looking for the expression `required_version = "~> 1.2.3"` to
// extract the semantic version.
func HclRequiredVersion(in io.Reader) (string, error) {
	var previousToken string
	var s scanner.Scanner
	s.Init(in)
	s.Mode = scanner.ScanIdents | scanner.ScanStrings | scanner.ScanRawStrings
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		// skip single character tokens
		if tok > 0 {
			continue
		}
		t := s.TokenText()
		// The needle being searched for is a single expression
		// required_version = "~> 1.2.3"
		if previousToken == "required_version" {
			ver, err := parseHclVersionExpression(t)
			if err != nil {
				return "", err
			}
			return ver, nil // needle found
		}
		previousToken = t
	}
	return "", fmt.Errorf("internal: %w", ErrRequiredVersionNotFound)
}

// parseHclVersionExpression '"~> 1.1.9"' -> '1.1.9'
func parseHclVersionExpression(expr string) (string, error) {
	expr, err := strconv.Unquote(expr)
	if err != nil {
		return "", fmt.Errorf("internal: %w", err)
	}
	// drop any spaces that might have been in the quoted string before the expression.
	expr = strings.TrimSpace(expr)
	if strings.HasPrefix(expr, "~>") {
		expr = strings.TrimPrefix(expr, "~>")
		expr = strings.TrimSpace(expr)
		return expr, nil
	}
	return "", fmt.Errorf("internal: unable to parse expression '%s'", expr)
}
