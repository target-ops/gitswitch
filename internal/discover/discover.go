// Package discover inspects the machine for existing git/ssh/gh state
// and returns a list of plausible identities. It never writes anything.
//
// The output of Scan() feeds `gitswitch init`, which then asks the user
// to name and confirm them.
package discover

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/target-ops/gitswitch/internal/gh"
	"github.com/target-ops/gitswitch/internal/git"
	"github.com/target-ops/gitswitch/internal/ssh"
)

// Detected is one identity that was found on the machine. The same
// underlying identity may show up in multiple sources — Scan() coalesces
// them by email when possible.
type Detected struct {
	Email      string
	GitName    string   // user.name from git config
	SSHKey     string   // private key path
	SigningKey string   // .pub used for SSH commit signing
	GHAccount  string
	Vendor     string   // github / gitlab / bitbucket / ""
	Sources    []string // human-readable provenance: ["git config", "~/.ssh/config", "gh", "id_rsa.pub"]
}

// Scan runs all the detectors and returns the coalesced result.
// Errors from individual detectors are swallowed (and the source is
// simply absent) — partial discovery is much more useful than no
// discovery, and most of these failures are "the file doesn't exist."
func Scan() []Detected {
	// Each detector emits one or more fragment Detecteds. We coalesce
	// after all fragments are collected — that way a fragment from
	// source N can be merged with a fragment from source N+M even if
	// they couldn't have been merged at the time they were emitted.
	var fragments []Detected

	addSource := func(d *Detected, s string) {
		for _, existing := range d.Sources {
			if existing == s {
				return
			}
		}
		d.Sources = append(d.Sources, s)
	}

	// 1. git config global
	gitName, gitEmail, _ := git.GlobalIdentity()
	signingKey, _ := git.GlobalSigning()
	if gitEmail != "" {
		d := Detected{
			Email:      gitEmail,
			GitName:    gitName,
			SigningKey: signingKey,
		}
		addSource(&d, "git config")
		fragments = append(fragments, d)
	}

	// 2. ~/.ssh/config — host blocks tell us which key is bound to which vendor
	blocks, _ := ssh.ParseConfig()
	for _, b := range blocks {
		if b.IdentityFile == "" {
			continue
		}
		vendor := vendorFromHost(b.HostName, b.Host)
		if vendor == "" {
			continue
		}
		d := Detected{SSHKey: b.IdentityFile, Vendor: vendor}
		addSource(&d, "~/.ssh/config")
		fragments = append(fragments, d)
	}

	// 3. ~/.ssh/*.pub — extract email comments
	for _, pub := range listPublicKeys() {
		email := commentEmail(pub.contents)
		if email == "" {
			continue
		}
		// Strip the .pub suffix to get the private-key path.
		priv := strings.TrimSuffix(pub.path, ".pub")
		d := Detected{Email: email, SSHKey: priv, SigningKey: pub.path}
		addSource(&d, filepath.Base(pub.path))
		fragments = append(fragments, d)
	}

	// 4. gh — active login
	if login, _ := gh.ActiveLogin(); login != "" {
		d := Detected{GHAccount: login, Vendor: "github"}
		addSource(&d, "gh")
		fragments = append(fragments, d)
	}

	return distributeWeak(coalesce(fragments))
}

// distributeWeak handles the case where a fragment carries only a gh
// account (or only a vendor + ssh key) and can't merge by strict-match
// rules — but the machine has exactly one strong identity with the
// same vendor. In that case the weak fragment is almost certainly
// describing the same person.
//
// Concretely: gh login octocat (vendor=github, no email/key) gets
// distributed onto a single id_rsa-bearing identity that also targets
// github.com. If two strong github identities exist, we leave the weak
// fragment standalone — better to ask the user than guess wrong.
func distributeWeak(in []Detected) []Detected {
	isStrong := func(d Detected) bool {
		// "strong" = carries enough to anchor an identity on its own.
		return d.Email != "" || d.SSHKey != ""
	}
	var weak []int
	strongByVendor := map[string][]int{}
	for i, d := range in {
		if !isStrong(d) {
			weak = append(weak, i)
			continue
		}
		v := d.Vendor
		if v == "" {
			continue
		}
		strongByVendor[v] = append(strongByVendor[v], i)
	}
	mergedInto := map[int]bool{}
	for _, wi := range weak {
		w := in[wi]
		if w.Vendor == "" {
			continue
		}
		candidates := strongByVendor[w.Vendor]
		if len(candidates) != 1 {
			continue
		}
		target := candidates[0]
		mergeFields(&in[target], w)
		for _, s := range w.Sources {
			in[target].Sources = appendUnique(in[target].Sources, s)
		}
		mergedInto[wi] = true
	}
	if len(mergedInto) == 0 {
		return in
	}
	out := make([]Detected, 0, len(in)-len(mergedInto))
	for i, d := range in {
		if !mergedInto[i] {
			out = append(out, d)
		}
	}
	return out
}

