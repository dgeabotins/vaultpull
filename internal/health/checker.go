package health

import (
	"fmt"
	"net/http"
	"time"
)

// Status represents the health of a dependency.
type Status struct {
	Name    string
	OK      bool
	Message string
}

// Report holds all health check results.
type Report struct {
	Statuses []Status
	Healthy  bool
}

// Checker runs health checks against Vault and local config.
type Checker struct {
	vaultAddr string
	client    *http.Client
}

// New creates a Checker for the given Vault address.
func New(vaultAddr string) *Checker {
	return &Checker{
		vaultAddr: vaultAddr,
		client:    &http.Client{Timeout: 5 * time.Second},
	}
}

// Run executes all health checks and returns a Report.
func (c *Checker) Run() Report {
	statuses := []Status{
		c.checkVault(),
	}

	healthy := true
	for _, s := range statuses {
		if !s.OK {
			healthy = false
			break
		}
	}

	return Report{Statuses: statuses, Healthy: healthy}
}

func (c *Checker) checkVault() Status {
	url := fmt.Sprintf("%s/v1/sys/health", c.vaultAddr)
	resp, err := c.client.Get(url)
	if err != nil {
		return Status{Name: "vault", OK: false, Message: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusTooManyRequests {
		return Status{Name: "vault", OK: true, Message: fmt.Sprintf("reachable (HTTP %d)", resp.StatusCode)}
	}
	return Status{Name: "vault", OK: false, Message: fmt.Sprintf("unexpected status %d", resp.StatusCode)}
}
