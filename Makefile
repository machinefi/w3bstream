MODULE_NAME = $(shell cat go.mod | grep "^module" | sed -e "s/module //g")

update_go_module:
	go mod tidy

install_toolkit: update_go_module
	@go install github.com/iotexproject/Bumblebee/gen/cmd/...

install_goimports: update_go_module
	@go install golang.org/x/tools/cmd/goimports@latest

## TODO add source format as a githook
format: install_goimports
	go mod tidy
	goimports -w -l -local "${MODULE_NAME}" ./

## gen code
generate: install_toolkit install_goimports
	go generate ./...
	goimports -w -l -local "${MODULE_NAME}" ./

## to migrate database models, if model defines changed, make this entry
migrate: install_toolkit install_goimports
	go run cmd/srv-applet-mgr/main.go migrate

## build srv-applet-mgr
build_server: update_go_module
	@cd cmd/srv-applet-mgr && go build
	@mkdir -p build
	@mv cmd/srv-applet-mgr/srv-applet-mgr build
	@cp -r cmd/srv-applet-mgr/config build/config
	@echo 'succeed! srv-applet-mgr =>build/srv-applet-mgr*'
	@echo 'succeed! config =>build/config/'
	@echo 'modify config/local.yaml to use your server config'

vendor: update_go_module
	@go mod vendor

# build docker image
build_image: update_go_module vendor
	@mkdir -p build_image/pgdata
	@git submodule update --init
	@docker build -t iotex/w3bstream:v3 .

#drop docker image
drop_image:
	@docker stop iotex_w3bstream
	@docker rm iotex_w3bstream

# run docker image
run_image:
	@docker run -d -it --name iotex_w3bstream -p 5432:5432 -p 8888:8888 -p 3000:3000 -v $(shell pwd)/build_image/pgdata:/var/lib/postgresql_data iotex/w3bstream:v3 /bin/bash /w3bstream/build_image/cmd/docker_init.sh

build_server_vendor: vendor
	@cd cmd/srv-applet-mgr && go build -mod vendor
	@mkdir -p build
	@mv cmd/srv-applet-mgr/srv-applet-mgr build
	@cp -r cmd/srv-applet-mgr/config build/config
	@echo 'succeed! srv-applet-mgr =>build/srv-applet-mgr*'
	@echo 'succeed! config =>build/config/'
	@echo 'modify config/local.yaml to use your server config'

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
	@cd examples && make all

build: build_server build_pub_client

