.PHONY: run build install deps core-deps local-deps nodejs-deps assets bindata staticfiles go-bindata

# shorthands:
# - go run -race ./cmd/snk-iris/*.go --config-dir=./shared/conf --resource-dir=`pwd`/shared # (required to check race conditions)
# - go run ./cmd/snk-iris/*.go --config-dir=./shared/conf --resource-dir=`pwd`/shared # (faster but not safer for debugging)
run: bindata
	@go run -race ./cmd/snk-iris/*.go \
		--config-dir=./shared/conf \
		--config-file=application.yml \
		--config-file=api.yml \
		--config-file=server.yml \
		--config-file=websocket.yml \
		--config-file=database.yml \
		--resource-dir=$(CURDIR)/shared

build: assets
	@go build -o ./bin/snk-iris ./cmd/snk-iris/*.go

build-run:
	@./bin/snk-iris \
		--config-dir=./shared/conf \
		--config-file=application.yml \
		--config-file=api.yml \
		--config-file=server.yml \
		--config-file=websocket.yml \
		--config-file=database.yml \
		--resource-dir=$(CURDIR)/shared

install:
	@go install cmd/snk-iris

assets: gz-bindata # go-bindata

gz-bindata:
	@bindata -ignore=\\.DS_Store -pkg main -o ./cmd/snk-iris/bindata.go ./shared/dist/web/...

go-bindata:
	@go-bindata -ignore=\\.DS_Store -pkg main -o ./cmd/snk-iris/bindata.go ./shared/dist/web/...

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