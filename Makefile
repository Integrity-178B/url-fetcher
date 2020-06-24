GOOS?=darwin
GOARCH?=amd64
APP?=api

VERSION=`git describe --abbrev=7 --always --tags`

.PHONY: help build clean clean-all lint

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build: clean ## Build application or tool. Envs: GOOS, GOARCH, APP
	@echo "Building ${GOOS}-${GOARCH}/${APP}..."
	@GOSUMDB=off GOOS=${GOOS} GOARCH=${GOARCH} go build -o build/${GOOS}-${GOARCH}/${APP} -ldflags "-X main.revision=${REVISION} -s -w" cmd/$(APP)/main.go
	@echo "Done"

clean: ## Clean build directory of application or tool. Envs: GOOS, GOARCH, APP
	@echo "Cleaning ${GOOS}-${GOARCH}/${APP} in build directory..."
	@rm -rf build/${GOOS}-${GOARCH}/${APP}
	@echo "Done"

clean-all: ## Clean build directory
	@echo "Cleaning everything in build directory..."
	@rm -rf build/*
	@echo "Done"

lint: ## Lint all project files
	@echo "Running linter..."
	@golangci-lint run --timeout=60s --verbose
	@echo "Done"
