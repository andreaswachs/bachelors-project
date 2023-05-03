.PHONY: setup
setup:
	./bin/setup.sh

.PHONY: clean-docker
clean-docker:
	./bin/clean-docker.sh

.PHONY: dev
dev:
	docker-compose up --build

.PHONY: install-proto
install-proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	@echo ""
	@echo "Ensure that you have ~/go/bin located in your PATH"
	@echo ""

.PHONY: pull-images
pull-images:
	find . -type f -name "*.yaml" \
		| xargs cat \
		| grep "image:" \
		| cut -d ":" -f 2-3 \
		| sed -e "s/\"//g" \
		| sort \
		| uniq \
		| xargs -n1 -P16 docker pull

.PHONY: install-tools
install-tools:
	sudo apt install docker docker-compose -y
	sudo ufw enable
	sudo ufw allow ssh
	sudo ufw allow 50052/tcp
	sudo ufw allow 40000:50000/tcp

.PHONY: deps
deps: install-tools pull-images
