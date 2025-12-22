GO_VERSION := 1.25

.PHONY: build dist run clean test build-linux-x86 build-macos-arm64

build:
	go build -o dist/file-compressor main.go

dist: build-linux-x86 build-macos-arm64

build-linux-x86:
	GOOS=linux GOARCH=amd64 go build -o dist/file-compressor-linux-x86 main.go

build-macos-arm64:
	GOOS=darwin GOARCH=arm64 go build -o dist/file-compressor-macos-arm64 main.go

run:
	go run main.go

clean:
	rm -rf dist/

test:
	go test ./...
