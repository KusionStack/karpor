// Copyright The Karbour Authors.
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

package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/golang/glog"
)

const (
	defaultStaticDir = "dist"
)

type StaticDirHandler struct {
	staticDir string
}

// NewStaticDirHandler loads the localization configuration and constructs a LocaleHandler.
func NewStaticDirHandler(staticDir string) *StaticDirHandler {
	if len(staticDir) == 0 {
		staticDir = defaultStaticDir
	}
	fmt.Println(staticDir)

	return &StaticDirHandler{
		staticDir: staticDir,
	}
}

// LocaleHandler serves different html versions based on the Accept-Language header.
func (handler *StaticDirHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello")
	// if r.URL.EscapedPath() == "/" || r.URL.EscapedPath() == "/index.html" {
	// 	// Do not store the html page in the cache. If the user is to click on 'switch language',
	// 	// we want a different index.html (for the right locale) to be served when the page refreshes.
	// 	w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
	// }

	// Disable directory listing.
	// if r.URL.Path != "/" && strings.HasSuffix(r.URL.Path, "/") {
	// 	http.NotFound(w, r)
	// 	return
	// }

	// Get assets directory
	// dirName := getStaticDirDir()
	dirName := handler.staticDir
	fmt.Println(dirName)

	http.FileServer(http.Dir(dirName)).ServeHTTP(w, r)
}

// getStaticDirDir determines the absolute path to the localized frontend assets
func getStaticDirDir() string {
	path, err := os.Executable()
	if err != nil {
		glog.Fatalf("Error determining path to executable: %#v", err)
	}

	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		glog.Fatalf("Error evaluating symlinks for path '%s': %#v", path, err)
	}

	return filepath.Join(filepath.Dir(path), defaultStaticDir)
}
