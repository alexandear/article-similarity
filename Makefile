MAKEFILE_PATH := $(abspath $(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
PATH := $(MAKEFILE_PATH):$(PATH)

export GOBIN := $(MAKEFILE_PATH)/bin

all: clean format swagger build test

help:
	@echo 'Usage: make <TARGETS> ... <OPTIONS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@echo '    build              Compile packages and dependencies.'
	@echo '    clean              Remove binary.'
	@echo '    format             Run gofmt on package sources.'
	@echo '    generate           Generate swagger code, mocks and other code.'
	@echo '    help               Show this help screen.'
	@echo '    swagger            Generate only swagger code.'
	@echo '    test               Run tests.'
	@echo ''
	@echo 'Targets run by default are: clean format swagger build test'
	@echo ''

clean:
	@echo clean
	@go clean

build:
	@echo build
	@go build -o $(GOBIN)/article-similarity

test:
	@echo test
	@go test -count=1 -v ./...

lint:
	@echo lint
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint
	@$(GOBIN)/golangci-lint run

format:
	@echo format
	@go fmt $(PKGS)

SWAGGER          = $(GOBIN)/swagger
SPEC             = $(MAKEFILE_PATH)/api/spec.yaml
SWAGGER_GEN_PATH = $(MAKEFILE_PATH)/internal

doc:
	@echo swagger doc
	@$(SWAGGER) serve --flavor=swagger $(SPEC)

generate: swagger

swagger:
	@echo swagger
	@go install github.com/go-swagger/go-swagger/cmd/swagger
	@rm -rf $(SWAGGER_GEN_PATH)/models
	@rm -rf $(SWAGGER_GEN_PATH)/restapi
	@$(SWAGGER) generate server -f $(SPEC) -t $(SWAGGER_GEN_PATH) --exclude-main --flag-strategy pflag

