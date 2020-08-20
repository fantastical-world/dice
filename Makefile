# Makefile for dicebag package
BRANCH_NAME := $(shell git branch 2> /dev/null | sed -e '/^[^*]/d' -e 's/* \(.*\)/(\1)/')
BUILD_COMMIT := $(shell git describe --tags --always --dirty --all --match=v*)
BUILD_DATE := $(shell date -u +%b-%d-%Y,%T-UTC)
BUILD_SEMVER := $(shell cat .SEMVER)

.PHONY: all build clean help install release test dirty-check

# target: all - default target, will trigger build
all: build

# target: build - runs build for local os/arch
build:
	CGO_ENABLED=0 go build -ldflags "-X main.buildBranch=$(BRANCH_NAME) -X main.buildCommit=$(BUILD_COMMIT) -X main.buildDate=$(BUILD_DATE) -X main.semVer=$(BUILD_SEMVER)" .

# target: clean - removes artifacts from tests, build, and install
clean:
	-rm -rf test
	go clean -i

# target: install - builds and installs package for local os/arch
install:
	CGO_ENABLED=0 go install -ldflags "-X main.buildBranch=$(BRANCH_NAME) -X main.buildCommit=$(BUILD_COMMIT) -X main.buildDate=$(BUILD_DATE) -X main.semVer=$(BUILD_SEMVER)" .

# target: release - will clean, build, install, and finally creates a git tag for the version
release: dirty-check clean test install
	git tag -a v$(BUILD_SEMVER) $(BUILD_COMMIT) -m "dicebag release - v$(BUILD_SEMVER)"
	git push

# target: test - runs tests and generates coverage reports
test:
	mkdir -p test
	go test -cover -coverprofile=test/c.out
	go tool cover -html=test/c.out -o test/coverage.html

# target: dirty-check - will check if repo is dirty
dirty-check:
ifneq (, $(findstring dirty, $(BUILD_COMMIT)))
	@echo "you're dirty check your repo status before releasing"
	false
endif

# target: help - displays help
help:
	@egrep "^#.?target:" Makefile
