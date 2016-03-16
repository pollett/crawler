project_name=crawler

export GOPATH := $(PWD)/vendor:$(PWD)

up: build
	docker-compose -p ${project_name} up

logs: 
	docker-compose -p ${project_name} logs

restart: 
	docker-compose -p ${project_name} restart

stop: 
	docker-compose -p ${project_name} stop

clean: stop
	docker-compose -p ${project_name} rm -f

build: 
	docker-compose -p ${project_name} build

.PHONY: up logs restart stop clean build
