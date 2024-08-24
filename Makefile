SHELL             := /usr/bin/env bash -Eeu -o pipefail
REPO_ROOT         := $(shell git rev-parse --show-toplevel)
MAKEFILE_DIR      := $(shell { cd "$(subst /,,$(dir $(lastword ${MAKEFILE_LIST})))" && pwd; } || pwd)
DOTLOCAL_DIR      := ${MAKEFILE_DIR}/.local
DOTLOCAL_BIN_DIR  := ${DOTLOCAL_DIR}/bin

export PATH := ${DOTLOCAL_BIN_DIR}:${REPO_ROOT}/.bin:${PATH}

.DEFAULT_GOAL := help
.PHONY: help
help:  ## Display this help documents
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' ${MAKEFILE_LIST} | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'

.PHONY: setup
setup:  ## Setup tools for development
	# == SETUP =====================================================
	# versenv
	make versenv
	# --------------------------------------------------------------

.PHONY: versenv
versenv:
	# direnv
	direnv allow .
	# golangci-lint
	golangci-lint --version

.PHONY: clean
clean:  ## Clean up cache, etc
	# reset tmp
	rm -rf ${MAKEFILE_DIR}/.tmp
	mkdir -p ${MAKEFILE_DIR}/.tmp
	# go build cache
	go env GOCACHE
	go clean -x -cache -testcache -modcache -fuzzcache
	# golangci-lint cache
	golangci-lint cache status
	golangci-lint cache clean

.PHONY: lint
lint:  ## Run secretlint, go mod tidy, golangci-lint
	# typo
	typos
	# gitleaks ref. https://github.com/gitleaks/gitleaks
	gitleaks detect --source . -v
	# tidy
	go-mod-tidy-all
	# golangci-lint
	# ref. https://golangci-lint.run/usage/linters/
	golangci-lint run -c "${REPO_ROOT}/.golangci.yml" --fix --sort-results --verbose --timeout=5m
	# diff
	git diff --exit-code


.PHONY: test
test:  ## Run go test and display coverage
	@[ -x "${DOTLOCAL_BIN_DIR}/godotnev" ] || GOBIN="${DOTLOCAL_BIN_DIR}" go install github.com/joho/godotenv/cmd/godotenv@latest

	# Unit testing
	godotenv -f .test.env go test -v -race -p=4 -parallel=8 -timeout=300s -cover -coverprofile=./coverage.txt.tmp ./... ; grep -Ev "\.deprecated\.go" ./coverage.txt.tmp > ./coverage.txt ; go tool cover -func=./coverage.txt

.PHONY: bench
bench: ## Run benchmarks
	cd integrationtest && go test -run "^NoSuchTestForBenchmark" -benchmem -bench . github.com/kunitsucom/util.go/integrationtest/database/sql -v -trimpath -race -p=4 -parallel=8 -timeout=30s

.PHONY: ci
ci: lint test ## CI command set

.PHONY: up
up:  ## Run docker compose up --wait -d
	# Run in background (If failed to start, output logs and exit abnormally)
	#if ! docker compose up --wait -d; then docker compose logs; exit 1; fi

.PHONY: ps
ps:  ## Run docker compose ps
	#docker compose ps

.PHONY: down
down:  ## Run docker compose down
	#docker compose down

.PHONY: reset
reset:  ## Run docker compose down and Remove volumes
	#docker compose down --volumes

.PHONY: rmi
rmi:  ## Run docker compose down and Remove all images, orphans
	#docker compose down --rmi all --remove-orphans

.PHONY: restart
restart:  ## Restart docker compose
	#-make down
	#make up

.PHONY: logs
logs:  ## Tail docker compose logs
	#@printf '[\033[36mNOTICE\033[0m] %s\n' "If want to go back prompt, enter Ctrl+C"
	#docker compose logs -f

.PHONY: release
release:  ## git tag per go modules for release
	${REPO_ROOT}/.bin/git-tag-go-mod

.PHONY: act-check
act-check:
	@if ! command -v act >/dev/null 2>&1; then \
		printf "\033[31;1m%s\033[0m\n" "act is not installed: brew install act" 1>&2; \
		exit 1; \
	fi

.PHONY: act-go-mod-tidy
act-go-mod-tidy: act-check ## Run go-mod-tidy workflow in act
	act pull_request --container-architecture linux/amd64 -P ubuntu-latest=catthehacker/ubuntu:act-latest -W .github/workflows/go-mod-tidy.yml

.PHONY: act-go-lint
act-go-lint: act-check ## Run go-lint workflow in act
	act pull_request --container-architecture linux/amd64 -P ubuntu-latest=catthehacker/ubuntu:act-latest -W .github/workflows/go-lint.yml

.PHONY: act-go-test
act-go-test: act-check ## Run go-test workflow in act
	act pull_request --container-architecture linux/amd64 -P ubuntu-latest=catthehacker/ubuntu:act-latest -W .github/workflows/go-test.yml

.PHONY: act-go-vuln
act-go-vuln: act-check ## Run go-vuln workflow in act
	act pull_request --container-architecture linux/amd64 -P ubuntu-latest=catthehacker/ubuntu:act-latest -W .github/workflows/go-vuln.yml
