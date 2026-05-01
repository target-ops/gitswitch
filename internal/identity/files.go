package identity

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// IdentitiesDir returns the directory where per-identity gitconfig
// files live. Mirrors the layout decision in Path() — honours
// $XDG_CONFIG_HOME.
func IdentitiesDir() string {
	if x := os.Getenv("XDG_CONFIG_HOME"); x != "" {
		return filepath.Join(x, "gitswitch", "identities")
	}
	return filepath.Join(os.Getenv("HOME"), ".config", "gitswitch", "identities")
}

// GitconfigPath returns the path of the per-identity gitconfig file
// for `name`.
func GitconfigPath(name string) string {
	return filepath.Join(IdentitiesDir(), name+".gitconfig")
}

// WriteGitconfig writes the per-identity gitconfig for the given
// identity. The file is included by ~/.gitconfig via includeIf when
// the user is inside the bound directory. SSH commit signing is
// enabled by default when SigningKey is set; .core.sshCommand is set
// only when SSHKey is non-empty (so users without a key still get a
// usable identity).
func WriteGitconfig(id Identity) error {
	if err := os.MkdirAll(IdentitiesDir(), 0o700); err != nil {
		return err
	}
	path := GitconfigPath(id.Name)

	var sb strings.Builder
	sb.WriteString("# gitswitch-managed file. Do not edit by hand.\n")
	sb.WriteString("# Regenerated whenever `gitswitch use " + id.Name + "` runs.\n\n")

	sb.WriteString("[user]\n")
	if id.GitName != "" {
		fmt.Fprintf(&sb, "    name = %s\n", id.GitName)
	}
	if id.Email != "" {
		fmt.Fprintf(&sb, "    email = %s\n", id.Email)
	}
	if id.SigningKey != "" {
		fmt.Fprintf(&sb, "    signingkey = %s\n", id.SigningKey)
	}

	if id.SigningKey != "" {
		sb.WriteString("\n[commit]\n    gpgsign = true\n")
		sb.WriteString("\n[tag]\n    gpgsign = true\n")
		sb.WriteString("\n[gpg]\n    format = ssh\n")
	}

	if id.SSHKey != "" {
		sb.WriteString("\n[core]\n")
		// Important: IdentitiesOnly=yes prevents ssh-agent from offering
		// every other key in the keychain to GitHub, which is both a
		// privacy leak and the source of the famous "Permission denied
		// (publickey)" failure when you have multiple accounts.
		fmt.Fprintf(&sb, "    sshCommand = ssh -i %s -o IdentitiesOnly=yes\n", id.SSHKey)
	}

	return os.WriteFile(path, []byte(sb.String()), 0o600)
}
