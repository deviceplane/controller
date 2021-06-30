OS = $(shell uname -s | tr '[:upper:]' '[:lower:]')
export CGO_ENABLED ?= 0

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Debug symbols
ifeq (${DEBUG},)
else
GOARGS=-gcflags="all=-N -l"
endif
LDFLAGS="-w -s"

GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

all: build

.PHONY: vendor
vendor: ## Download and vendor all modules
	@go mod vendor

.PHONY: build
build: export GOPRIVATE = github.com/Edgeworx/*
build: GOARGS += -mod=vendor -ldflags $(LDFLAGS)
build: fmt ## Build binary
	go build $(GOARGS) -o bin/controller cmd/main.go

.PHONY: fmt
fmt: ## Run go fmt against code
	@gofmt -s -w $(GOFILES_NOVENDOR)

.PHONY: lint
lint: export GOFLAGS=-mod=vendor
lint: golangci-lint fmt ## Lint the source
	@$(GOLANGCI_LINT) run --timeout 5m0s

golangci-lint: ## Install golangci
ifeq (, $(shell which golangci-lint))
	@{ \
	set -e ;\
	GOLANGCI_TMP_DIR=$$(mktemp -d) ;\
	cd $$GOLANGCI_TMP_DIR ;\
	go mod init tmp ;\
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.33.0 ;\
	rm -rf $$GOLANGCI_TMP_DIR ;\
	}
GOLANGCI_LINT=$(GOBIN)/golangci-lint
else
GOLANGCI_LINT=$(shell which golangci-lint)
endif

help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
