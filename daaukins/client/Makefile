DKN_INSTALL_DIR ?= ~/opt/bin

.PHONY: build install

build:
	@echo "Building dkn..."
	@go build -v -o dkn

install: build
	@echo "Installing dkn..."
	@mkdir -p ${DKN_INSTALL_DIR}
	@cp dkn ${DKN_INSTALL_DIR}
	@echo "Ensure that ${DKN_INSTALL_DIR} is in your PATH."
