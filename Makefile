include go.mk

# Override the variables in the go.mk
APPROOT=karpor
GOSOURCE_PATHS = ./pkg/... ./cmd/...
LICENSE_CHECKER ?= license-eye
LICENSE_CHECKER_VERSION ?= main

# Front-End tools
UIFORMATER			?= prettier

# Default architecture for building binaries.
# Override this variable by setting GOARCH=<your-architecture> before invoking the make command.
# To find this list of possible platforms, run the following:
#   go tool dist list
GOARCH ?= amd64

# Default setting for CGO_ENABLED to disable the use of cgo.
# Can be overridden by setting CGO_ENABLED=1 before invoking the make command.
CGO_ENABLED ?= 0


# Check if the SKIP_UI_BUILD flag is set to control the UI building process.
# If the flag is not set, the BUILD_UI variable is assigned the value 'build-ui'.
# If the flag is set, the BUILD_UI variable remains empty.
ifndef SKIP_UI_BUILD
    BUILD_UI = build-ui
else
    BUILD_UI =
endif

# If you encounter an error like "panic: permission denied" on MacOS,
# please visit https://github.com/eisenxp/macos-golink-wrapper to find the solution.

.PHONY: gen-version
gen-version: ## Generate version file
	@echo "🛠️ Updating the version file ..."
	@cd pkg/version/scripts && $(GO) run gen/gen.go

.PHONY: test
test:  ## Run the tests
	@PKG_LIST=$${TARGET_PKG:-$(GOSOURCE_PATHS)}; \
	$(GO) test -gcflags=all=-l -timeout=10m `$(GO) list -e $${PKG_LIST} | grep -vE "cmd|internal|internalimport|generated|handler|middleware|registry|openapi|apis|version|gitutil|server|elasticsearch"` ${TEST_FLAGS}


# cover: Generates a coverage report for the specified TARGET_PKG or default GOSOURCE_PATHS.
# Usage:
#   make cover TARGET_PKG=<go-package-path>
# Example:
#   make cover                              # use the default GOSOURCE_PATHS
#   make cover TARGET_PKG='./pkg/util/...'  # specify a custom package path
.PHONY: cover
cover: ## Generates coverage report
	@PKG_LIST=$${TARGET_PKG:-$(GOSOURCE_PATHS)}; \
	echo "🚀 Executing unit tests for $${PKG_LIST}:"; \
	$(GO) test -gcflags=all=-l -timeout=10m `$(GO) list $${PKG_LIST} | grep -vE "cmd|internal|internalimport|generated|handler|middleware|registry|openapi|apis|version|gitutil|server|elasticsearch"` -coverprofile $(COVERAGEOUT) ${TEST_FLAGS} && \
	(echo "\n📊 Calculating coverage rate:"; $(GO) tool cover -func=$(COVERAGEOUT)) || (echo "\n💥 Running go test failed!"; exit 1)


.PHONY: format
format:  ## Format source code of frontend and backend
	@which $(GOFORMATER) > /dev/null || (echo "Installing $(GOFORMATER)@$(GOFORMATER_VERSION) ..."; $(GO) install mvdan.cc/gofumpt@$(GOFORMATER_VERSION) && echo -e "Installation complete!\n")
	@for path in $(GOSOURCE_PATHS); do $(GOFORMATER) -l -w -e `echo $${path} | cut -b 3- | rev | cut -b 5- | rev`; done;
	@which $(UIFORMATER) > /dev/null || (echo "Installing $(UIFORMATER) ..."; npm install -g prettier && echo -e "Installation complete!\n")
	@cd ui && npx prettier --write .


# Target: update-codegen
# Description: Updates the generated code using the 'hack/update-codegen.sh' script.
# Usage: make update-codegen
.PHONY: update-codegen
update-codegen: ## Update generated code
	hack/update-codegen.sh