// coalesce union-finds the fragments: any two that share a matching
// field (email, canonical SSH-key path, or gh account) collapse into
// one. We iterate until no more merges happen — single-pass would miss
// cases where fragment A merges with B via field X, then B should also
// merge with C via field Y newly inherited from A.
func coalesce(in []Detected) []Detected {
	out := append([]Detected(nil), in...)
	for {
		merged := false
		for i := 0; i < len(out); i++ {
			for j := i + 1; j < len(out); j++ {
				if shouldMerge(&out[i], &out[j]) {
					mergeFields(&out[i], out[j])
					for _, s := range out[j].Sources {
						out[i].Sources = appendUnique(out[i].Sources, s)
					}
					out = append(out[:j], out[j+1:]...)
					merged = true
					break
				}
			}
			if merged {
				break
			}
		}
		if !merged {
			break
		}
	}
	return out
}

func shouldMerge(a, b *Detected) bool {
	return matchEmail(a, *b) || matchSSHKey(a, *b) || matchGH(a, *b)
}

func appendUnique(slice []string, s string) []string {
	for _, x := range slice {
		if x == s {
			return slice
		}
	}
	return append(slice, s)
}

func matchEmail(a *Detected, b Detected) bool {
	return a.Email != "" && b.Email != "" && strings.EqualFold(a.Email, b.Email)
}

func matchSSHKey(a *Detected, b Detected) bool {
	if a.SSHKey == "" || b.SSHKey == "" {
		return false
	}
	if a.SSHKey == b.SSHKey {
		return true
	}
	// Symlink resolution: id_rsa_github_work might be a symlink to
	// id_rsa, and we want those to coalesce.
	ar, errA := filepath.EvalSymlinks(a.SSHKey)
	br, errB := filepath.EvalSymlinks(b.SSHKey)
	return errA == nil && errB == nil && ar == br
}

func matchGH(a *Detected, b Detected) bool {
	return a.GHAccount != "" && b.GHAccount != "" && a.GHAccount == b.GHAccount
}

// mergeFields copies fields from `b` into `a` for any field that's empty
// on `a`. Doesn't overwrite — first-detected wins for conflicting values.
func mergeFields(a *Detected, b Detected) {
	if a.Email == "" {
		a.Email = b.Email
	}
	if a.GitName == "" {
		a.GitName = b.GitName
	}
	if a.SSHKey == "" {
		a.SSHKey = b.SSHKey
	}
	if a.SigningKey == "" {
		a.SigningKey = b.SigningKey
	}
	if a.GHAccount == "" {
		a.GHAccount = b.GHAccount
	}
	if a.Vendor == "" {
		a.Vendor = b.Vendor
	}
}

// vendorFromHost returns the canonical vendor name for known git hosts,
// or "" if we don't recognize it.
func vendorFromHost(hostname, alias string) string {
	for _, s := range []string{hostname, alias} {
		s = strings.ToLower(s)
		switch {
		case strings.Contains(s, "github.com"):
			return "github"
		case strings.Contains(s, "gitlab.com"):
			return "gitlab"
		case strings.Contains(s, "bitbucket.org"):
			return "bitbucket"
		}
	}
	return ""
}

type publicKeyFile struct {
	path     string
	contents string
}

// listPublicKeys reads every *.pub file in ~/.ssh/. Returns an empty
// slice on any error — discovery is best-effort.
func listPublicKeys() []publicKeyFile {
	dir := filepath.Join(os.Getenv("HOME"), ".ssh")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []publicKeyFile
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".pub") {
			continue
		}
		full := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(full)
		if err != nil {
			continue
		}
		out = append(out, publicKeyFile{path: full, contents: string(data)})
	}
	return out
}

// commentEmail extracts the email-looking trailing comment from an SSH
// public key file. Format: "<algo> <base64> <comment>" where the comment
// is conventionally the email used during ssh-keygen.
func commentEmail(contents string) string {
	scanner := bufio.NewScanner(strings.NewReader(contents))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		comment := fields[len(fields)-1]
		if strings.Contains(comment, "@") && strings.Contains(comment, ".") {
			return comment
		}
	}
	return ""
}
