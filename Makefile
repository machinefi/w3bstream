DOCKER_COMPOSE_FILE=./docker-compose.yaml
WS_BACKEND_IMAGE = $(USER)/w3bstream:main
WS_WORKING_DIR=$(shell pwd)/working_dir

.DEFAULT_GOAL := all

## cmd build entries

## update and download go module dependency
.PHONY: update
update:
	@go mod tidy
	@go mod download

## build all targets to ./build/
.PHONY: targets
targets: update
	@cd cmd && for target in * ; \
	do \
		echo "\033[32mbuilding $$target ... \033[0m" ; \
		if [ -d $$target ] && [ -e $$target/Makefile ]; then \
			cd $$target; \
			make target --no-print-directory; \
			cd ..; \
		else \
			echo "\033[31mno entry\033[0m" ; \
		fi; \
		echo "\033[32mdone!\033[0m\n"; \
	done

## build all docker images
.PHONY: images
images:
	@cd cmd && for target in * ; \
	do \
		echo "\033[32mbuilding $$target docker image ... \033[0m" ; \
		if [ -d $$target ] && [ -e $$target/Dockerfile ]; then \
			cd $$target; \
			make image --no-print-directory || true; \
			cd ..; \
		else \
			echo "\033[31mno entry\033[0m" ; \
		fi; \
		echo "\033[32mdone!\033[0m\n"; \
	done

.PHONY: all
all: update targets test images

.PHONY: clean
clean:
	@cd cmd && for target in * ; \
	do \
		echo "\033[32mcleaning $$target ... \033[0m" ; \
		if [ -d $$target ] && [ -e $$target/Makefile ]; then \
			cd $$target; \
			make clean --no-print-directory || true; \
			cd ..; \
		else \
			echo "\033[31mno entry\033[0m" ; \
		fi; \
		echo "\033[32mdone!\033[0m\n" ; \
	done


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

## toolkit for code generation
.PHONY: toolkit
toolkit:
	@go install github.com/machinefi/w3bstream/pkg/depends/gen/cmd/...@toolkit-patch-0.0.3
	@echo installed `which toolkit`

## gomock for generating mock code
.PHONY: gomock
gomock:
	@go install github.com/golang/mock/mockgen@v1.6.0

.PHONY: generate
generate: toolkit gomock
	@cd pkg/models              && go generate ./...
	@cd pkg/enums               && go generate ./...
	@cd pkg/errors              && go generate ./...
	@cd pkg/depends/util/strfmt && go generate ./...
	@cd pkg/test                && go generate ./...

.PHONY: precommit
precommit: toolkit targets test
	@toolkit fmt
	@cd cmd/srv-applet-mgr && make openapi --no-print-directory
	@git add -u

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
