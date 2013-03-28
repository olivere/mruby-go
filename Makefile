.PHONY: all mruby

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
INC := ./include
LIB := ./lib/${GOOS}_${GOARCH}

all: build

build:
	C_INCLUDE_PATH=${INC} LIBRARY_PATH=${LIB} go build -ldflags='-linkmode=external'

test:
	C_INCLUDE_PATH=${INC} LIBRARY_PATH=${LIB} go test -ldflags='-linkmode=external'

bench:
	C_INCLUDE_PATH=${INC} LIBRARY_PATH=${LIB} go test -test.bench . -ldflags='-linkmode=external'

mruby:
	pushd mruby && make && popd
