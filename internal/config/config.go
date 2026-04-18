package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Profile represents a named Vault sync configuration.
type Profile struct {
	Name      string            `yaml:"name"`
	VaultAddr string            `yaml:"vault_addr"`
	VaultPath string            `yaml:"vault_path"`
	OutputFile string           `yaml:"output_file"`
	Mapping   map[string]string `yaml:"mapping,omitempty"`
}

// Config is the top-level configuration structure.
type Config struct {
	DefaultProfile string    `yaml:"default_profile"`
	Profiles       []Profile `yaml:"profiles"`
}

// Load reads and parses a vaultpull config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// GetProfile returns the profile matching the given name.
func (c *Config) GetProfile(name string) (*Profile, error) {
	if name == "" {
		name = c.DefaultProfile
	}
	for i := range c.Profiles {
		if c.Profiles[i].Name == name {
			return &c.Profiles[i], nil
		}
	}
	return nil, fmt.Errorf("profile %q not found", name)
}

func (c *Config) validate() error {
	if len(c.Profiles) == 0 {
		return fmt.Errorf("at least one profile must be defined")
	}
	seen := map[string]bool{}
	for _, p := range c.Profiles {
		if p.Name == "" {
			return fmt.Errorf("profile name must not be empty")
		}
		if seen[p.Name] {
			return fmt.Errorf("duplicate profile name: %q", p.Name)
		}
		seen[p.Name] = true
		if p.VaultAddr == "" {
			return fmt.Errorf("profile %q: vault_addr is required", p.Name)
		}
		if p.VaultPath == "" {
			return fmt.Errorf("profile %q: vault_path is required", p.Name)
		}
		if p.OutputFile == "" {
			return fmt.Errorf("profile %q: output_file is required", p.Name)
		}
	}
	return nil
}
