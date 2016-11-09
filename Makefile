.PHONY: start-containers stop-containers test
.DEFAULT_GOAL := test

export ETCD_ENDPOINTS=http://localhost:2379
DOCKER_MACHINE_IP=$(shell docker-machine ip default 2> /dev/null)
ifneq ($(DOCKER_MACHINE_IP),)
	ETCD_ENDPOINTS=http://$(DOCKER_MACHINE_IP):2379
endif

start-containers:
	docker-compose up -d

stop-containers:
	docker-compose down

test: start-containers
	@echo Running Tests
	@go test -v ./...
