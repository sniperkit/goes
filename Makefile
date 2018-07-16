.PHONY: run build install deps core-deps local-deps

run:
	@go run -race ./cmd/snk-goes/*.go \
		--config-dir=./shared/conf \
		--config-file=config.yaml \
		--resource-dir=./shared/data

build:
	@go build -o ./bin/snk-goes ./cmd/snk-goes/*.go

install:
	@go install cmd/snk-goes

deps: local-deps

core-deps:
	@go get -u github.com/jteeuwen/go-bindata/...
	@go get -u github.com/mitchellh/gox
	@go get -u github.com/Masterminds/glide

local-deps:
	@glide install --strip-vendor