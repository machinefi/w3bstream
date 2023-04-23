.DEFAULT_GOAL := build

REMOTE=$(shell git remote -v | head -n 1 | awk -F" " '{ print $$2 }')
BRANCH=$(shell git branch | grep '*' | awk '{ print $$2; }')
COMMIT_ID=$(shell git rev-parse --short HEAD)
COMMIT_AUTHOR=$(shell git log | head -n 2 | grep Author | awk -F "[<>]" '{ print $$2;  }')
TAG=$(shell git describe --tags `git rev-list --tags --max-count=1`)
NOW=$(shell date +%Y%m%d%H%M%S)
BUILD_PATH=github.com/machinefi/w3bstream/cmd/srv-applet-mgr/types


.PHONY: build
build:
	@go build -ldflags "-X ${BUILD_PATH}.Name=srv-applet-mgr\
 -X ${BUILD_PATH}.Remote=${REMOTE}\
 -X ${BUILD_PATH}.Branch=${BRANCH}\
 -X ${BUILD_PATH}.Commit=${COMMIT_ID}\
 -X ${BUILD_PATH}.Version=${TAG}\
 -X ${BUILD_PATH}.Timestamp=${NOW}"
	@rm -rf build
	@mkdir -p build
	@mv srv-applet-mgr build/
	@cp -r config build/
	@rm -rf openapi.json
	@go generate .
	@cp openapi.json build/
	@rm -rf ../../build/srv-applet-mgr
	@mkdir -p ../../build
	@mv build ../../build/srv-applet-mgr

.PHONY: run
run: build
	@cd ../../build/srv-applet-mgr && ./srv-applet-mgr

.PHONY: light
light:
	@go build -ldflags "-s -w -X ./global.ServiceInfo=${REMOTE}..${BRANCH}..${COMMIT_ID}..${COMMIT_AUTHOR}" -o srv-applet-mgr
	@rm -rf build
	@mkdir -p build
	@mv srv-applet-mgr build/
	@cp -r config build/
	@cp openapi.json build/
	@rm -rf ../../build/srv-applet-mgr
	@mkdir -p ../../build && mv build ../../build/srv-applet-mgr
