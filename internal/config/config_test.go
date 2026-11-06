package config

import (
	"os"
	"path/filepath"
	"slices"
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
