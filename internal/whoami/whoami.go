package whoami

import (
	"fmt"
	"os"
	"os/user"
	"time"
)

// Info holds contextual information about the current user and environment.
type Info struct {
	Username  string
	Hostname  string
	VaultAddr string
	Profile   string
	Timestamp time.Time
}

// Gather collects current user/environment context.
func Gather(profile string) (*Info, error) {
	u, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("whoami: get user: %w", err)
	}

	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}

	vaultAddr := os.Getenv("VAULT_ADDR")
	if vaultAddr == "" {
		vaultAddr = "(not set)"
	}

	return &Info{
		Username:  u.Username,
		Hostname:  host,
		VaultAddr: vaultAddr,
		Profile:   profile,
		Timestamp: time.Now().UTC(),
	}, nil
}

// Format returns a human-readable summary of the Info.
func (i *Info) Format() string {
	return fmt.Sprintf(
		"User:       %s\nHost:       %s\nVault Addr: %s\nProfile:    %s\nTime:       %s",
		i.Username,
		i.Hostname,
		i.VaultAddr,
		i.Profile,
		i.Timestamp.Format(time.RFC3339),
	)
}
