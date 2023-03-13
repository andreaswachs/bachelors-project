setup:
	./bin/setup.sh

clean-docker:
	./bin/clean-docker.sh

build:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative service/service.proto

install-proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	@echo ""
	@echo "Ensure that you have ~/go/bin located in your PATH"
	@echo ""
