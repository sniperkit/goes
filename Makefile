.PHONY: run build install deps

run:
	@go run -race *.go

build:
	@go build -o ./bin/goes *.go

install:
	@go install

deps:
	@glide install --strip-vendor