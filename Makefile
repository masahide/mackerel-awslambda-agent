export AWS_DEFAULT_REGION ?= ap-northeast-1

NAME:= $(notdir $(PWD))

DIST=.dist
BIN?=$(DIST)/mackerel-awslambda-agent
MAIN?=./cmd/mackerel-awslambda-agent
TEST_PATTERN?=.
TEST_OPTIONS?=
PKG?=./pkg/config
ENV  := dev

CF_STACKNAME := $(ENV)-$(NAME)
CF_FILE      := template.yml
IAM_CF_FILE  := iam-template.yml


PKG_CF_FILE := .package_template.$(ENV).yml
S3BUCKET    := test-yamasaki-masahide
S3PREFIX    := packages/AWS-SAM/$(ENV)-$(NAME)
S3SRCPREFIX := $(S3PREFIX)/src

# mackerel plugins
PLUGINS     := check-aws-cloudwatch-logs check-aws-sqs-queue-size

PLUGIN_DIST := $(PWD)/$(DIST)
PLUGIN_FILES := $(PLUGINS:%=$(DIST)/%)

export GO111MODULE := on

# Install all the build and lint dependencies
setup:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh
	go mod download
.PHONY: setup


# gofmt and goimports all go files
fmt:
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done
.PHONY: fmt

test:
	go test $(TEST_OPTIONS) -v -race -coverpkg=$(MAIN) \
		-covermode=atomic -coverprofile=main_coverage.txt \
		$(MAIN) -run $(TEST_PATTERN) -timeout=2m
	go test $(TEST_OPTIONS) -v -race -coverpkg=$(PKG) \
		-covermode=atomic -coverprofile=pkg_coverage.txt \
		$(PKG)  -run $(TEST_PATTERN) -timeout=2m
	cat main_coverage.txt pkg_coverage.txt > coverage.txt
.PHONY: test

cover: test
	go tool cover -html=coverage.txt
	rm coverage.txt
.PHONY: cover

# Run all the linters
lint:
	./bin/golangci-lint run --tests=false --enable-all --disable=lll ./...
.PHONY: lint

# Run all the tests and code checks
ci: build test lint
.PHONY: ci

# Build a beta version of $(BIN)
build: clean $(BIN)
.PHONY: build

clean:
	rm -rf $(DIST)
	rm -f coverage.txt *coverage.txt
.PHONY: clean

$(BIN):
	GOOS=linux GOARCH=amd64 go build -o $@ $(MAIN)/main.go

.DEFAULT_GOAL := build






.PHONY: extract-plugins
extract-plugins: $(PLUGIN_FILES)

$(PLUGIN_FILES):
	mkdir -p $(DIST)
	plugins="$(PLUGINS)" \
	plugin_dist=$(PLUGIN_DIST) \
			./hack/get_plugins.sh

.PHONY: package
package: $(PKG_CF_FILE)

$(PKG_CF_FILE): $(BIN) $(PLUGIN_FILES)
	aws cloudformation package \
		--template-file $(CF_FILE) \
		--s3-bucket $(S3BUCKET) \
		--s3-prefix $(S3SRCPREFIX) \
		--output-template-file $(PKG_CF_FILE)

.PHONY: deploy
deploy:
	aws cloudformation deploy \
		--template-file $(PKG_CF_FILE) \
		--stack-name $(CF_STACKNAME) \
		--capabilities CAPABILITY_IAM


.PHONY: create-iam
create-iamrole:
	aws cloudformation create-stack --stack-name $(CF_STACKNAME)-iam \
		--template-body file://$(IAM_CF_FILE) \
		--capabilities CAPABILITY_NAMED_IAM CAPABILITY_AUTO_EXPAND


