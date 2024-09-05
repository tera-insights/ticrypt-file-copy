# Copyright 2024 Tera Insights, LLC. All rights reserved.

## Do not modify; this is maintained by a Github Action. 
VERSION := $(shell grep 'current_version =' .bumpversion.toml | sed 's/^[[:space:]]*current_version = //'| tr -d '"')
VERSION_WO_DASHES = $(subst -,.,$(VERSION))

GOLANG_PKG := github.com/tera-insights/ticrypt-file-copy
GOTAGS_TEST := -gcflags="all=-N -l" -tags no_openssl
GOTAGS := -tags no_openssl
GO_LDFLAGS = -ldflags '-X "main.Version=$(VERSION)" -X "main.BuildDate=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")" -X "main.GitHash=$(shell git rev-parse HEAD)" $(GO_EXTRA_LDFLAGS)'

## Building
clean:
	rm -rf bin

bin/: clean
	mkdir -p bin/

bin/debug/: clean
	mkdir -p bin/debug/

bin/ticrypt-file-copy: bin/
	go build $(GO_LDFLAGS) $(GOTAGS) -o bin/ticp $(GOLANG_PKG)

bin/debug/ticrypt-file-copy: bin/debug/
	go build $(GO_LDFLAGS) $(GOTAGS_TEST) -o bin/debug/ticp $(GOLANG_PKG)

build: bin/ticrypt-file-copy

build_debug: bin/debug/ticrypt-file-copy

## Installation

install: build
	cp bin/ticp /usr/local/bin/ticp
	mkdir -p /etc/ticp
	cp config/ticp.conf /etc/ticp/ticp.conf

## Testing
test: mocks
	go test ./...