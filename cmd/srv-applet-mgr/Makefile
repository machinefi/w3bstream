.DEFAULT_GOAL := build

REMOTE=$(shell git remote -v | head -n 1 | awk -F" " '{ print $$2 }')
BRANCH=$(shell git branch | grep '*' | awk '{ print $$2; }')
COMMIT_ID=$(shell git rev-parse --short HEAD)
COMMIT_AUTHOR=$(shell git log | head -n 2 | grep Author | awk -F "[<>]" '{ print $$2;  }')

.PHONY: build
build:
	@go build -ldflags "-s -w -X ./global.ServiceInfo=${REMOTE}..${BRANCH}..${COMMIT_ID}..${COMMIT_AUTHOR}" -o srv-applet-mgr
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
