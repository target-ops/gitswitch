// Package identity is the in-memory model of a gitswitch identity and the
// JSON-backed config store. The store is the single source of truth for
// "what identities exist and where are they bound." Other layers
// (~/.gitconfig, ~/.ssh/config) are derived state — populated by
// `gitswitch use` based on what's stored here.
package identity

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ConfigVersion is the on-disk schema version. Bump whenever the JSON
// shape changes in a way that requires migration.
const ConfigVersion = 1

// Identity is one git/ssh/gh identity that travels together.
type Identity struct {
	Name       string `json:"name"`                  // user-given alias: "work", "personal"
	Email      string `json:"email"`
	GitName    string `json:"git_name,omitempty"`    // git config user.name
	SSHKey     string `json:"ssh_key,omitempty"`     // absolute path to private key
	SigningKey string `json:"signing_key,omitempty"` // absolute path to .pub used for SSH commit signing
	GHAccount  string `json:"gh_account,omitempty"`
	Vendor     string `json:"vendor,omitempty"` // github | gitlab | bitbucket | ...
}

// Binding is a directory ↔ identity association. The directory is an
// absolute path; an identity binds to all subdirectories beneath it.
type Binding struct {
	Directory string `json:"directory"`
	Identity  string `json:"identity"` // matches Identity.Name
}

// Config is the on-disk JSON document at ~/.config/gitswitch/config.json.
type Config struct {
	Version    int        `json:"version"`
	Identities []Identity `json:"identities,omitempty"`
	Bindings   []Binding  `json:"bindings,omitempty"`
}

// FindByName returns a pointer into the slice (or nil), so callers can
// mutate without repeating the search.
func (c *Config) FindByName(name string) *Identity {
	for i := range c.Identities {
		if c.Identities[i].Name == name {
			return &c.Identities[i]
		}
	}
	return nil
}

// FindByEmail returns the first identity matching this email, or nil.
// Email comparison is case-insensitive.
func (c *Config) FindByEmail(email string) *Identity {
	for i := range c.Identities {
		if equalFold(c.Identities[i].Email, email) {
			return &c.Identities[i]
		}
	}
	return nil
}

// Upsert adds an identity or replaces an existing one with the same Name.
func (c *Config) Upsert(id Identity) {
	for i := range c.Identities {
		if c.Identities[i].Name == id.Name {
			c.Identities[i] = id
			return
		}
	}
	c.Identities = append(c.Identities, id)
}

// Path returns the absolute path to ~/.config/gitswitch/config.json.
// Honours $XDG_CONFIG_HOME when set.
func Path() string {
	if x := os.Getenv("XDG_CONFIG_HOME"); x != "" {
		return filepath.Join(x, "gitswitch", "config.json")
	}
	return filepath.Join(os.Getenv("HOME"), ".config", "gitswitch", "config.json")
}

// Load reads the config from disk. Returns an empty Config (not an error)
// when the file doesn't exist — that's the legitimate "fresh machine"
// state and callers shouldn't have to special-case it.
func Load() (*Config, error) {
	path := Path()
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{Version: ConfigVersion}, nil
		}
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if c.Version == 0 {
		c.Version = ConfigVersion
	}
	return &c, nil
}

// Save atomically writes the config back to disk with 0600 perms.
// Atomic = write to a temp file + rename, so a crash mid-write can never
// produce a half-corrupt config.
func Save(c *Config) error {
	c.Version = ConfigVersion
	path := Path()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	tmp, err := os.CreateTemp(filepath.Dir(path), ".config.json.*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName) // no-op if rename succeeded

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}

func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if 'A' <= ca && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if 'A' <= cb && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}
