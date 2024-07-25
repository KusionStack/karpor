package gitutil

import (
	"bytes"
	"os/exec"
	"strings"
)

// RunGitCommand runs a git command and returns the output as a string
func runGitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

// GetLatestTag returns the latest tag from Git
func GetLatestTag() (string, error) {
	return runGitCommand("describe", "--tags", "--abbrev=0")
}

// GetCurrentCommit returns the current commit SHA
func GetCurrentCommit() (string, error) {
	return runGitCommand("rev-parse", "--short", "HEAD")
}

// GetCurrentHead returns the current HEAD SHA
func GetCurrentHead() (string, error) {
	return runGitCommand("rev-parse", "HEAD")
}

// GetTagCommit returns the commit SHA for a given tag
func GetTagCommit(tag string) (string, error) {
	return runGitCommand("rev-list", "-n", "1", tag)
}
