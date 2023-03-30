DKN_ADDR ?= localhost
DKN_PORT ?= 50052

.PHONY: build
build:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative service.proto

.PHONY: debug
debug:
	evans --path . --proto service.proto --host ${DKN_ADDR} --port ${DKN_PORT}
