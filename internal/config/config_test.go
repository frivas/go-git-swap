package config

import (
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/frivas/go-git-swap/internal/model"
)

func TestProfileLifecycle(t *testing.T) {
	path := filepath.Join(t.TempDir(), FileName)
	cfg, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"work", "personal", "client"} {
		if err := cfg.AddProfile(model.GitProfile{Name: name, Email: name + "@x.com"}); err != nil {
			t.Fatal(err)
		}
	}
	if err := cfg.AddProfile(model.GitProfile{Name: "work"}); err == nil {
		t.Fatal("duplicate profile accepted")
	}
	if got := cfg.ListProfiles(); !slices.Equal(got, []string{"client", "personal", "work"}) {
		t.Fatalf("ListProfiles = %v, want sorted", got)
	}
	if err := cfg.SetActiveProfile("work"); err != nil {
		t.Fatal(err)
	}
	if err := cfg.RemoveProfile("work"); err != nil {
		t.Fatal(err)
	}

	// Reload from disk: active profile was cleared with its profile.
	cfg, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.GetActiveProfile() != "" {
		t.Fatalf("active = %q after deleting it", cfg.GetActiveProfile())
	}
	if len(cfg.ListProfiles()) != 2 {
		t.Fatalf("profiles = %v", cfg.ListProfiles())
	}
}

func TestFileIsPrivate(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix permissions do not apply")
	}
	path := filepath.Join(t.TempDir(), FileName)
	if err := os.WriteFile(path, []byte(`{"profiles":{}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := cfg.AddProfile(model.GitProfile{Name: "a"}); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Fatalf("mode = %o, want 600", perm)
	}
}

func TestSaveKeepsSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "dotfiles.json")
	link := filepath.Join(dir, FileName)
	if err := os.WriteFile(target, []byte(`{"profiles":{}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	cfg, err := Open(link)
	if err != nil {
		t.Fatal(err)
	}
	if err := cfg.AddProfile(model.GitProfile{Name: "a"}); err != nil {
		t.Fatal(err)
	}
	if fi, err := os.Lstat(link); err != nil || fi.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("symlink replaced by a regular file: %v", err)
	}
	if data, _ := os.ReadFile(target); !strings.Contains(string(data), `"a"`) {
		t.Fatalf("target not updated: %s", data)
	}
}
