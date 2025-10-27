
.PHONY: help image tests lint

BINARY_NAME=anniedad
GOLANG_CI_VERSION=v1.59.1
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
CGO=0

# grep the version from the mix file
VERSION=$(shell ./version.sh)

#default: tests

help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

build: ## Build app
	CGO_ENABLED=$(CGO) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -installsuffix cgo -o $(BINARY_NAME) ./cmd/server

# TOOLS

lint: ## Run linters
	golangci-lint run -v ./...

install_lint: ## Get GLOLANGCI_LINT and install
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANG_CI_VERSION)

cert: ## generate TLC certs
	cd cert; ./gen.sh; cd ..

tests: ## Run local tests
	go test -v -coverprofile=cover.txt `go list ./... | egrep -v 'proto|mock'`

coverage: ## Check coverage
	go tool cover -func cover.txt

test-cover: tests coverage ## Run local tests with coverage checking

proto: ## Generate grpc files
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./internal/proto/anniedad/anniedad.proto ./internal/proto/auth/auth.proto

# VERSIONS

bump-patch: tests ## Bump patch version
	bumpversion patch

bump-minor: tests ## Bump minor version
	bumpversion minor

bump-major: tests ## Bump major version
	bumpversion major

version: ## Output the current version
	@echo $(VERSION)
