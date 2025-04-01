package internal

import (
	"strings"
	"testing"
)

func TestHclRequiredVersion(t *testing.T) {
	got, err := HclRequiredVersion(strings.NewReader(sampleVersion))
	if err != nil {
		t.Errorf("HclRequiredVersion() error = %v", err)
		return
	}
	if got != "1.2.3" {
		t.Errorf("HclRequiredVersion() got = %v", got)
	}
}

const sampleVersion = `
terraform {
  required_version = "~> 1.2.3"
}
`
