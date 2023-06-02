.DEFAULT_GOAL := target

VERSION=$(shell git describe --tags --always)
FEATURE=$(shell git branch | grep '*' | awk '{ print $$2; }')
BUILD_TIME=$(shell date +%Y%m%d%H%M%S)
VERSION_PATH=github.com/machinefi/w3bstream/cmd/srv-applet-mgr/types
BUILD_ROOT=../../build/pub_client

.PHONY: debug
debug:
	@echo ${FEATURE}
	@echo ${VERSION}
	@echo ${BUILD_TIME}
	@echo ${BUILD_ROOT}

.PHONY: build
build: clean
	@go build -ldflags "-X ${VERSION_PATH}.Name=pub_client\
 -X ${VERSION_PATH}.Feature=${FEATURE}\
 -X ${VERSION_PATH}.Version=${VERSION}\
 -X ${VERSION_PATH}.Timestamp=${BUILD_TIME}"
	@mkdir -p ${BUILD_ROOT}
	@mv pub_client ${BUILD_ROOT}/pub_client

.PHONY: clean
clean:
	@rm -rf ${BUILD_ROOT}/pub_client
	@rm -rf ${BUILD_ROOT}/config

.PHONY: target
target: build