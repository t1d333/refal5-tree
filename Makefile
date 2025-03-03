all: lint build


BUILD_DIR = ./build
SRC_PATH = ./cmd/app/main.go
COMPILER_BIN = refal5-tree


.PHONY: lint
lint:
	golangci-lint run -v


.PHONY: build
$(BUILD_DIR):
	mkdir -p $@
	go build -o $(BUILD_DIR)/$(COMPILER_BIN) $(SRC_PATH)
	

.PHONY: clean
clean:
	rm -rf build
