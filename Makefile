GOSOURCE_PATHS ?= ./cmd/...

include go.mk

APPROOT=karbour


.PHONY: clean
clean:  ## Clean build bundles
	-rm -rf ./build

.PHONY: update-codegen
update-codegen: ## Update generated code
	hack/update-codegen.sh

.PHONY: build-all
build-all: build-darwin build-linux build-windows ## Build for all platforms

.PHONY: build-darwin
build-darwin: update-codegen ## Build for MacOS
	-rm -rf ./build/darwin
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 \
		go build -o ./build/darwin/$(APPROOT) \
		./cmd/main.go

.PHONY: build-linux
build-linux: update-codegen ## Build for Linux
	-rm -rf ./build/linux
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
		go build -o ./build/linux/$(APPROOT) \
		./cmd/main.go

.PHONY: build-windows
build-windows: update-codegen ## Build for Windows
	-rm -rf ./build/windows
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
		go build -o ./build/windows/$(APPROOT).exe \
		./cmd/main.go
