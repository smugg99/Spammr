PROJECT_NAME := $(notdir $(CURDIR))
BINARY_NAME := $(PROJECT_NAME)

PROJECT_ROOT := $(CURDIR)
BACKEND_SOURCE_DIR := $(PROJECT_ROOT)/app

BUILD_DIR := $(PROJECT_ROOT)/build
CURRENT_BUILD_DIR := $(BUILD_DIR)/$(BINARY_NAME)

.PHONY: all
all: build run

.PHONY: build
build: | $(CURRENT_BUILD_DIR)
	@echo "Building backend..."
	@go build -o $(CURRENT_BUILD_DIR)/$(BINARY_NAME) $(BACKEND_SOURCE_DIR)
	# @cp -r $(BACKEND_SOURCE_DIR)/config.json $(CURRENT_BUILD_DIR)/config.json
	# @cp -r $(BACKEND_SOURCE_DIR)/.env $(CURRENT_BUILD_DIR)/.env

.PHONY: 
run: | $(CURRENT_BUILD_DIR)
	@echo "Running the application..."
	@cd $(CURRENT_BUILD_DIR) && ./$(BINARY_NAME)

$(CURRENT_BUILD_DIR):
	@mkdir -p $(CURRENT_BUILD_DIR)

.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

.PHONY: help
help:
	@echo "Usage:"
	@echo "  make                      Build and run the appplication"
	@echo "  make build                Build the application"
	@echo "  make run                  Run the application"
	@echo "  make clean                Clean up build files"
	@echo "  make help                 Show this help message"

.PHONY: build build-web build-web-install run clean help
