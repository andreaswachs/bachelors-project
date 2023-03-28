.PHONY: build install

build:
	@echo "Building dkn..."
	@go build -v -o dkn

install: build
	@echo "Installing dkn..."
	@mkdir -p ~/opt/bin
	@cp dkn ~/opt/bin
	@echo "Ensure that ~/opt/bin is in your PATH."
