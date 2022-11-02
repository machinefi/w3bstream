MODULE_NAME = $(shell cat go.mod | grep "^module" | sed -e "s/module //g")

update_go_module:
	go mod tidy

install_toolkit: update_go_module
	@go install github.com/machinefi/Bumblebee/gen/cmd/...@latest

install_goimports: update_go_module
	@go install golang.org/x/tools/cmd/goimports@latest

install_easyjson: update_go_module
	@go install github.com/mailru/easyjson/...@latest

## TODO add source format as a githook
format: install_goimports
	go mod tidy
	goimports -w -l -local "${MODULE_NAME}" ./

## gen code
generate: install_toolkit install_easyjson install_goimports
	go generate ./...
	goimports -w -l -local "${MODULE_NAME}" ./
	toolkit patch goid

## to migrate database models, if model defines changed, make this entry
migrate: install_toolkit install_easyjson install_goimports
	go run cmd/srv-applet-mgr/main.go migrate

## build srv-applet-mgr
build_server: update_go_module generate format
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

build_server_for_docker: update_go_module vendor
	@cd cmd/srv-applet-mgr && GOOS=linux GOWORK=off CGO_ENABLED=1 go build -mod vendor
	@mkdir -p build
	@mv cmd/srv-applet-mgr/srv-applet-mgr build
	@cp -r cmd/srv-applet-mgr/config build/config
	@rm -rf vendor

vendor: update_go_module
	@go mod vendor

#
update_frontend:
	@cd frontend &&	git pull origin main

init_frontend:
	@git submodule update --init

# build docker image
build_image: update_go_module vendor init_frontend update_frontend
	@mkdir -p build_image/pgdata
	@mkdir -p build_image/asserts
	@docker build -t iotex/w3bstream:v3 .
	@rm -rf vendor

# drop docker container
drop_image:
	@docker-compose -f ./docker-compose.yaml down
	#@docker stop iotex_w3bstream
	#@docker rm iotex_w3bstream

# restart docker container
restart_image:
	@docker-compose -f ./docker-compose.yaml down
	@echo "The container was shut down before, now restart it"
	@WS_WORKING_DIR=$(shell pwd)/build_image docker-compose -p w3bstream -f ./docker-compose.yaml up -d

# run docker image
run_image:
	@WS_WORKING_DIR=$(shell pwd)/build_image docker-compose -p w3bstream -f ./docker-compose.yaml up -d
	#@docker run -d -it --name iotex_w3bstream  -e DATABASE_URL="postgresql://test_user:test_passwd@127.0.0.1/test?schema=applet_management" -e NEXT_PUBLIC_API_URL="http://127.0.0.1:8888" -p 5432:5432 -p 8888:8888 -p 1883:1883  -p 3000:3000 -v $(shell pwd)/build_image/pgdata:/var/lib/postgresql_data -v $(shell pwd)/build_image/asserts:/w3bstream/cmd/srv-applet-mgr/asserts -v $(shell pwd)/build_image/conf/srv-applet-mgr/config/local.yml:/w3bstream/cmd/srv-applet-mgr/config/local.yml iotex/w3bstream:v3 /bin/sh /init.sh

## run docker image fast
#run_fast:
#	@sh ./build_image/run_fast/run_w3bstream.sh start


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
	docker-compose -f testutil/docker-compose-pg.yaml up -d
	docker-compose -f testutil/docker-compose-mqtt.yaml up -d

wasm_demo: update_go_module
	@cd _examples && make all

build: build_server build_pub_client

