.PHONY: build test dev clean

# 变量定义
BINARY_NAME=singdns
BINARY_DIR=bin
CORE_DIR=core

build:
	@echo "Building..."
	@mkdir -p $(BINARY_DIR)
	@go build -o $(BINARY_DIR)/$(BINARY_NAME) ./cmd/singdns

test:
	@echo "Running tests..."
	@go test -v ./...

dev:
	@echo "Starting development server..."
	@go run ./cmd/singdns

clean:
	@echo "Cleaning..."
	@rm -rf $(BINARY_DIR)
	@go clean

install:
	@echo "Installing..."
	@mkdir -p /usr/local/$(BINARY_NAME)/core
	@cp $(BINARY_DIR)/$(BINARY_NAME) /usr/local/bin/
	@cp -r configs /etc/$(BINARY_NAME)
	@cp -r core/* /usr/local/$(BINARY_NAME)/core/

uninstall:
	@echo "Uninstalling..."
	@rm -f /usr/local/bin/$(BINARY_NAME)
	@rm -rf /usr/local/$(BINARY_NAME)
	@rm -rf /etc/$(BINARY_NAME) 