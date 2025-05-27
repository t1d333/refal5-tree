all: lint build


BUILD_DIR = ./build
SRC_PATH = ./cmd/app/main.go
COMPILER_BIN = refal5t


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

.PHONY: autotests
autotests: build
	cd autotests && bash run.sh
	
	
rmcc: rmcc.ref
	./build/refal5t rmcc.ref
	go build rmcc.go

generate-test-inputs:
	python3 ./scripts/gen-files.py -o inputs 1kb 2kb 5kb 10kb 32kb 64kb 128kb 512kb 1mb
	
fab: fab.ref
	refal5t fab.ref
	go build fab.go
	
.PHONY: test-fab
test-fab: fab
	@BIN=./fab && \
	DIR=./inputs && \
	for file in $$(ls $$DIR); do \
		echo "Running on $$file..."; \
		/usr/bin/time -f "\nReal %e sec\nTime: %E CPU: %P Mem: %MKB" $$BIN "$$DIR/$$file"; \
		echo; \
	done
	
.PHONY: test-fab-lambda
test-fab-lambda: fab.ref
	@BIN=./fabl && \
	DIR=./inputs && \
	for file in $$(ls $$DIR); do \
		echo "Running on $$file..."; \
		/usr/bin/time -f "\nReal %e sec\nTime: %E CPU: %P Mem: %MKB" $$BIN "$$DIR/$$file"; \
		echo; \
	done
