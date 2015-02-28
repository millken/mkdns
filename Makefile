export GOPATH=$(cd "$(dirname "$0")"; pwd)

BINDATA_IGNORE = $(shell git ls-files -io --exclude-standard $< | sed 's/^/-ignore=/;s/[.]/[.]/g')

usage:
	@echo ""
	@echo "Task                 : Description"
	@echo "-----------------    : -------------------"
	@echo "make setup           : Install all necessary dependencies"
	@echo "make dev             : Generate development build"
	@echo "make test            : Run tests"
	@echo "make format          : formater code"	
	@echo "make build           : Generate production build for current OS"
	@echo "make bootstrap       : Install cross-compilation toolchain"
	@echo "make release         : Generate binaries for all supported OSes"
	@echo "make clean           : Remove all build files and reset assets"
	@echo "make assets          : Generate production assets file"
	@echo "make dev-assets      : Generate development assets file"
	@echo "make docker          : Build docker image"
	@echo ""

test:
	godep go test

assets: static/
	go-bindata $(BINDATA_OPTS) $(BINDATA_IGNORE) -ignore=[.]gitignore -ignore=[.]gitkeep $<...

dev-assets:
	@$(MAKE) --no-print-directory assets BINDATA_OPTS="-debug"

dev: dev-assets
	godep go build
	@echo "You can now execute ./mkdns"

format:	
	go fmt ./...
build: 
	godep go build
	@echo "You can now execute ./mkdns"

release: assets
	gox -osarch="darwin/amd64 darwin/386 linux/amd64 linux/386 windows/amd64 windows/386" -output="./bin/mkdns_{{.OS}}_{{.Arch}}"

bootstrap:
	gox -build-toolchain

setup:
	go get github.com/tools/godep
#	godep get github.com/mitchellh/gox
#	godep get github.com/jteeuwen/go-bindata/...
	godep restore

clean:
	rm -f ./mkdns
	rm -f ./bin/*
	rm -f bindata.go

docker:
	docker build -t mkdns .
