MAKEFLAGS += --warn-undefined-variables
.DEFAULT_GOAL := test

.PHONY: start-containers stop-containers test build release

export ETCD_ENDPOINTS=http://localhost:2379
DOCKER_MACHINE_IP=$(shell docker-machine ip default 2> /dev/null)
ifneq ($(DOCKER_MACHINE_IP),)
	ETCD_ENDPOINTS=http://$(DOCKER_MACHINE_IP):2379
endif

VERSION ?= dev-build

start-containers:
	docker-compose up -d

stop-containers:
	docker-compose down

test: start-containers
	@echo Running Tests
	@go test -v ./...

build:
	@mkdir -p build
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/etcd-template -v github.com/thrawn01/etcd-template/cmd/etcd-template

release: build
	mkdir -p release
	git tag $(VERSION)
	git push --tags
	cd build && tar -czf ../release/etcd-template-$(VERSION).tar.gz etcd-template
	@echo
	@cd release && shasum etcd-template-$(VERSION).tar.gz
	@cd release && shasum etcd-template-$(VERSION).tar.gz > etcd-template-$(VERSION).sha1.txt
	@echo Upload files in release/ directory to GitHub release.