# VERSION file handling targets
# These targets are used to manage the VERSION file during build process.
# save-version: Creates a backup of the current VERSION file
# restore-version: Restores the VERSION file from backup and removes the backup file
VERSION_FILE := pkg/version/VERSION
VERSION_BACKUP := $(VERSION_FILE).bak
.PHONY: save-version
save-version:
	@if [ -f $(VERSION_FILE) ]; then \
		echo "📦 Backing up VERSION file..."; \
		cp $(VERSION_FILE) $(VERSION_BACKUP); \
	fi

.PHONY: restore-version
restore-version:
	@if [ -f $(VERSION_BACKUP) ]; then \
		echo "📦 Restoring VERSION file..."; \
		cp $(VERSION_BACKUP) $(VERSION_FILE); \
		rm $(VERSION_BACKUP); \
	fi

# Build-related targets

# Internal build targets without version handling
# These targets perform the actual build operation for each platform.
# They are prefixed with '_' to indicate they are internal and should not be called directly.
# Each target is responsible for:
# 1. Cleaning the platform-specific build directory
# 2. Building the binary with correct GOOS and GOARCH
.PHONY: _build-darwin
_build-darwin:
	@rm -rf ./_build/darwin
	@echo "🚀 Building karpor-server for darwin platform ..."
	GOOS=darwin GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) \
		$(GO) build -o ./_build/darwin/$(APPROOT) \
		./cmd/karpor || exit 1

.PHONY: _build-linux
_build-linux:
	@rm -rf ./_build/linux
	@echo "🚀 Building karpor-server for linux platform ..."
	GOOS=linux GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) \
		$(GO) build -o ./_build/linux/$(APPROOT) \
		./cmd/karpor || exit 1

.PHONY: _build-windows
_build-windows:
	@rm -rf ./_build/windows
	@echo "🚀 Building karpor-server for windows platform ..."
	GOOS=windows GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) \
		$(GO) build -o ./_build/windows/$(APPROOT).exe \
		./cmd/karpor || exit 1

# Public build targets with version handling
# These targets provide the complete build workflow for each platform:
# 1. Backup the current VERSION file (save-version)
# 2. Generate a new version (gen-version)
# 3. Build the binary (_build-xxx)
# 4. Restore the original VERSION file (restore-version)
# If build fails, the VERSION file will still be restored

# Target: build-darwin
# Description: Builds for macOS platform.
# Usage:
#   make build-darwin GOARCH=<your-architecture> SKIP_UI_BUILD=<true,false>
# Example:
#   make build-darwin
#   make build-darwin GOARCH=arm64
#   make build-darwin GOARCH=arm64 SKIP_UI_BUILD=true
.PHONY: build-darwin
build-darwin: save-version gen-version $(BUILD_UI) ## Build for MacOS (Darwin)
	@$(MAKE) _build-darwin || ($(MAKE) restore-version && exit 1)
	@$(MAKE) restore-version

# Target: build-linux
# Description: Builds for Linux platform.
# Usage:
#   make build-linux GOARCH=<your-architecture> SKIP_UI_BUILD=<true,false>
# Example:
#   make build-linux
#   make build-linux GOARCH=arm64
#   make build-linux GOARCH=arm64 SKIP_UI_BUILD=true
.PHONY: build-linux
build-linux: save-version gen-version $(BUILD_UI) ## Build for Linux
	@$(MAKE) _build-linux || ($(MAKE) restore-version && exit 1)
	@$(MAKE) restore-version

# Target: build-windows
# Description: Builds for Windows platform.
# Usage:
#   make build-windows GOARCH=<your-architecture> SKIP_UI_BUILD=<true,false>
# Example:
#   make build-windows
#   make build-windows GOARCH=arm64
#   make build-windows GOARCH=arm64 SKIP_UI_BUILD=true
.PHONY: build-windows
build-windows: save-version gen-version $(BUILD_UI) ## Build for Windows
	@$(MAKE) _build-windows || ($(MAKE) restore-version && exit 1)
	@$(MAKE) restore-version

