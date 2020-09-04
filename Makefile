
NAME := wpdump
VERSION := v0.0.1
REVISION := $(shell git rev-parse --short HEAD)

##
# options
##

GOARCH := amd64
LDFLAGS := -X 'main.Name=$(NAME)' \
           -X 'main.Version=$(VERSION)' \
           -X 'main.Revision=$(REVISION)'

ifeq ($(DEBUG), 1)
	BUILD_OPTIONS := -race -tags DEBUG -ldflags="$(LDFLAGS)"
	BUILD_MODE := debug
else
	BUILD_OPTIONS := -ldflags="-s -w $(LDFLAGS)"
	BUILD_MODE := release
endif

DOCKER_GO := docker run -it -v "$(PWD):/go" -e GOPATH= -e GOOS=$$GOOS -e GOARCH=$$GOARCH golang:latest go
GO_SRCS := $(shell find . -type f -name '*.go')

##
# build
##

all: linux darwin windows
linux: target/$(BUILD_MODE)/$(NAME)-linux-$(GOARCH)
darwin: target/$(BUILD_MODE)/$(NAME)-darwin-$(GOARCH)
windows: target/$(BUILD_MODE)/$(NAME)-windows-$(GOARCH).exe

##
# artifact
##

target/$(BUILD_MODE)/$(NAME)-linux-$(GOARCH): $(GO_SRCS)
	GOOS=linux; \
	GOARCH=$(GOARCH); \
	$(DOCKER_GO) build $(BUILD_OPTIONS) -o target/$(BUILD_MODE)/$(NAME)-linux-$(GOARCH)

target/$(BUILD_MODE)/$(NAME)-darwin-$(GOARCH): $(GO_SRCS)
	GOOS=darwin; \
	GOARCH=$(GOARCH); \
	$(DOCKER_GO) build $(BUILD_OPTIONS) -o target/$(BUILD_MODE)/$(NAME)-darwin-$(GOARCH)

target/$(BUILD_MODE)/$(NAME)-windows-$(GOARCH).exe: $(GO_SRCS)
	GOOS=windows; \
	GOARCH=$(GOARCH); \
	$(DOCKER_GO) build $(BUILD_OPTIONS) -o target/$(BUILD_MODE)/$(NAME)-windows-$(GOARCH).exe

clean:
	-rm -rf target/*

##
# check
##

check:
	@echo [tool fix]
	@go tool fix -diff .

	@echo
	@echo [fmt]
	@go fmt ./...

	@echo
	@echo [gofumpt]
	@go get -u mvdan.cc/gofumpt
	@gofumpt -l -w .

	@echo
	@echo [golangci-lint]
	@golangci-lint run --enable 'asciicheck,bodyclose,depguard,dogsled,dupl,exhaustive,exportloopref,funlen,gochecknoinits,gocognit,goconst,gocritic,gocyclo,godot,godox,goerr113,gofmt,gofumpt,goheader,goimports,golint,gomnd,gomodguard,goprintffuncname,interfacer,misspell,nakedret,nestif,nlreturn,noctx,nolintlint,prealloc,rowserrcheck,scopelint,sqlclosecheck,stylecheck,unconvert,unparam,whitespace,wsl'

	@echo
	@echo [mod tidy]
	@go mod tidy

	@echo
	@echo [mod verify]
	@go mod verify

	@echo
	@echo [test]
	@go test ./...

.PHONY: all linux darwin windows clean
