package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Profile represents a named sync target.
type Profile struct {
	Name   string `yaml:"name"`
	Path   string `yaml:"path"`
	Output string `yaml:"output"`
}

// Config holds the top-level configuration.
type Config struct {
	VaultAddr  string    `yaml:"vault_addr"`
	VaultToken string    `yaml:"vault_token"`
	Default    string    `yaml:"default"`
	Profiles   []Profile `yaml:"profiles"`
}

// Load reads and parses the YAML config file at path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if cfg.VaultAddr == "" {
		return nil, errors.New("vault_addr is required")
	}

	return &cfg, nil
}

// GetProfile returns the profile matching name, falling back to Default, or error.
func (c *Config) GetProfile(name string) (*Profile, error) {
	if name == "" {
		name = c.Default
	}
	for i := range c.Profiles {
		if c.Profiles[i].Name == name {
			return &c.Profiles[i], nil
		}
	}
	return nil, fmt.Errorf("profile %q not found", name)
}
