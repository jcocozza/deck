APP_NAME=deck
BUILD_DIR=bin

compress-tar:
	tar -czvf $(BUILD_DIR)/$(APP_NAME)_$(OS)_$(ARCH).tar.gz $(BUILD_DIR)/$(APP_NAME)_$(OS)_$(ARCH)

compress-zip:
	cd $(BUILD_DIR) && zip $(APP_NAME)_$(OS)_$(ARCH).zip $(APP_NAME)_$(OS)_$(ARCH)


.PHONY: build
build:
	@echo "Building $(APP_NAME) for $(OS)/$(ARCH)"
	GOOS=$(OS) GOARCH=$(ARCH) go build -o "$(BUILD_DIR)/$(APP_NAME)_$(OS)_$(ARCH)"


linux:
	$(MAKE) OS=linux ARCH=amd64 build compress-tar

darwin-amd:
	$(MAKE) OS=darwin ARCH=amd64 build compress-tar

darwin-arm:
	$(MAKE) OS=darwin ARCH=arm64 build compress-tar

windows:
	$(MAKE) OS=windows ARCH=amd64 build compress-zip

clean:
	rm -rf $(BUILD_DIR)

all: clean linux darwin-amd darwin-arm windows
