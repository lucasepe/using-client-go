# Set the shell to bash always
SHELL := /bin/bash

# Variables
KIND_CLUSTER_NAME ?= local-dev

# Tools
KIND=$(shell which kind)

# Sets the default goal to be used if no targets were specified on the command line
.DEFAULT_GOAL := help

## help: Print this help
.PHONY: help
help: Makefile
	@echo
	@echo " Choose a command:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## kind.up: Starts the KinD cluster
.PHONY: kind.up
kind.up:
	@$(KIND) get kubeconfig --name $(KIND_CLUSTER_NAME) > /dev/null 2>&1 \
	|| $(KIND) create cluster --name=$(KIND_CLUSTER_NAME)

## kind.down: Shuts down the KinD cluster
.PHONY: kind.down
kind.down:
	@$(KIND) delete cluster --name=$(KIND_CLUSTER_NAME)
