# Set the shell to bash always
SHELL := /bin/bash

# Go installs the above commands to get installed in $GOBIN if defined, and $GOPATH/bin otherwise:
#GOBIN := $(shell go env GOBIN)
#GOBIN ?= $(shell go env GOPATH)/bin

VCS_HOST=github.com
ORG_NAME := lucasepe
PROJECT_NAME := using-client-go/using-codegen

UNAME := $(uname -s)

.PHONY: help
help:	### Show targets documentation
ifeq ($(UNAME), Linux)
	@grep -P '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
else
	@awk -F ':.*###' '$$0 ~ FS {printf "%15s%s\n", $$1 ":", $$2}' \
		$(MAKEFILE_LIST) | grep -v '@awk' | sort
endif

.PHONY: vendor
vendor: ### Vendor dependencies
	@go mod vendor

.PHONY: deps
deps:	### Optimize dependencies
	@go mod tidy

.PHONY: generate
generate: vendor ### Generate code
	@bash ./hack/hack.sh

.PHONY: generate2
generate2: ### Generate code
	set -e;\
	CODEGEN_DIR=./code-generator;\
	MODULE=$(VCS_HOST)/$(ORG_NAME)/$(PROJECT_NAME);\
	$${CODEGEN_DIR}/generate-groups.sh "deepcopy,client" \
	    $$MODULE/pkg/generated $$MODULE/pkg/apis \
	    expression:v1alpha1 \
	    -h ./hack/boilerplate.go.txt \
	    -o ../using-codegen;\
	cp -r $$MODULE/pkg/* ./pkg;\
	rm -rf $(VCS_HOST)

.PHONY: install
install: ## Install CRDs
	@kubectl apply -f manifests/crds

.PHONY: clean
clean: ### Clean build files
	@rm -rf $(VCS_HOST) ./pkg/generated
	@rm -rf ./bin
	@go clean

.PHONY: build
build: generate clean ### Build binary
	@go build -tags netgo -a -v -ldflags "${LD_FLAGS}" -o ./bin/expressions-operator ./cmd/*.go
	@chmod +x ./bin/*

