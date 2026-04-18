package version_test

import (
	"strings"
	"testing"

	"github.com/user/vaultpull/internal/version"
)

func TestInfo_ContainsVersion(t *testing.T) {
	info := version.Info()
	if !strings.Contains(info, version.Version) {
		t.Errorf("Info() = %q, want it to contain version %q", info, version.Version)
	}
}

func TestInfo_ContainsCommit(t *testing.T) {
	info := version.Info()
	if !strings.Contains(info, version.Commit) {
		t.Errorf("Info() = %q, want it to contain commit %q", info, version.Commit)
	}
}

func TestInfo_ContainsDate(t *testing.T) {
	info := version.Info()
	if !strings.Contains(info, version.Date) {
		t.Errorf("Info() = %q, want it to contain date %q", info, version.Date)
	}
}

func TestShort_ReturnsVersion(t *testing.T) {
	if got := version.Short(); got != version.Version {
		t.Errorf("Short() = %q, want %q", got, version.Version)
	}
}
