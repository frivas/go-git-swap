package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"github.com/frivas/go-git-swap/internal/model"
)

// FileName is the profile store, relative to the user's home directory.
const FileName = ".go-git-swap.json"

type Config struct {
	filepath string
	profiles *model.Profiles
}

// New opens the profile store in the user's home directory.
func New() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return Open(filepath.Join(homeDir, FileName))
}

// Open loads the profile store at path, creating it if it does not exist.
func Open(path string) (*Config, error) {
	cfg := &Config{
		filepath: path,
		profiles: &model.Profiles{Profiles: make(map[string]model.GitProfile)},
	}

	if err := cfg.load(); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
		if err := cfg.save(); err != nil {
			return nil, err
		}
	}
	if cfg.profiles.Profiles == nil {
		cfg.profiles.Profiles = make(map[string]model.GitProfile)
	}
	// Older versions created the file world-readable.
	if err := os.Chmod(path, 0o600); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) load() error {
	data, err := os.ReadFile(c.filepath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, c.profiles); err != nil {
		return fmt.Errorf("parse %s: %w", c.filepath, err)
	}
	return nil
}

// save writes atomically so a crash never leaves a truncated store.
func (c *Config) save() error {
	data, err := json.MarshalIndent(c.profiles, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(c.filepath), ".go-git-swap-*.json")
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(tmp.Name()) }() // no-op after a successful rename
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	// CreateTemp already uses 0600.
	return os.Rename(tmp.Name(), c.filepath)
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

// ListProfiles returns profile names sorted, so menu numbers are stable.
func (c *Config) ListProfiles() []string {
	names := make([]string, 0, len(c.profiles.Profiles))
	for name := range c.profiles.Profiles {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

func (c *Config) RemoveProfile(name string) error {
	if _, exists := c.profiles.Profiles[name]; !exists {
		return fmt.Errorf("profile %s does not exist", name)
	}
	delete(c.profiles.Profiles, name)
	if c.profiles.Active == name {
		c.profiles.Active = ""
	}
	return c.save()
}
