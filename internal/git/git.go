package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/frivas/go-git-swap/internal/model"
)

type GitConfig struct{}

func New() *GitConfig {
	return &GitConfig{}
}

// ExpandPath expands $VARS and a leading ~/ in path.
func ExpandPath(path string) (string, error) {
	path = os.ExpandEnv(path)
	if rest, ok := strings.CutPrefix(path, "~/"); ok {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, rest)
	}
	return path, nil
}

// ReadPublicKey returns "<type> <base64>" from an OpenSSH public key file.
// It rejects anything else, including private keys.
func ReadPublicKey(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	fields := strings.Fields(string(data))
	isKeyType := func(p string) bool { return len(fields) > 0 && strings.HasPrefix(fields[0], p) }
	if len(fields) < 2 || !slices.ContainsFunc([]string{"ssh-", "ecdsa-", "sk-"}, isKeyType) {
		return "", fmt.Errorf("%s is not an OpenSSH public key (expected the .pub file)", path)
	}
	return fields[0] + " " + fields[1], nil
}

// SetGlobalConfig points the global git identity at profile. A profile without
// a signing key turns signing off, so a previous profile's key is never reused.
func (g *GitConfig) SetGlobalConfig(p model.GitProfile) error {
	for key, value := range map[string]string{
		"user.email":    p.Email,
		"user.name":     p.FullName,
		"user.username": p.Username,
	} {
		if value == "" {
			if err := unset(key); err != nil {
				return err
			}
			continue
		}
		if err := set(key, value); err != nil {
			return err
		}
	}

	if p.SigningKey == "" {
		for _, key := range []string{"user.signingkey", "commit.gpgsign"} {
			if err := unset(key); err != nil {
				return err
			}
		}
		return nil
	}

	keyPath, err := ExpandPath(p.SigningKey)
	if err != nil {
		return err
	}
	pubKey, err := ReadPublicKey(keyPath)
	if err != nil {
		return err
	}

	signersPath, err := g.GetGlobalConfig("gpg.ssh.allowedSignersFile")
	if err != nil || signersPath == "" {
		signersPath = "~/.ssh/allowed_signers"
	}
	if signersPath, err = ExpandPath(signersPath); err != nil {
		return err
	}
	if err := addAllowedSigner(signersPath, p.Email, pubKey); err != nil {
		return err
	}

	for _, kv := range [][2]string{
		{"gpg.format", "ssh"},
		{"user.signingkey", keyPath},
		{"commit.gpgsign", "true"},
		{"gpg.ssh.allowedSignersFile", signersPath},
	} {
		if err := set(kv[0], kv[1]); err != nil {
			return err
		}
	}
	return nil
}

// addAllowedSigner appends "<email> <key>" to the allowed signers file unless
// it is already there. Existing entries are kept.
func addAllowedSigner(path, email, pubKey string) error {
	line := email + ` namespaces="git" ` + pubKey
	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("read allowed signers: %w", err)
	}
	if slices.Contains(strings.Split(string(data), "\n"), line) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create %s: %w", filepath.Dir(path), err)
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("open allowed signers: %w", err)
	}
	if len(data) > 0 && !strings.HasSuffix(string(data), "\n") {
		line = "\n" + line
	}
	_, err = f.WriteString(line + "\n")
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	return err
}

func set(key, value string) error {
	if out, err := exec.Command("git", "config", "--global", key, value).CombinedOutput(); err != nil {
		return fmt.Errorf("set %s: %w: %s", key, err, strings.TrimSpace(string(out)))
	}
	return nil
}

func unset(key string) error {
	err := exec.Command("git", "config", "--global", "--unset", key).Run()
	// Exit status 5: the key was not set.
	if exitErr, ok := errors.AsType[*exec.ExitError](err); ok && exitErr.ExitCode() == 5 {
		return nil
	}
	if err != nil {
		return fmt.Errorf("unset %s: %w", key, err)
	}
	return nil
}

func (g *GitConfig) GetGlobalConfig(key string) (string, error) {
	output, err := exec.Command("git", "config", "--global", "--get", key).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
