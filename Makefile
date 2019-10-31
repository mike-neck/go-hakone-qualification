.DEFAULT_GOAL := build
.PHONY: build

build: clean build-96

test: clean hakone-test test-96

clean:
	rm -rf build/

build-96:
	go build -o build/hakone-96 ./cmd/hakone-96/

hakone-test:
	@echo test for Record type
	cd hakone && go test

test-96:
	@echo test for hakone-96
	go test ./cmd/hakone-96/...
