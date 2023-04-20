.DEFAULT_GOAL := build

REMOTE=$(shell git remote -v | head -n 1 | awk -F" " '{ print $$2 }')
BRANCH=$(shell git branch | grep '*' | awk '{ print $$2; }')
COMMIT_ID=$(shell git rev-parse --short HEAD)
COMMIT_AUTHOR=$(shell git log | head -n 2 | grep Author | awk -F "[<>]" '{ print $$2;  }')

.PHONY: build
build:
	@go build -ldflags "-s -w -X ./global.ServiceInfo=${REMOTE}..${BRANCH}..${COMMIT_ID}..${COMMIT_AUTHOR}" -o pub_client
	@rm -rf build
	@mkdir -p build
	@mv pub_client build/
	@cp -r config build/
	@rm -rf ../../build/pub_client
	@mkdir -p ../../build
	@mv build ../../build/pub_client

