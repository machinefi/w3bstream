MODULE_NAME = $(shell cat go.mod | grep "^module" | sed -e "s/module //g")
DOCKER_IMAGE = $(USER)/w3bstream:main

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
	@mkdir -p build
	@mv cmd/srv-applet-mgr/srv-applet-mgr build
	@rm -rf build/config
	@mkdir -p build/config
	@cp cmd/srv-applet-mgr/config/default.yml build/config/default.yml
	@cp build_image/etc/srv-applet-mgr/config/local.yml build/config/local.yml
	@echo 'succeed! srv-applet-mgr =>build/srv-applet-mgr*'
	@echo 'succeed! config =>build/config/'
	@echo 'modify config/local.yaml to use your server config'

build_server_for_docker: update_go_module
	@cd cmd/srv-applet-mgr && GOOS=linux GOWORK=off CGO_ENABLED=1 go build
	@mkdir -p build
	@mv cmd/srv-applet-mgr/srv-applet-mgr build
	@cp -r cmd/srv-applet-mgr/config build/config

#
update_frontend:
	@cd frontend &&	git pull origin main

init_frontend:
	@git submodule update --init

# build docker image
build_image: update_go_module init_frontend update_frontend
	@mkdir -p build_image/pgdata
	@mkdir -p build_image/asserts
	@docker build -t ${DOCKER_IMAGE} .

# drop docker container
drop_image:
	@docker-compose -f ./docker-compose.yaml down

# restart docker container
restart_image:
	@docker-compose -f ./docker-compose.yaml down
	@echo "The container was shut down before, now restart it"
	@WS_WORKING_DIR=$(shell pwd)/build_image docker-compose -p w3bstream -f ./docker-compose.yaml up -d

# run docker image
run_image:
	@WS_WORKING_DIR=$(shell pwd)/build_image DOCKER_IMAGE=${DOCKER_IMAGE} docker-compose -p w3bstream -f ./docker-compose.yaml up -d

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

