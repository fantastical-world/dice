# Makefile for dice package
BUILD_COMMIT := $(shell git describe --tags --always --dirty --match=v*)
BUILD_DATE := $(shell date -u +%b-%d-%Y,%T-UTC)
BUILD_SEMVER := $(shell cat .SEMVER)

.PHONY: all build clean help release test dirty-check

# target: all - default target, will trigger build
all: build

# target: build - runs build for local os/arch
build:
	go build ./...

# target: clean - removes artifacts from tests
clean:
	-rm -rf results

# target: release - will clean, build, test, and finally creates a git tag for the version
release: dirty-check clean build test
	git tag v$(BUILD_SEMVER) $(BUILD_COMMIT)
	git push origin v$(BUILD_SEMVER)

# target: test - runs tests and generates coverage reports
test:
	mkdir -p results
	go test -tags test ./... -cover -coverprofile=results/tc.out
	go tool cover -html=results/tc.out -o results/coverage.html

# target: test-race - runs tests with race checks (slower) and generates coverage reports
test-race:
	mkdir -p results
	go test -race -tags test ./... -cover -coverprofile=results/tc.out
	go tool cover -html=results/tc.out -o results/coverage.html

# target: dirty-check - will check if repo is dirty
dirty-check:
ifneq (, $(findstring dirty, $(BUILD_COMMIT)))
	@echo "you're dirty check your repo status before releasing"
	false
endif

# target: help - displays help
help:
	@egrep "^#.?target:" Makefile
