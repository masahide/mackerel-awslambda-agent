export AWS_DEFAULT_REGION ?= ap-northeast-1

NAME:= $(notdir $(PWD))

DIST=.dist
CHECKS_BIN=$(DIST)/checks/checks
INVOKER_BIN=$(DIST)/invoker/invoker
SENDER_BIN=$(DIST)/sender/sender
BINS=$(CHECKS_BIN) $(INVOKER_BIN) $(SENDER_BIN)
TEST_OPTIONS?=
PKG?= $(wildcard pkg/*)
ENV  := dev

CF_STACKNAME := $(ENV)-$(NAME)
CF_FILE      := template.yml
IAM_CF_FILE  := iam-template.yml

PLUGIN_DIST := $(dir $(CHECKS_BIN))
PLUGIN_FILES := $(PLUGIN_DIST)/check-aws-cloudwatch-logs-insights

export GO111MODULE := on

# Install all the build and lint dependencies
setup:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh
	go mod download
	go get github.com/golang/mock/mockgen@v1.4.4
	go generate ./...
.PHONY: setup


# gofmt and goimports all go files
fmt:
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done
.PHONY: fmt

test:
	#go test $(TEST_OPTIONS) -v -race -coverpkg=$(PKG) \
	#	-covermode=atomic -coverprofile=pkg_coverage.txt \
	#	$(PKG)  -run . -timeout=2m
	go test -race -cover -covermode=atomic -coverprofile=pkg_coverage.txt -timeout=2m ./pkg/...
	cat pkg_coverage.txt > coverage.txt
.PHONY: test

cover: test
	go tool cover -html=coverage.txt
	rm coverage.txt
.PHONY: cover

# Run all the linters
lint:
	./bin/golangci-lint run --tests=false --enable-all --disable=lll,wsl ./...
.PHONY: lint

# Run all the tests and code checks
ci: build test lint
.PHONY: ci

# Build a beta version of $(BINS)
build: clean $(BINS) extract-plugins
.PHONY: build

clean:
	rm -rf $(DIST)
	rm -f coverage.txt *coverage.txt
.PHONY: clean

$(INVOKER_BIN):
	GOOS=linux GOARCH=amd64 go build -o $@ cmd/$(notdir $@)/main.go
$(SENDER_BIN):
	GOOS=linux GOARCH=amd64 go build -o $@ cmd/$(notdir $@)/main.go
$(CHECKS_BIN):
	GOOS=linux GOARCH=amd64 go build -o $@ cmd/$(notdir $@)/main.go

.DEFAULT_GOAL := build

.PHONY: extract-plugins
extract-plugins: $(PLUGIN_FILES)

$(PLUGIN_FILES):
	mkdir -p $(PLUGIN_DIST)
	DIST=$(PLUGIN_DIST) ./hack/get_plugins.sh

