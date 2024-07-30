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

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/KusionStack/karpor/pkg/util/gitutil"
	"github.com/KusionStack/karpor/pkg/version"
)

func calculateVersion() (string, error) {
	latestTag, err := gitutil.GetLatestTag()
	if err != nil {
		fmt.Println("Error getting latest tag:", err)
		return "v0.1.0", nil
	}
	currentCommit, err := gitutil.GetCurrentCommit()
	if err != nil {
		fmt.Println("Error getting current commit:", err)
		return "", err
	}
	currentHead, err := gitutil.GetCurrentHead()
	if err != nil {
		fmt.Println("Error getting current head:", err)
		return "", err
	}
	tagCommit, err := gitutil.GetTagCommit(latestTag)
	if err != nil {
		fmt.Println("Error getting tag commit:", err)
		return "", err
	}
	if currentHead == tagCommit {
		return latestTag, nil
	}
	return fmt.Sprintf("%s-%s", latestTag, currentCommit), nil
}

func main() {
	currentDir, err := os.Getwd()
	versionStr, err := calculateVersion()
	if err != nil {
		fmt.Println("Error calculating version:", err)
		os.Exit(1)
	}
	version := version.Version{Version: versionStr}

	rootDir := strings.Replace(currentDir, "/pkg/version/scripts", "", 1)
	versionFilePath := filepath.Join(rootDir, "pkg", "version", "VERSION")

	// Ensure the pkg/version directory exists
	err = os.MkdirAll(filepath.Dir(versionFilePath), 0o755)
	if err != nil {
		fmt.Println("Error creating version directory:", err)
		os.Exit(1)
	}
	// Write the version to the VERSION file, replacing any existing content
	err = os.WriteFile(versionFilePath, []byte(version.Version+"\n"), 0o644)
	if err != nil {
		fmt.Println("Error writing version file:", err)
		os.Exit(1)
	}
	fmt.Println("Version file updated to:", version.Version)
}