# Target: build-ui
# Description: Builds the UI for the dashboard.
# Usage: make build-ui
# TODO: use env var to define the build mode
.PHONY: build-ui
build-ui: gen-version ## Build UI for the dashboard
	@echo "🧀 Building UI for the dashboard ..."
	cd ui && npm install && npm run build && touch build/.gitkeep

# Target: build-all
# Description: Builds for all supported platforms (Darwin, Linux, Windows).
# Note: Uses recursive make calls to ensure each platform build has its own
# version handling context, preventing interference between builds.
# Usage: make build-all
.PHONY: build-all
build-all: ## Build for all platforms
	@echo "🚀 Building for all platforms..."
	@$(MAKE) build-darwin
	@$(MAKE) build-linux
	@$(MAKE) build-windows

# Target: build
# Description: Automatically builds for the current platform.
# Detects the current OS and calls the appropriate platform-specific build target.
# Usage: make build
.PHONY: build
build: ## Build for current platform
	@echo "🔍 Detecting current platform..."
	@case "$$(uname -s)" in \
		Darwin*) \
			echo "🚀 Detected macOS platform, building for darwin..." && \
			$(MAKE) build-darwin ;; \
		Linux*) \
			echo "🚀 Detected Linux platform, building for linux..." && \
			$(MAKE) build-linux ;; \
		MINGW*|MSYS*|CYGWIN*) \
			echo "🚀 Detected Windows platform, building for windows..." && \
			$(MAKE) build-windows ;; \
		*) \
			echo "❌ Unsupported platform: $$(uname -s)" && exit 1 ;; \
	esac

.PHONY: check-license
check-license:  ## Checks if repo files contain valid license header
	@which $(LICENSE_CHECKER) > /dev/null || (echo "Installing $(LICENSE_CHECKER)@$(LICENSE_CHECKER_VERSION) ..."; $(GO) install github.com/apache/skywalking-eyes/cmd/$(LICENSE_CHECKER)@$(LICENSE_CHECKER_VERSION) && echo -e "Installation complete!\n")
	@${GOPATH}/bin/$(LICENSE_CHECKER) header check

.PHONY: fix-license
fix-license:  ## Adds missing license header to repo files
	@which $(LICENSE_CHECKER) > /dev/null || (echo "Installing $(LICENSE_CHECKER)@$(LICENSE_CHECKER_VERSION) ..."; $(GO) install github.com/apache/skywalking-eyes/cmd/$(LICENSE_CHECKER)@$(LICENSE_CHECKER_VERSION) && echo -e "Installation complete!\n")
	@${GOPATH}/bin/$(LICENSE_CHECKER) header fix

.PHONY: gen-api-spec
gen-api-spec: ## Generate API Specification with OpenAPI format
	@which $(GOPATH)/bin/swag > /dev/null || (echo "Installing swag@v1.7.8 ..."; $(GO) install github.com/swaggo/swag/cmd/swag@v1.7.8 && echo "Installation complete!\n")
	# Generate API documentation with OpenAPI format
	@$(GOPATH)/bin/swag init --parseDependency --parseInternal --parseDepth 1 -g cmd/karpor/main.go -o api/openapispec/ && echo "🎉 Done!" || (echo "💥 Fail!"; exit 1)
	# Format swagger comments
	@$(GOPATH)/bin/swag fmt -g pkg/**/*.go && echo "🎉 Done!" || (echo "💥 Failed!"; exit 1)

.PHONY: gen-api-doc
gen-api-doc: gen-api-spec ## Generate API Documentation by API Specification
	@which $(GOPATH)/bin/swagger > /dev/null || (echo "Installing swagger@v0.30.5 ..."; $(GO) install github.com/go-swagger/go-swagger/cmd/swagger@v0.30.5 && echo "Installation complete!\n")
	@$(GOPATH)/bin/swagger generate markdown -f ./api/openapispec/swagger.json --output=docs/api.md && echo "🎉 Done!" || (echo "💥 Fail!"; exit 1)

