MODULE_NAME = $(shell cat go.mod | grep "^module" | sed -e "s/module //g")
DOCKER_IMAGE = $(USER)/w3bstream:main
STUDIO_DOCKER_IMAGE = $(USER)/w3bstream-studio:main
DOCKER_COMPOSE_FILE = ./docker-compose.yaml

update_go_module:
	go mod tidy

install_toolkit: update_go_module
	@go install github.com/machinefi/w3bstream/pkg/depends/gen/cmd/...

install_easyjson: update_go_module
	@go install github.com/mailru/easyjson/...@latest

## TODO add source format as a githook
format: install_toolkit
	@toolkit fmt

## gen code
generate: install_toolkit install_easyjson
	@go generate ./...
	@toolkit fmt

## to migrate database models, if model defines changed, make this entry
migrate: install_toolkit install_easyjson
	go run cmd/srv-applet-mgr/main.go migrate

## build srv-applet-mgr
build_server:
	@cd cmd/srv-applet-mgr && go build
	@echo 'succeed! srv-applet-mgr =>cmd/srv-applet-mgr/srv-applet-mgrr*'
	@echo 'succeed! config =>cmd/srv-applet-mgr/config'
	@echo 'modify cmd/srv-applet-mgr/config/local.yaml to use your server config'

build_server_for_docker: update_go_module
	@cd cmd/srv-applet-mgr && GOOS=linux GOWORK=off CGO_ENABLED=1 go build
	@mkdir -p build
	@mv cmd/srv-applet-mgr/srv-applet-mgr build
	@cp -r cmd/srv-applet-mgr/config build/config

#
update_studio:
	@cd studio && git pull origin main

init_submodules:
	@git submodule update --init

# build docker images
build_backend_image: update_go_module
	@docker build -f Dockerfile -t ${DOCKER_IMAGE} .

build_studio_image: init_submodules update_studio
	@cd studio && docker build -f Dockerfile -t ${STUDIO_DOCKER_IMAGE} .

build_docker_images: build_backend_image build_studio_image

# stop server running in docker containers
stop_docker:
	@docker-compose -f ${DOCKER_COMPOSE_FILE} stop

# stop docker and delete docker resouces
drop_docker:
	@docker-compose -f ${DOCKER_COMPOSE_FILE} down

# restart server in docker containers
restart_docker: drop_docker run_docker

# run server in docker containers
run_docker:
	@WS_WORKING_DIR=$(shell pwd)/working_dir WS_BACKEND_IMAGE=${DOCKER_IMAGE} WS_STUDIO_IMAGE=${STUDIO_DOCKER_IMAGE} docker-compose -p w3bstream -f ${DOCKER_COMPOSE_FILE} up -d

## migrate first
run_server: build_server
	@cd build && ./srv-applet-mgr

## make pub_client
build_pub_client: update_go_module
	@cd cmd/pub_client && go build
	@mkdir -p build
	@mv cmd/pub_client/pub_client build
	@echo 'succeed! pub_client => build/pub_client*'

clean:
	@rm -rf build/{config,pub_client,srv-applet-mgr}
	@echo 'remove build/{config,pub_client,srv-applet-mgr}'

run_depends:
	@docker-compose -f testutil/docker-compose-pg.yaml up -d
	@docker-compose -f testutil/docker-compose-mqtt.yaml up -d

stop_depends:
	@docker-compose -f testutil/docker-compose-pg.yaml stop
	@docker-compose -f testutil/docker-compose-mqtt.yaml stop

wasm_demo: update_go_module
	@cd _examples && make all

build: build_server build_pub_client

