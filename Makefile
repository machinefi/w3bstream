DOCKER_COMPOSE_FILE = ./docker-compose.yaml
WS_BACKEND_IMAGE = $(USER)/w3bstream:main
WS_WORKING_DIR=$(shell pwd)/working_dir

.DEFAULT_GOAL := build

## cmd build entries

.PHONY: update
update:
	@go mod tidy
	@go mod download

.PHONY: toolkit
toolkit:
	@go install github.com/machinefi/w3bstream/pkg/depends/gen/cmd/...@toolkit-patch-0.0.3
	@echo installed `which toolkit`

.PHONY: srv_applet_mgr
srv_applet_mgr:
	@toolkit fmt
	@cd cmd/srv-applet-mgr && make --no-print-directory
	@echo srv-applet-mgr is built to "\033[31m ./build/srv-applet-mgr/... \033[0m"

.PHONY: srv_applet_mgr_lite
srv_applet_mgr_lite:
	@cd cmd/srv-applet-mgr && make lite --no-print-directory
	@echo srv-applet-mgr is built to "\033[31m ./build/srv-applet-mgr/... \033[0m"


.PHONY: pub_client
pub_client:
	@cd cmd/pub_client && make --no-print-directory
	@echo pub_client is built to "\033[31m ./build/pub_client/... \033[0m"

.PHONY: build
build: update test toolkit srv_applet_mgr pub_client

.PHONY: build_lite
build_lite: update srv_applet_mgr_lite

.PHONY: clean
clean:
	@rm -rf ./build/config ./build/pub_client ./build/srv-applet-mgr

## docker build entries

.PHONY: build_docker_images
build_docker_images: build_backend_image

.PHONY: build_backend_image
build_backend_image: update
	@docker build -f Dockerfile -t ${WS_BACKEND_IMAGE} .

# run server in docker containers
.PHONY: run_docker
run_docker:
	@WS_WORKING_DIR=${WS_WORKING_DIR} WS_BACKEND_IMAGE=${WS_BACKEND_IMAGE} WS_STUDIO_IMAGE=${WS_STUDIO_IMAGE} docker-compose -p w3bstream -f ${DOCKER_COMPOSE_FILE} up -d

# stop server running in docker containers
.PHONY: stop_docker
stop_docker:
	@docker-compose -f ${DOCKER_COMPOSE_FILE} stop

# stop docker and delete docker resources
.PHONY: drop_docker
drop_docker:
	@docker-compose -f ${DOCKER_COMPOSE_FILE} down

# restart server in docker containers
.PHONY: restart_docker
restart_docker: drop_docker run_docker

## developing stage entries

.PHONY: generate
generate: toolkit
	@cd pkg/models              && go generate ./...
	@cd pkg/enums               && go generate ./...
	@cd pkg/errors              && go generate ./...
	@cd pkg/errors              && go generate ./...
	@cd pkg/depends/util/strfmt && go generate ./...

## to migrate database models, if model defines changed, make this entry
.PHONY: migrate
migrate: toolkit
	go run cmd/srv-applet-mgr/main.go migrate

.PHONY: test
test: test_depends
	@go test -cover ./...
	@docker stop mqtt_test postgres_test || true && docker container rm mqtt_test postgres_test || true

.PHONY: test_depends
test_depends: postgres_test mqtt_test

.PHONY: cleanup_test_depends
cleanup_test_depends:
	@docker stop mqtt_test postgres_test || true && docker container rm mqtt_test postgres_test || true

.PHONY: postgres_test
postgres_test:
	docker run --name postgres_test -e POSTGRES_PASSWORD=test_passwd -e POSTGRES_USER=root -p 5432:5432 -d postgres:14-alpine

.PHONY: mqtt_test
mqtt_test:
	docker run --name mqtt_test -p 1883:1883 -d eclipse-mosquitto:1.6.15

