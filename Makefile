GOSOURCE_PATHS ?= ./cmd/...

include go.mk


.PHONY: clean
clean:  ## Clean build bundles
	-rm -rf ./build

.PHONY: build-all
build-all: build-darwin build-linux build-windows ## Build for all platforms

.PHONY: build-darwin
build-darwin: ## Build for MacOS
	-rm -rf ./build/darwin
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 \
		go build -o ./build/darwin/$(APPROOT) \
		./cmd/main.go

.PHONY: build-linux
build-linux: ## Build for Linux
	-rm -rf ./build/linux
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
		go build -o ./build/linux/$(APPROOT) \
		./cmd/main.go

.PHONY: build-windows
build-windows: ## Build for Windows
	-rm -rf ./build/windows
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
		go build -o ./build/windows/$(APPROOT).exe \
		./cmd/main.go
