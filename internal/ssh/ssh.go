// Package ssh probes ssh and parses ~/.ssh/config blocks.
package ssh

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// HostBlock is a parsed Host stanza from ~/.ssh/config — only the
// fields gitswitch cares about today.
type HostBlock struct {
	Host         string // value after "Host "
	HostName     string
	IdentityFile string
}

// ParseConfig reads ~/.ssh/config and returns one HostBlock per Host stanza.
// Missing file → ([]HostBlock{}, nil) — that's a normal "fresh machine" state.
func ParseConfig() ([]HostBlock, error) {
	path := filepath.Join(os.Getenv("HOME"), ".ssh", "config")
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var blocks []HostBlock
	var cur *HostBlock

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := splitKV(line)
		if !ok {
			continue
		}
		switch strings.ToLower(key) {
		case "host":
			if cur != nil {
				blocks = append(blocks, *cur)
			}
			cur = &HostBlock{Host: val}
		case "hostname":
			if cur != nil {
				cur.HostName = val
			}
		case "identityfile":
			if cur != nil {
				cur.IdentityFile = expandTilde(val)
			}
		}
	}
	if cur != nil {
		blocks = append(blocks, *cur)
	}
	return blocks, scanner.Err()
}

// TestAuth runs `ssh -T <host>` and returns the welcome message line
// if the server greeted us (e.g. GitHub's "Hi <user>!"). Empty string
// + nil error when the auth was rejected with no message; non-nil error
// only on shell-level failures.
func TestAuth(host string, timeout time.Duration) (string, error) {
	ctx, cancel := contextTimeout(timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ssh", "-T",
		"-o", "BatchMode=yes",
		"-o", "StrictHostKeyChecking=accept-new",
		host)
	out, _ := cmd.CombinedOutput() // ssh -T to git hosts always exits non-zero
	return firstWelcomeLine(string(out)), nil
}

var welcomeRE = regexp.MustCompile(`(?m)^(?:Hi|Welcome|Hello)[^\n]*$`)

func firstWelcomeLine(s string) string {
	m := welcomeRE.FindString(s)
	return strings.TrimSpace(m)
}

func splitKV(line string) (string, string, bool) {
	for i, r := range line {
		if r == ' ' || r == '\t' || r == '=' {
			key := strings.TrimSpace(line[:i])
			val := strings.TrimSpace(line[i+1:])
			val = strings.TrimLeft(val, " \t=")
			return key, val, key != "" && val != ""
		}
	}
	return "", "", false
}

func expandTilde(p string) string {
	if strings.HasPrefix(p, "~/") {
		return filepath.Join(os.Getenv("HOME"), p[2:])
	}
	return p
}
