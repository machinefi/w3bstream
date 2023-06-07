.DEFAULT_GOAL := target

VERSION=$(shell git describe --tags --always)
FEATURE=$(shell git branch | grep '*' | awk '{ print $$2; }')
BUILD_TIME=$(shell date +%Y%m%d%H%M%S)
VERSION_PATH=github.com/machinefi/w3bstream/cmd/srv-applet-mgr/types
BUILD_ROOT=../../build/srv-applet-mgr
REPO_ROOT=../..

.PHONY: debug
debug:
	@echo ${FEATURE}
	@echo ${VERSION}
	@echo ${BUILD_TIME}
	@echo ${BUILD_ROOT}

.PHONY: toolkit
toolkit:
	@go install github.com/machinefi/w3bstream/pkg/depends/gen/cmd/...@toolkit-patch-0.0.3
	@echo installed `which toolkit`

.PHONY: openapi
openapi: toolkit
	@go generate .

.PHONY: build
build: clean
	@go build -ldflags "-X ${VERSION_PATH}.Name=srv-applet-mgr\
 -X ${VERSION_PATH}.Feature=${FEATURE}\
 -X ${VERSION_PATH}.Version=${VERSION}\
 -X ${VERSION_PATH}.Timestamp=${BUILD_TIME}"
	@mkdir -p ${BUILD_ROOT}
	@mv srv-applet-mgr  ${BUILD_ROOT}/srv-applet-mgr
	@cp -r openapi.json ${BUILD_ROOT}/openapi.json

.PHONY: target
target: build

.PHONY: run
run: build
	@cd ${BUILD_ROOT} && ./srv-applet-mgr

.PHONY: clean
clean:
	@rm -rf ${BUILD_ROOT}/srv-applet-mgr
	@rm -rf ${BUILD_ROOT}/openapi.json

.PHONY: image
image:
	@cd ${REPO_ROOT} && docker build -f cmd/srv-applet-mgr/Dockerfile . -t $(USER)/w3bstream:${VERSION}
