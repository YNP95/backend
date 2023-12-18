default: help

help:
	@echo 'usage: make [target] ...'
	@echo ''
	@echo 'targets:'
	@egrep '^(.+)\:\ .*##\ (.+)' ${MAKEFILE_LIST} | sed 's/:.*##/#/' | column -t -c 2 -s '#'

all:
	make clean
	make docs
	make build-server
	make build-batch

clean:
	go clean

docs:	dummy
	swag init -g cmd/server/server.go

build-server:	dummy
	go build api/*.go
	go build cmd/server/*.go

build-batch:	dummy
	go build api/*.go
	go build cmd/batch/*.go

dummy:
