package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/frivas/go-git-swap/internal/model"
)

// sandbox points HOME at a temp dir so the real ~/.gitconfig is never touched.
func sandbox(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(home, ".gitconfig"))
	return home
}

func get(t *testing.T, key string) string {
	t.Helper()
	v, _ := New().GetGlobalConfig(key)
	return v
}

func TestSwitchingDropsPreviousSigningKey(t *testing.T) {
	home := sandbox(t)
	sshDir := filepath.Join(home, ".ssh")
	if err := os.MkdirAll(sshDir, 0o700); err != nil {
		t.Fatal(err)
	}
	keyPath := filepath.Join(sshDir, "work")
	if out, err := exec.Command("ssh-keygen", "-q", "-t", "ed25519", "-N", "", "-f", keyPath).CombinedOutput(); err != nil {
		t.Skipf("ssh-keygen unavailable: %v %s", err, out)
	}
	signers := filepath.Join(sshDir, "allowed_signers")
	if err := os.WriteFile(signers, []byte("someone@else ssh-ed25519 AAAAexisting"), 0o600); err != nil {
		t.Fatal(err)
	}

	g := New()
	work := model.GitProfile{Name: "work", FullName: "Work Me", Username: "wk", Email: "me@work.com", SigningKey: "~/.ssh/work.pub"}
	for range 2 { // second run must not duplicate the signer line
		if err := g.SetGlobalConfig(work); err != nil {
			t.Fatal(err)
		}
	}
	if get(t, "user.signingkey") != keyPath+".pub" || get(t, "commit.gpgsign") != "true" || get(t, "gpg.format") != "ssh" {
		t.Fatalf("signing not configured: key=%q", get(t, "user.signingkey"))
	}
	data, _ := os.ReadFile(signers)
	pub, _ := ReadPublicKey(keyPath + ".pub")
	if !strings.HasPrefix(string(data), "someone@else ") || strings.Count(string(data), pub) != 1 {
		t.Fatalf("allowed_signers clobbered or duplicated:\n%s", data)
	}

	// A real signed commit verifies against allowed_signers.
	repo := filepath.Join(home, "repo")
	for _, args := range [][]string{{"init", "-q", repo}, {"-C", repo, "commit", "-q", "--allow-empty", "-m", "x"}, {"-C", repo, "verify-commit", "HEAD"}} {
		if out, err := exec.Command("git", args...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}

	if err := g.SetGlobalConfig(model.GitProfile{Name: "personal", FullName: "Me", Username: "me", Email: "me@home.org"}); err != nil {
		t.Fatal(err)
	}
	if k := get(t, "user.signingkey"); k != "" {
		t.Fatalf("work signing key leaked into personal profile: %q", k)
	}
	if get(t, "commit.gpgsign") != "" || get(t, "user.email") != "me@home.org" {
		t.Fatal("personal profile not applied")
	}
}

func TestReadPublicKeyRejectsPrivateKey(t *testing.T) {
	p := filepath.Join(t.TempDir(), "id")
	if err := os.WriteFile(p, []byte("-----BEGIN OPENSSH PRIVATE KEY-----\nabc\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadPublicKey(p); err == nil {
		t.Fatal("private key accepted as public key")
	}
}
