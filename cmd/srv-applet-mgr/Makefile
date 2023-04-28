.DEFAULT_GOAL := build

VERSION=$(shell git describe --tags --always)
FEATURE=$(shell git branch | grep '*' | awk '{ print $$2; }')
BUILD_TIME=$(shell date +%Y%m%d%H%M%S)
VERSION_PATH=github.com/machinefi/w3bstream/cmd/srv-applet-mgr/types

.PHONY: debug
debug:
	@echo ${FEATURE}
	@echo ${VERSION}
	@echo ${BUILD_TIME}

.PHONY: openapi
openapi:
	@go generate .

.PHONY: lite
lite:
	@go build -ldflags "-X ${VERSION_PATH}.ProjectName=srv-applet-mgr\
 -X ${VERSION_PATH}.Feature=${FEATURE}\
 -X ${VERSION_PATH}.Version=${VERSION}\
 -X ${VERSION_PATH}.Timestamp=${BUILD_TIME}"
	@rm -rf ../../build/srv-applet-mgr/srv-applet-mgr
	@rm -rf ../../build/srv-applet-mgr/config
	@rm -rf ../../build/srv-applet-mgr/openapi.json
	@mkdir -p ../../build/srv-applet-mgr
	@mv srv-applet-mgr ../../build/srv-applet-mgr/srv-applet-mgr
	@cp -r config ../../build/srv-applet-mgr/config
	@cp -r openapi.json ../../build/srv-applet-mgr/openapi.json

.PHONY: build
build: openapi lite

.PHONY: run
run: build
	@cd ../../build/srv-applet-mgr && ./srv-applet-mgr

