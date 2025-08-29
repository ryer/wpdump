
NAME := wpdump
REVISION := $(shell git rev-parse --short HEAD)

##
# options
##

GOARCH := amd64
LDFLAGS := -X 'main.Name=$(NAME)' \
           -X 'main.Revision=$(REVISION)'

ifeq ($(DEBUG), 1)
	BUILD_OPTIONS := -gcflags=all="-N -l" -tags DEBUG -ldflags="$(LDFLAGS)"
	BUILD_MODE := debug
else
	BUILD_OPTIONS := -ldflags="-s -w $(LDFLAGS)"
	BUILD_MODE := release
endif

DOCKER_GO := docker run -it -v "$(PWD):/go" -e GOPATH= -e GOOS=$$GOOS -e GOARCH=$$GOARCH golang:1.24 go
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
	@echo [go fix]
	@go fix ./...

	@echo
	@echo [go fmt]
	@go fmt ./...

	@echo
	@echo [gofumpt]
	@go run mvdan.cc/gofumpt -l -w .

	@echo
	@echo [golangci-lint]
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run --enable-all --disable 'wsl,lll,varnamelen,ireturn,cyclop,paralleltest,wrapcheck,gosec,forbidigo,testpackage,gochecknoglobals,depguard,exhaustruct,intrange,mnd,tenv,funlen'

	@echo
	@echo [go mod tidy]
	@go mod tidy

	@echo
	@echo [go mod verify]
	@go mod verify

	@echo
	@echo [go test]
	@go test ./...

.PHONY: all linux darwin windows clean check
