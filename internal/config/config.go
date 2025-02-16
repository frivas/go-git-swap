package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/frivas/go-git-swap/internal/model"
)

type Config struct {
	filepath string
	profiles *model.Profiles
}

func New() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".go-git-swap.json")
	cfg := &Config{
		filepath: configPath,
		profiles: &model.Profiles{
			Active:   "",
			Profiles: make(map[string]model.GitProfile),
		},
	}

	if err := cfg.load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		// Create empty config file
		if err := cfg.save(); err != nil {
			return nil, err
		}
	}

	// Ensure active profile is empty if no profiles exist
	if len(cfg.profiles.Profiles) == 0 {
		cfg.profiles.Active = ""
		if err := cfg.save(); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func (c *Config) load() error {
	data, err := os.ReadFile(c.filepath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, c.profiles)
}

func (c *Config) save() error {
	data, err := json.MarshalIndent(c.profiles, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.filepath, data, 0644)
}

func (c *Config) AddProfile(profile model.GitProfile) error {
	if _, exists := c.profiles.Profiles[profile.Name]; exists {
		return fmt.Errorf("profile %s already exists", profile.Name)
	}

	c.profiles.Profiles[profile.Name] = profile
	return c.save()
}

func (c *Config) GetProfile(name string) (model.GitProfile, bool) {
	profile, exists := c.profiles.Profiles[name]
	return profile, exists
}

func (c *Config) SetActiveProfile(name string) error {
	if _, exists := c.profiles.Profiles[name]; !exists {
		return fmt.Errorf("profile %s does not exist", name)
	}

	c.profiles.Active = name
	return c.save()
}

func (c *Config) GetActiveProfile() string {
	return c.profiles.Active
}

func (c *Config) ListProfiles() []string {
	profiles := make([]string, 0, len(c.profiles.Profiles))
	for name := range c.profiles.Profiles {
		profiles = append(profiles, name)
	}
	return profiles
}

func (c *Config) RemoveProfile(name string) error {
	if _, exists := c.profiles.Profiles[name]; !exists {
		return fmt.Errorf("profile %s does not exist", name)
	}

	delete(c.profiles.Profiles, name)
	return c.save()
}
