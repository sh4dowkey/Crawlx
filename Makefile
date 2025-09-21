# Variables
BINARY_NAME=crawlx
CMD_PATH=./cmd/crawlx
DIST_PATH=./dist

# Detect Windows
ifeq ($(OS),Windows_NT)
    MKDIR_CMD=powershell -Command "if (-not (Test-Path '$(DIST_PATH)')) { New-Item -ItemType Directory -Path '$(DIST_PATH)' }"
    RM_CMD=powershell -Command "if (Test-Path '$(DIST_PATH)') { Remove-Item -Recurse -Force '$(DIST_PATH)' }"
    BINEXT=.exe
else
    MKDIR_CMD=mkdir -p $(DIST_PATH)
    RM_CMD=rm -rf $(DIST_PATH)
    BINEXT=
endif

.PHONY: help
help:
	@echo "Usage: make <target>"
	@echo "Targets: run, build, build-all, clean"

.PHONY: run
run:
	go run $(CMD_PATH) -u https://toscrape.com -d 1

.PHONY: build
build:
	@echo Building for current OS...
	@$(MKDIR_CMD)
	go build -o $(DIST_PATH)/$(BINARY_NAME)$(BINEXT) $(CMD_PATH)
	@echo Build complete!

.PHONY: build-all
build-all:
	@echo Building for all platforms...
	@$(MKDIR_CMD)
	GOOS=linux GOARCH=amd64 go build -o $(DIST_PATH)/$(BINARY_NAME)-linux-amd64 $(CMD_PATH)
	GOOS=windows GOARCH=amd64 go build -o $(DIST_PATH)/$(BINARY_NAME)-windows-amd64.exe $(CMD_PATH)
	GOOS=darwin GOARCH=amd64 go build -o $(DIST_PATH)/$(BINARY_NAME)-darwin-amd64 $(CMD_PATH)
	@echo All builds complete!

.PHONY: clean
clean:
	@echo Cleaning build artifacts...
	@$(RM_CMD)
	@echo Clean complete.
