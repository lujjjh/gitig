MODULE=github.com/lujjjh/gitig/cmd/gitig

all: build

build:
	go build -o bin/gitig $(MODULE)

run:
	go run -race $(MODULE)

test:
	go test -race ./...

.PHONY: all build run test
