include .env

CONTAINER = core
DOCKER_COMPOSE_EXISTS := $(shell command -v docker-compose 2> /dev/null)

.PHONY: start stop ps simulate ssh ssh-redis run build deps test

start:
	@echo "starting application environment. will start in daemon mode. run 'make ps' for container status"
	@docker-compose up -d

stop:
	@docker-compose stop

ps:
	@docker-compose ps

simulate:
	@echo "running outside of container."
	@source .env; go run cmd/simulation/main.go 5 ${RIDER_API_KEY} ${CLIENT_API_KEY} ${CORE_URL}
# core
ssh:
	@docker-compose exec $(CONTAINER) bash

run:
	@CompileDaemon -build='make build' -graceful-kill -command='./build/app'

build:
	@echo "running build for development."
	@CGO_ENABLED=0 go build -o build/app -v cmd/pkg/main.go

ssh-redis:
	@echo "make sure to authenticate using AUTH default $REDIS_PASSWORD"
	@docker compose exec redis redis-cli

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