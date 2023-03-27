setup:
	./bin/setup.sh

clean-docker:
	./bin/clean-docker.sh

grpc-build:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative service/service.proto

grpc-debug:
	evans --path service --proto service.proto --host localhost

dev:
	go build
	./server

dev-docker:
	docker-compose up --build

install-proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	@echo ""
	@echo "Ensure that you have ~/go/bin located in your PATH"
	@echo ""

system-integration-test-leader-leader:
	@echo ""
	@echo "Success criteria: an error log message is posted complaining that a follower is a leader, but should be a follower"
	@echo "You will have to manually CTRL-C to stop the test when you see the message, until I figure out how to do this automatically"
	@echo ""
	docker-compose -f system_integration_tests/leader-leader.yaml up
