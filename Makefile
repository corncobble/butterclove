# Local dependencies: mkdir, awk, echo, printf, rm, go

OUT = butterclove
.DEFAULT_GOAL := help

##@ Building
.PHONY: all build
all: tidy build

build: ## Build for production
	CGO_ENABLED=0 go build -ldflags='-s' -o $(OUT)

##@ Development
.PHONY: dev-build
dev-build: ## Build for development
	go build -race -gcflags='all=-N -l' -o $(OUT)


##@ Cleanup
.PHONY: clean
clean: ## Remove all build and download artifacts
	@echo "Clearing build..."
	@rm -rf $(OUT)


##@ Helpers
.PHONY: tidy help
tidy: ## Tidy up the go.mod file
	@go mod tidy

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
