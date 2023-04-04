include go.mk

# Override the variables in the go.mk
APPROOT=karbour
GOSOURCE_PATHS = ./pkg/...


## Build-related targets
.PHONY: build-all
build-all: build-darwin build-linux build-windows ## Build for all platforms

.PHONY: build-darwin
build-darwin: ## Build for MacOS
	-rm -rf ./_build/darwin
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 \
		go build -o ./_build/darwin/$(APPROOT) \
		./cmd

.PHONY: build-linux
build-linux: ## Build for Linux
	-rm -rf ./_build/linux
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
		go build -o ./_build/linux/$(APPROOT) \
		./cmd

.PHONY: build-windows
build-windows: ## Build for Windows
	-rm -rf ./_build/windows
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
		go build -o ./_build/windows/$(APPROOT).exe \
		./cmd

