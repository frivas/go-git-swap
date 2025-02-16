package git

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

type GitConfig struct{}

func New() *GitConfig {
    return &GitConfig{}
}

func (g *GitConfig) SetGlobalConfig(email, username, name, signingKey string) error {
    configs := map[string]string{
        "user.email":                    email,
        "user.name":                     name,
        "user.username":                 username,
        "commit.gpgsign":                "true",
        "gpg.format":                    "ssh",
        "user.signingkey":               signingKey,
        "gpg.ssh.allowedSignersFile":    "~/.ssh/allowed_signers",
    }

    homeDir, err := os.UserHomeDir()
    if err != nil {
        return fmt.Errorf("failed to get home directory: %w", err)
    }

    for key, value := range configs {
        if value != "" {
            // Expand ~ to home directory for paths
            if strings.HasPrefix(value, "~/") {
                value = filepath.Join(homeDir, value[2:])
            }
            cmd := exec.Command("git", "config", "--global", key, value)
            if err := cmd.Run(); err != nil {
                return fmt.Errorf("failed to set %s: %w", key, err)
            }
        }
    }

    // Create or update allowed_signers file
    if signingKey != "" {
        allowedSignersPath := filepath.Join(homeDir, ".ssh", "allowed_signers")
        allowedSignersDir := filepath.Dir(allowedSignersPath)

        if err := os.MkdirAll(allowedSignersDir, 0700); err != nil {
            return fmt.Errorf("failed to create .ssh directory: %w", err)
        }

        allowedSignersContent := fmt.Sprintf("%s %s", email, signingKey)
        if err := os.WriteFile(allowedSignersPath, []byte(allowedSignersContent), 0600); err != nil {
            return fmt.Errorf("failed to create allowed_signers file: %w", err)
        }
    }

    return nil
}

func (g *GitConfig) GetGlobalConfig(key string) (string, error) {
    cmd := exec.Command("git", "config", "--global", "--get", key)
    output, err := cmd.Output()
    if err != nil {
        return "", err
    }
    return string(output), nil
}
