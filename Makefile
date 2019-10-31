.DEFAULT_GOAL := build
.PHONY: build

build: clean build-96

test: clean test-96

clean:
	rm -rf build/

build-96:
	go build -o build/hakone-96 ./cmd/hakone-96/

test-96:
	go test ./cmd/hakone-96/...
