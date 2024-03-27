include go.mk

# Override the variables in the go.mk
APPROOT=karbour
GOSOURCE_PATHS = ./pkg/...
LICENSE_CHECKER ?= license-eye
LICENSE_CHECKER_VERSION ?= main

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
.PHONY: test
test:  ## Run the tests
	go test -gcflags=all=-l -timeout=10m `go list $(GOSOURCE_PATHS) | grep -v "internalimport"` ${TEST_FLAGS}

.PHONY: cover
cover:  ## Generates coverage report
	go test -gcflags=all=-l -timeout=10m `go list $(GOSOURCE_PATHS) | grep -v "internalimport"` -coverprofile $(COVERAGEOUT) ${TEST_FLAGS}


# Target: update-codegen
# Description: Updates the generated code using the 'hack/update-codegen.sh' script.
# Example: make update-codegen
.PHONY: update-codegen
update-codegen: ## Update generated code
	hack/update-codegen.sh

# Build-related targets

# Target: build-all
# Description: Builds for all supported platforms (Darwin, Linux, Windows).
# Example: make build-all
.PHONY: build-all
build-all: build-darwin build-linux build-windows ## Build for all platforms

# Target: build-darwin
# Description: Builds for macOS platform.
# Dependencies: BUILD_UI
# Example:
# - make build-darwin
# - make build-darwin SKIP_UI_BUILD=true
.PHONY: build-darwin
build-darwin: $(BUILD_UI) ## Build for MacOS (Darwin)
	-rm -rf ./_build/darwin
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 \
		go build -o ./_build/darwin/$(APPROOT) \
		./cmd

# Target: build-linux
# Description: Builds for Linux platform.
# Dependencies: BUILD_UI
# Example: make build-linux
.PHONY: build-linux
build-linux: $(BUILD_UI) ## Build for Linux
	-rm -rf ./_build/linux
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
		go build -o ./_build/linux/$(APPROOT) \
		./cmd

# Target: build-windows
# Description: Builds for Windows platform.
# Dependencies: BUILD_UI
# Example: make build-windows
.PHONY: build-windows
build-windows: $(BUILD_UI) ## Build for Windows
	-rm -rf ./_build/windows
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
		go build -o ./_build/windows/$(APPROOT).exe \
		./cmd

# Target: build-server-all
# Description: Builds server for all supported platforms (Darwin, Linux, Windows).
# Example: make build-server-all
.PHONY: build-server-all
build-server-all: build-server-darwin build-server-linux build-server-windows ## Build server for all platforms

# Target: build-server-darwin
# Description: Builds server for the macOS platform.
# Example: make build-server-darwin
.PHONY: build-server-darwin
build-server-darwin: ## Build server for MacOS (Darwin)
	-rm -rf ./_build/darwin
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 \
		go build -o ./_build/darwin/$(APPROOT) \
		./cmd

# Target: build-server-linux
# Description: Builds server for the Linux platform.
# Example: make build-server-linux
.PHONY: build-server-linux
build-server-linux: ## Build server for Linux
	-rm -rf ./_build/linux
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
		go build -o ./_build/linux/$(APPROOT) \
		./cmd

# Target: build-server-windows
# Description: Builds server for the Windows platform.
# Example: make build-server-windows
.PHONY: build-server-windows
build-server-windows: ## Build server for Windows
	-rm -rf ./_build/windows
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
		go build -o ./_build/windows/$(APPROOT).exe \
		./cmd

# Target: build-ui
# Description: Builds the UI for the dashboard.
# Example: make build-ui
.PHONY: build-ui
build-ui: ## Build UI for the dashboard
	@echo "Building UI for the dashboard ..."
	cd ui && npm install && npm run build

.PHONY: check-license
check-license:  ## Checks if repo files contain valid license header
	@which $(LICENSE_CHECKER) > /dev/null || (echo "Installing $(LICENSE_CHECKER)@$(LICENSE_CHECKER_VERSION) ..."; go install github.com/apache/skywalking-eyes/cmd/$(LICENSE_CHECKER)@$(LICENSE_CHECKER_VERSION) && echo -e "Installation complete!\n")
	@${GOPATH}/bin/$(LICENSE_CHECKER) header check

.PHONY: fix-license
fix-license:  ## Adds missing license header to repo files
	@which $(LICENSE_CHECKER) > /dev/null || (echo "Installing $(LICENSE_CHECKER)@$(LICENSE_CHECKER_VERSION) ..."; go install github.com/apache/skywalking-eyes/cmd/$(LICENSE_CHECKER)@$(LICENSE_CHECKER_VERSION) && echo -e "Installation complete!\n")
	@${GOPATH}/bin/$(LICENSE_CHECKER) header fix

.PHONY: gen-api-spec
gen-api-spec: ## Generate API Specification with OpenAPI format
	@which swag > /dev/null || (echo "Installing swag@v1.7.8 ..."; go install github.com/swaggo/swag/cmd/swag@v1.7.8 && echo "Installation complete!\n")
	# Generate API documentation with OpenAPI format
	-swag init --parseDependency --parseInternal --parseDepth 1 -g cmd/main.go -o api/openapispec/
	# Format swagger comments
	-swag fmt -g pkg/**/*.go
	@echo "🎉 Done!"

.PHONY: gen-api-doc
gen-api-doc: ## Generate API Documentation by API Specification
	@which swagger > /dev/null || (echo "Installing swagger@v0.30.5 ..."; go install github.com/go-swagger/go-swagger/cmd/swagger@v0.30.5 && echo "Installation complete!\n")
	-swagger generate markdown -f ./api/openapispec/swagger.json && echo "🎉 Done!" || (echo "💥 Failed!"; exit 1)
