include .env

CONTAINER = core
DOCKER_COMPOSE_EXISTS := $(shell command -v docker-compose 2> /dev/null)

.PHONY: start stop ps simulate ssh-core run-core build-core deps-core test-core ssh-rider run-rider build-rider deps-rider test-rider ssh-client run-client build-client deps-client test-client

start:
	@docker-compose up # -d # @todo enable -d

stop:
	@docker-compose stop

ps:
	@docker-compose ps

simulate:
	@echo "[+] Running outside of container. Spawning 1000 riders."
	#@for i in {1..20}; do go run cmd/rider/main.go 1000 ${RIDER_API_KEY} ${CORE_URL}; done;
	@go run cmd/rider/main.go 1000 ${RIDER_API_KEY} ${CORE_URL}
# core
ssh-core:
	@docker-compose exec $(CONTAINER) bash

run-core:
	@CompileDaemon -build='make build-core' -graceful-kill -command='./build/core'

build-core:
	@CGO_ENABLED=0 go build -o build/core -v cmd/pkg/main.go

deps-core:
	${call core, mod vendor}

test-core:
	${call core, test -v ./...}


#---- docker enviroment ----
ifdef DOCKER_COMPOSE_EXISTS
define core
	@docker-compose exec ${CONTAINER} go ${1}
endef
endif