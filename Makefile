# Makefile for the CrawlX project

# Variables
BINARY_NAME=crawlx
CMD_PATH=./cmd/crawlx
DIST_PATH=./dist

# Detect OS (Windows_NT is set by Windows)
ifeq ($(OS),Windows_NT)
    MKDIR=if not exist $(DIST_PATH) mkdir $(DIST_PATH)
    RM=rmdir /s /q $(DIST_PATH)
    BINEXT=.exe
else
    MKDIR=mkdir -p $(DIST_PATH)
    RM=rm -rf $(DIST_PATH)
    BINEXT=
endif

# Default target
.PHONY: help
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  run          - Run the crawler with default test arguments"
	@echo "  build        - Build the binary for the current operating system"
	@echo "  build-all    - Cross-compile binaries for Windows, Linux, and macOS"
	@echo "  clean        - Remove the build artifacts from the dist/ folder"

# Run the application
.PHONY: run
run:
	go run $(CMD_PATH) -u https://toscrape.com -d 1

# Build for the current OS
.PHONY: build
build:
	@echo Building for current OS...
	@$(MKDIR)
	go build -o $(DIST_PATH)/$(BINARY_NAME)$(BINEXT) $(CMD_PATH)
	@echo Build complete!

# Build for all target platforms
.PHONY: build-all
build-all:
	@echo Building for all platforms...
	@$(MKDIR)
	GOOS=linux GOARCH=amd64 go build -o $(DIST_PATH)/$(BINARY_NAME)-linux-amd64 $(CMD_PATH)
	GOOS=windows GOARCH=amd64 go build -o $(DIST_PATH)/$(BINARY_NAME)-windows-amd64.exe $(CMD_PATH)
	GOOS=darwin GOARCH=amd64 go build -o $(DIST_PATH)/$(BINARY_NAME)-darwin-amd64 $(CMD_PATH)
	@echo All builds complete!

# Clean the dist directory
.PHONY: clean
clean:
	@echo Cleaning build artifacts...
	@$(RM)
	@echo Clean complete.
