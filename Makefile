.PHONY: help dev clean build

all: clean build

build:
	@mkdir -p bin
	go build ${BUILD_FLAGS} -C bin ..

dev:
	go run . ${FLAGS}

clean:
	rm -rf bin

help:
	@npx -y makeman help.json ${TARGET}