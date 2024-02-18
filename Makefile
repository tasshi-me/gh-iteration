SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := help

GO_FILES:=$(shell find . -type f -name '*.go' -print)

BIN:=gh-iteration
BIN_EXT:=
ifeq ($(shell go env GOOS), windows)
BIN_EXT:=.exe
endif

DOCS_DIR:=docs

.PHONY: all build lint fix test docs-gen docs-lint docs-update help
all: lint test build docs-gen docs-lint

build: $(BIN)$(BIN_EXT)  ## Build executable file
	@:

$(BIN)$(BIN_EXT): $(GO_FILES) go.mod go.sum
	@go build -o $@ ./cmd/gh-iteration

lint: ## Lint codes
	@golangci-lint run --tests ./...

fix: ## Fix lint errors
	@golangci-lint run --tests ./... --fix

test: ## Run tests
	@go test $(option) ./...

docs-gen: ## Generate documentation files
	@go run ./cmd/gen-docs --out-dir $(DOCS_DIR)

docs-lint: docs-gen ## Check if the docs outdated
	@if test -n "$(shell git status $(DOCS_DIR) -s)"; then \
		echo "$(DOCS_DIR) is outdated"; \
		git status $(DOCS_DIR) -s; \
		exit 1; \
	fi

docs-update: docs-gen ## Update and commit document
	git add docs
	git commit -m "docs: update document"

help: ## Display this help
	# https://postd.cc/auto-documented-makefile/
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
