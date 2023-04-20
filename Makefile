DOCKER_COMPOSE_FILE = ./docker-compose.yaml
WS_BACKEND_IMAGE = $(USER)/w3bstream:main
WS_WORKING_DIR=$(shell pwd)/working_dir

.PHONY: update_go_module
update_go_module:
	@go mod tidy

.PHONY: build
build: build_server build_pub_client

.PHONY: build_server
build_server:
	@mkdir -p build
	@cd cmd/srv-applet-mgr && go build
	@cd cmd/srv-applet-mgr && go generate .
	@rm -rf build/{config,srv-applet-mgr,openapi.json}
	@mv cmd/srv-applet-mgr/srv-applet-mgr build/
	@cp -r cmd/srv-applet-mgr/config build/config
	@cp -r cmd/srv-applet-mgr/openapi.json build/openapi.json
	@echo 'succeed! srv-applet-mgr =>cmd/srv-applet-mgr/srv-applet-mgr'
	@echo 'succeed! config =>cmd/srv-applet-mgr/config'
	@echo 'modify cmd/srv-applet-mgr/config/local.yaml to use your server config'

.PHONY: build_pub_client
build_pub_client: update_go_module
	@cd cmd/pub_client && go build
	@mkdir -p build/pub_client build/pub_client/config
	@mv -f cmd/pub_client/pub_client build/pub_client
	@cp -u cmd/pub_client/config/config.yml build/pub_client/config/config.yml
	@echo 'succeed! pub_client => build/pub_client/pub_client'

.PHONY: build_docker_images
build_docker_images: build_backend_image

.PHONY: build_backend_image
build_backend_image: update_go_module
	@docker build -f Dockerfile -t ${WS_BACKEND_IMAGE} .

# run server in docker containers
.PHONY: run_docker
run_docker:
	@WS_WORKING_DIR=${WS_WORKING_DIR} WS_BACKEND_IMAGE=${WS_BACKEND_IMAGE} WS_STUDIO_IMAGE=${WS_STUDIO_IMAGE} docker-compose -p w3bstream -f ${DOCKER_COMPOSE_FILE} up -d

## migrate first
.PHONY: run_server
run_server: build_server
	@cd build && ./srv-applet-mgr

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

.PHONY: clean
clean:
	@rm -rf ./build/config ./build/pub_client ./build/srv-applet-mgr

.PHONY: install_toolkit
install_toolkit:
	@if [ ! -f "$$GOBIN/toolkit" ] ; \
	then \
		go install github.com/machinefi/w3bstream/pkg/depends/gen/cmd/...@toolkit ; \
		echo "toolkit installed" ; \
	fi
	@echo `which toolkit`

update_toolkit:
	@go install github.com/machinefi/w3bstream/pkg/depends/gen/cmd/...@toolkit

.PHONY: generate
generate: install_toolkit
	@go generate ./...

## to migrate database models, if model defines changed, make this entry
.PHONY: migrate
migrate: install_toolkit 
	go run cmd/srv-applet-mgr/main.go migrate
