.PHONY: all

all: build

build:
	go build -ldflags='-linkmode=external'

test:
	go test -ldflags='-linkmode=external'

bench:
	go test -test.bench . -ldflags='-linkmode=external'