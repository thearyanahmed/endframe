include .env

CORE_CONTAINER = core
RIDER_CONTAINER = rider
CLIENT_CONTAINER = client
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
	@docker-compose exec $(CORE_CONTAINER) bash

run-core:
	@CompileDaemon -build='make build-core' -graceful-kill -command='./build/core'

build-core:
	@CGO_ENABLED=0 go build -o build/core -v cmd/core/main.go

deps-core:
	${call core, mod vendor}

test-core:
	${call core, test -v ./...}

# rider
ssh-rider:
	@docker-compose exec $(CORE_CONTAINER) bash

run-rider:
	@CompileDaemon -build='make build-rider' -graceful-kill -command='./build/rider'

build-rider:
	@CGO_ENABLED=0 go build -o build/rider  -v cmd/rider/main.g

deps-rider:
	${call rider, mod vendor}

test-rider:
	${call rider, test -v ./...}

# client
ssh-client:
	@docker-compose exec $(CORE_CONTAINER) bash

run-client:
	@CompileDaemon -build='make build-client' -graceful-kill -command='./build/client'

build-client:
	@CGO_ENABLED=0 go build -o build/client -v cmd/client/main.g

deps-client:
	${call client, mod vendor}

test-client:
	${call client, test -v ./...}


#---- docker enviroment ----
ifdef DOCKER_COMPOSE_EXISTS
define core
	@docker-compose exec ${CORE_CONTAINER} go ${1}
endef
define rider
	@docker-compose exec ${RIDER_CONTAINER} go ${1}
endef
define client
	@docker-compose exec ${CLIENT_CONTAINER} go ${1}
endef
endif