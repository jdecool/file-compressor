GO_VERSION := 1.25

.PHONY: build run clean test

build:
	go build -o file-compressor main.go

run:
	go run main.go

clean:
	rm -f file-compressor

test:
	go test ./...
