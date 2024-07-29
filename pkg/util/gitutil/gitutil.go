// Copyright The Karpor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


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
