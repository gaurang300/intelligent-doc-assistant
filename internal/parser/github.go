package parser

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CloneGitHubRepo clones a GitHub repository to a temporary directory and returns the path
func (p *Parser) CloneGitHubRepo(repoURL string) (string, error) {
	// Extract repo name from URL
	parts := strings.Split(repoURL, "/")
	repoName := parts[len(parts)-1]
	if strings.HasSuffix(repoName, ".git") {
		repoName = repoName[:len(repoName)-4]
	}

	// Create temp directory
	tempDir := filepath.Join(os.TempDir(), "intelligent-doc-assistant", repoName)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Clone the repository
	cmd := exec.Command("git", "clone", repoURL, tempDir)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to clone repository: %w", err)
	}

	return tempDir, nil
}

// IsGitHubURL checks if a string is a GitHub repository URL
func IsGitHubURL(path string) bool {
	return strings.HasPrefix(path, "https://github.com/") ||
		strings.HasPrefix(path, "git@github.com:")
}
