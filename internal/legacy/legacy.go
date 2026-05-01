// Package legacy reads the v0.2.x Python config at ~/.config.ini and
// converts it into the v1.0 identity model.
//
// The on-disk format from v0.2.x is a stock Python configparser file:
//
//	[github]
//	ofirhaim1 = ofir474@gmail.com,/Users/x/.ssh/id_rsa_github_OfirHaim1
//
//	[gitlab]
//	work = work@example.com,/Users/x/.ssh/id_rsa_gitlab_work
//
//	[current]
//	vendor = github
//	username = ofirhaim1
//
// Each section header is a vendor name. Each key=value entry inside is
// `<username> = <email>,<key_path>`. The [current] section is just
// metadata about which user was last "switched to" — it doesn't carry
// over into v1.0 (the new model is directory-bound, not active-user).
//
// We avoid pulling in an INI library for this: the format is small,
// the parser fits in 30 lines, and it eliminates a transitive
// dependency for what is essentially one-shot migration code.
package legacy

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/target-ops/gitswitch/internal/identity"
)

// LegacyPath returns the conventional location of the v0.2.x config.
// It lived in the user's home directory (not under ~/.config) because
// the Python implementation hard-coded that path.
func LegacyPath() string {
	return filepath.Join(os.Getenv("HOME"), ".config.ini")
}

// Exists reports whether a legacy config file is present and readable.
func Exists() bool {
	_, err := os.Stat(LegacyPath())
	return err == nil
}

// Parse reads ~/.config.ini and returns the contained identities.
// Returns (nil, nil) if the file doesn't exist — that's the
// "fresh-install user, nothing to migrate" path and shouldn't be
// treated as an error.
func Parse() ([]identity.Identity, error) {
	data, err := os.ReadFile(LegacyPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var (
		out          []identity.Identity
		section      string
		currentVendor string // tracks which vendor [section] we're inside
	)

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		// [section] line
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.TrimSpace(line[1 : len(line)-1])
			currentVendor = section
			continue
		}

		// Skip the [current] section entirely — it has different
		// semantics ("vendor", "username") that don't fit the
		// vendor=section / user=key model.
		if section == "current" || section == "DEFAULT" {
			continue
		}

		// key = value
		eq := strings.IndexByte(line, '=')
		if eq < 0 {
			continue
		}
		username := strings.TrimSpace(line[:eq])
		value := strings.TrimSpace(line[eq+1:])
		if username == "" || value == "" {
			continue
		}

		// Value is "email,key_path"
		parts := strings.SplitN(value, ",", 2)
		email := strings.TrimSpace(parts[0])
		keyPath := ""
		if len(parts) == 2 {
			keyPath = strings.TrimSpace(parts[1])
		}

		id := identity.Identity{
			Name:    sanitizeLegacyName(username),
			Email:   email,
			Vendor:  currentVendor,
			SSHKey:  keyPath,
		}
		// For github vendor specifically, the v0.2.x username is also
		// the gh account name (that's how `gh auth switch` matched it).
		if currentVendor == "github" {
			id.GHAccount = username
		}
		// Default the .pub path next to the private key — it's how
		// gitswitch and ssh-keygen both write it. Still empty if the
		// .pub file isn't there; downstream callers handle that.
		if keyPath != "" {
			pub := keyPath + ".pub"
			if _, err := os.Stat(pub); err == nil {
				id.SigningKey = pub
			}
		}
		out = append(out, id)
	}

	return out, scanner.Err()
}

// sanitizeLegacyName strips characters that shouldn't appear in an
// identity name. The Python version was lenient with case; the new
// JSON store treats names case-sensitively, so we lowercase to keep
// behaviour aligned with what `init` produces today.
func sanitizeLegacyName(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case 'a' <= c && c <= 'z', '0' <= c && c <= '9', c == '-', c == '_':
			out = append(out, c)
		}
	}
	return string(out)
}
