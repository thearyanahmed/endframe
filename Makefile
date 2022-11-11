include .env

CONTAINER = core
DOCKER_COMPOSE_EXISTS := $(shell command -v docker-compose 2> /dev/null)

.PHONY: start stop ps simulate ssh run build deps test

start:
	@docker-compose up # -d # @todo enable -d

stop:
	@docker-compose stop

ps:
	@docker-compose ps

simulate:
	@echo "[+] Running outside of container. Spawning 1000 riders."
	#@for i in {1..20}; do go run cmd/simulation.go 1000 ${RIDER_API_KEY} ${CORE_URL}; done;
	@go run cmd/simulation.go 1000 ${RIDER_API_KEY} ${CORE_URL}
# core
ssh:
	@docker-compose exec $(CONTAINER) bash

run:
	@CompileDaemon -build='make build' -graceful-kill -command='./build/app'

build:
	@CGO_ENABLED=0 go build -o build/app -v cmd/pkg/main.go

deps:
	${call app_container, mod vendor}

test:
	${call app_container, test -v ./...}


#---- docker enviroment ----
ifdef DOCKER_COMPOSE_EXISTS
define app_container
	@docker-compose exec ${CONTAINER} go ${1}
endef
endif