.PHONY: gen-cli-doc
gen-cli-doc: ## Generate CLI Documentation
	@$(GO) run ./hack/gen-cli-docs/main.go && echo "🎉 Done!"

# Target: add-contributor
# Description: Adds a new contributor to the project's list of contributors using the all-contributors-cli tool.
# Usage:
#   make add-contributor user=<github-username> role=<contributor-roles>
# Example:
#   make add-contributor user=mike role=code
#   make add-contributor user=john role=code,doc
# Where:
#   <github-username> is the GitHub username of the contributor.
#   <contributor-roles> is a comma-separated list of roles the contributor has (e.g., code, doc, design, ideas),
#     with all values listed in the https://allcontributors.org/docs/en/emoji-key.
.PHONY: add-contributor
add-contributor: ## Add a new contributor
	@if [ -z "$(user)" ] || [ -z "$(role)" ]; then \
		echo "Error: 'user' and 'role' must be specified."; \
		echo "Usage: make add-contributor user=<github-username> role=<contributor-roles>"; \
		exit 1; \
	fi
	@which all-contributors > /dev/null || (echo "Installing all-contributors-cli ..."; npm i -g all-contributors-cli && echo -e "Installation complete!\n")
	@all-contributors add $(user) $(role) && echo "🎉 Done!" || (echo "💥 Fail!"; exit 1)

# Target: update-contributors
# Description: Generate the latest list of contributors and update it in README.
# Usage:
#   make update-contributors
.PHONY: update-contributors
update-contributors: ## Update the list of contributors
	@which all-contributors > /dev/null || (echo "Installing all-contributors-cli ..."; npm i -g all-contributors-cli && echo -e "Installation complete!\n")
	-all-contributors generate && echo "🎉 Done!" || (echo "💥 Fail!"; exit 1)

# Target: check
# Description: Run all checks to ensure code quality.
# The checks are run in the following order:
# 1. 🔨 lint: Check code style and potential issues using golangci-lint
# 2. 🧪 cover: Run tests and generate coverage report
# 3. 📦 build: Build the binary for the current platform
# If any check fails, the subsequent checks will not run.
# Usage:
#   make check
.PHONY: check
check: ## Check the lint, test, and build
	@echo "🔍 Running all checks..."
	@echo "🔨 1/3 Running lint check..."
	@$(MAKE) lint || (echo "❌ Lint check failed!" && exit 1)
	@echo "✅ Lint check passed!"
	@echo "🧪 2/3 Running test coverage..."
	@$(MAKE) cover || (echo "❌ Test coverage check failed!" && exit 1)
	@echo "✅ Test coverage check passed!"
	@echo "📦 3/3 Running build check..."
	@$(MAKE) build || (echo "❌ Build check failed!" && exit 1)
	@echo "✅ Build check passed!"
	@echo "🎉 All checks passed successfully!"

# controller-gen path
CONTROLLER_GEN = ${GOPATH}/bin/controller-gen
# controller-gen version
CONTROLLER_GEN_VERSION = v0.17.1

# Target: install-controller-gen
# Description: Install controller_gen.
# Usage:
#   make install-controller-gen
.PHONY: install-controller-gen
install-controller-gen:
	$(GO) install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_GEN_VERSION)

# Target: generate-crds
# Description: Generate CRDs into a special dir.
# Usage:
#   make generate-crds
.PHONY: generate-crds
generate-crds:
	@# generate rbac, webhook and crds
	$(CONTROLLER_GEN) crd  paths="./pkg/kubernetes/apis/cluster/v1beta1/..." output:crd:artifacts:config=config/crds/
	$(CONTROLLER_GEN) crd  paths="./pkg/kubernetes/apis/search/v1beta1/..." output:crd:artifacts:config=config/crds/

# Target: manifests
# Description: Install controller_gen and generate CRDs into a special dir.
# Usage:
#   make manifests
.PHONY: manifests
manifests: install-controller-gen generate-crds
