MODULE_NAME = $(shell cat go.mod | grep "^module" | sed -e "s/module //g")
TOOLKIT_PKG = ${MODULE_NAME}/gen/cmd/toolkit

install_toolkit:
	@go install "${TOOLKIT_PKG}/..."

## TODO add source format as a githook
format: install_toolkit
	go mod tidy
	toolkit fmt

generate: install_toolkit
	go mod tidy
	go generate ./...
	toolkit fmt


export PG_TEST_DB_NAME=test
export PG_TEST_DB_USER=test_user
export PG_TEST_DB_PASSWD=test_passwd
export PG_TEST_HOSTNAME='postgres://$(PG_TEST_DB_USER):$(PG_TEST_DB_PASSWD)@127.0.0.1:5432'
export PG_TEST_MASTER_EP='$(PG_TEST_HOSTNAME)/$(PG_TEST_DB_NAME)'
export PG_TEST_SLAVE_EP=$(PG_TEST_HOSTNAME)


pg_envs:
	@echo "=== print env variable ==="
	@echo 'PG_TEST_DB_NAME   = $(PG_TEST_DB_NAME)'
	@echo 'PG_TEST_DB_USER   = $(PG_TEST_DB_USER)'
	@echo 'PG_TEST_DB_PASSWD = $(PG_TEST_DB_PASSWD)'
	@echo 'PG_TEST_HOSTNAME  = $(PG_TEST_HOSTNAME)'
	@echo 'PG_TEST_MASTER_EP = $(PG_TEST_MASTER_EP)'
	@echo 'PG_TEST_SLAVE_EP  = $(PG_TEST_SLAVE_EP)'
	@echo "=== print env variable end  ===\n"

pg_start:
	@if [[ $$(pg_isready -h localhost) != "localhost:5432 - accepting connections" ]] ; \
	then \
		echo "=== start postgres server ==="; \
		docker-compose -f testutil/docker-compose-pg.yaml up -d ; \
		echo "=== init database ===" ; \
		for i in {1..5} ; \
		do \
			if [[ $$(pg_isready -h localhost) =~ "accepting connections" ]] ; \
			then \
				psql $(PG_TEST_HOSTNAME) -c 'create database $(PG_TEST_DB_NAME)' && \
				psql $(PG_TEST_HOSTNAME) -c 'create schema $(PG_TEST_DB_NAME)' ; \
				break ; \
			else \
				echo "server not ready, retry in 10 second" ; sleep 10 ; \
			fi \
		done ; \
		if [[ $$(pg_isready -h localhost) != "localhost:5432 - accepting connections" ]] ; \
		then \
			echo "=== database init failed ==="  ; \
			exit 1;  \
		fi \
	fi ; \

test:
	go test ./...
	@echo "=================TEST FINISHED================="
