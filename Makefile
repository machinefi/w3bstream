DOCKER_COMPOSE_FILE = ./docker-compose.yaml
WS_BACKEND_IMAGE = $(USER)/w3bstream:main
WS_WORKING_DIR=$(shell pwd)/working_dir

.DEFAULT_GOAL := all

## cmd build entries

## update and download go module dependency
.PHONY: update
update:
	@go mod tidy
	@go mod download

## toolkit for code generation
.PHONY: toolkit
toolkit:
	@cd pkg/depends/gen/cmd
	@go install ./...
	@echo installed `which toolkit`

## build cmd/srv-applet-mgr
.PHONY: srv_applet_mgr
srv_applet_mgr:
	@toolkit fmt
	@cd cmd/srv-applet-mgr && make --no-print-directory
	@echo srv-applet-mgr is built to "\033[31m ./build/srv-applet-mgr/... \033[0m"

## build cmd/pub_client
.PHONY: pub_client
pub_client:
	@cd cmd/pub_client && make --no-print-directory
	@echo pub_client is built to "\033[31m ./build/pub_client/... \033[0m"

.PHONY: all
all: build test

.PHONY: build
build: update toolkit srv_applet_mgr pub_client

.PHONY: clean
clean:
	@rm -rf ./build/config ./build/pub_client ./build/srv-applet-mgr

## docker build entries

.PHONY: build_image
build_image:
	@docker build -f cmd/srv-applet-mgr/Dockerfile -t ${WS_BACKEND_IMAGE} .

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
	@go test -cover -coverprofile=coverage.out ./...
	@docker stop mqtt_test postgres_test redis_test || true && docker container rm mqtt_test postgres_test redis_test || true

bench: test_depends
	@cd ./cmd/srv-applet-mgr/tests/integrations/ && go test event_benchmark_test.go -bench=.
	@docker stop mqtt_test postgres_test redis_test || true && docker container rm mqtt_test postgres_test redis_test || true

.PHONY: test_depends
test_depends: cleanup_test_depends postgres_test mqtt_test redis_test

.PHONY: cleanup_test_depends
cleanup_test_depends:
	@docker stop mqtt_test postgres_test redis_test || true && docker container rm mqtt_test postgres_test redis_test || true

.PHONY: postgres_test
postgres_test:
	docker run --name postgres_test -e POSTGRES_PASSWORD=test_passwd -e POSTGRES_USER=root -p 15432:5432 -d postgres:14-alpine

.PHONY: mqtt_test
mqtt_test:
	docker run --name mqtt_test -p 11883:1883 -d eclipse-mosquitto:1.6.15

.PHONY: redis_test
redis_test:
	docker run --name redis_test -p 16379:6379 -d redis:6.2

