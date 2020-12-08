export AWS_DEFAULT_REGION ?= ap-northeast-1

DIST=.dist
CHECKER_BIN=$(DIST)/checker/checker
INVOKER_BIN=$(DIST)/invoker/invoker
SENDER_BIN=$(DIST)/sender/sender
BINS=$(CHECKER_BIN) $(INVOKER_BIN) $(SENDER_BIN)
PKG?= $(wildcard pkg/*)
ENV  := dev

PLUGIN_DIST := $(dir $(CHECKER_BIN))
PLUGIN_FILES := $(PLUGIN_DIST)/check-aws-cloudwatch-logs-insights

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
	go test -race -cover -covermode=atomic -coverprofile=pkg_coverage.txt -timeout=2m ./pkg/...
	cat pkg_coverage.txt > coverage.txt
.PHONY: test

cover: test
	go tool cover -html=coverage.txt
	rm coverage.txt
.PHONY: cover

# Run all the linters
lint:
	./bin/golangci-lint run --tests=false --enable-all --disable=lll,wsl,exhaustivestruct,nlreturn,gochecknoglobals ./...
.PHONY: lint

# Run all the tests and code checker
ci: build test lint
.PHONY: ci

# Build a beta version of $(BINS)
build: clean $(BINS) extract-plugins
.PHONY: build

clean:
	rm -rf $(BINS)
	rm -f coverage.txt *coverage.txt
.PHONY: clean

$(INVOKER_BIN):
	GOOS=linux GOARCH=amd64 go build -o $@ cmd/$(notdir $@)/main.go
$(SENDER_BIN):
	GOOS=linux GOARCH=amd64 go build -o $@ cmd/$(notdir $@)/main.go
$(CHECKER_BIN):
	GOOS=linux GOARCH=amd64 go build -o $@ cmd/$(notdir $@)/main.go

.DEFAULT_GOAL := build

.PHONY: extract-plugins
extract-plugins: $(PLUGIN_FILES)

$(PLUGIN_FILES):
	mkdir -p $(PLUGIN_DIST)
	DIST=$(PLUGIN_DIST) ./hack/get_plugins.sh

