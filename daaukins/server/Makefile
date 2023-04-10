setup:
	./bin/setup.sh

clean-docker:
	./bin/clean-docker.sh


grpc-debug: evans --path service --proto service.proto --host localhost

grpc-debug-docker:
	evans --path service --proto service.proto --host localhost --port 50052
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

admin-deps:
	@go install github.com/cloudflare/cfssl/cmd/...@latest

pull-images:
	find . -type f -name "*.yaml" \
		| xargs cat \
		| grep "image:" \
		| cut -d ":" -f 2-3 \
		| sed -e "s/\"//g" \
		| sort \
		| uniq \
		| xargs -n1 -P16 docker pull

gen-rootca:
	cfssl selfsign -config cfssl.json --profile rootca "My Root CA" csr.json | cfssljson -bare root

gen-certs:
	./bin/gen-certs.sh

clean-certs:
	rm -rf certs || true
	mkdir certs

full-gen-certs: clean-certs gen-rootca gen-certs

