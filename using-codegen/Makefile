# Set the shell to bash always
SHELL := /bin/bash

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


.PHONY: install
install: ## Install CRDs
	@kubectl apply -f manifests/crds

.PHONY: clean
clean: ### Clean build files
	@go clean

