.PHONY: run build install deps core-deps local-deps nodejs-deps assets bindata staticfiles go-bindata

run: bindata
	@go run -race ./cmd/snk-goes/*.go \
		--config-dir=./shared/conf \
		--config-file=config.yaml \
		--resource-dir=./shared/data

build: assets
	@go build -o ./bin/snk-goes ./cmd/snk-goes/*.go

build-run:
	@./bin/snk-goes \
		--config-dir=./shared/conf \
		--config-file=config.yaml \
		--resource-dir=./shared/data

install:
	@go install cmd/snk-goes

assets: gz-bindata # go-bindata

gz-bindata:
	@bindata -ignore=\\.DS_Store -pkg main -o ./cmd/snk-goes/bindata.go ./shared/dist/web/...

go-bindata:
	@go-bindata -ignore=\\.DS_Store -pkg main -o ./cmd/snk-goes/bindata.go ./shared/dist/web/...

web: web-deps

web-deps:
	@cd ./web && yarn install

web-dev:
	@cd ./web && yarn run dev

web-run:
	@cd ./web && yarn run start

web-generate:
	@cd ./web && yarn run generate

deps: local-deps

nodejs-deps:
	@brew install yarn

core-deps:
	@go get -u github.com/kataras/bindata/cmd/bindata
	@go get -u github.com/bouk/staticfiles
	@go get -u github.com/jteeuwen/go-bindata/...
	@go get -u github.com/mitchellh/gox
	@go get -u github.com/Masterminds/glide

local-deps:
	@glide install --strip-vendor