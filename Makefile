GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
APPS=$(shell ls app)
PROJ_ROOT=$(shell pwd)
ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find app -name *.proto")
	API_PROTO_FILES=$(shell $(Git_Bash) -c 'find proto -path proto/third_party -prune -o -name "*.proto" | grep -v "^proto/third_party"')
else
	INTERNAL_PROTO_FILES=$(shell find app -name *.proto)
	API_PROTO_FILES=$(shell find proto -path proto/third_party -prune -o -name "*.proto" | grep -v "^proto/third_party")
endif

.PHONY: init
# init env. install essential tools and plugins, including goxenith, protoc plugins and wire.
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/favadi/protoc-go-inject-tag@latest

.PHONY: mod
# add missing and remove unused modules
mod:
	go mod tidy

.PHONY: fmt
# go format *.go files
fmt:
	gofmt -s -w .

.PHONY: ent
ent:
	cd app/models/ && go generate ./...

.PHONY: proto

# generate proto
proto:
	protoc --proto_path=./proto \
			--proto_path=./proto/third_party \
 	       --go_out=paths=source_relative:./proto \
 	       --go-http_out=paths=source_relative:./proto \
 	       --go-grpc_out=paths=source_relative:./proto \
	       $(API_PROTO_FILES)
	       find ./proto -name "*.pb.go" -exec protoc-go-inject-tag -input={} \;

.PHONY: build-direct
# build app for Linux
build-direct:
	GOOS=linux GOARCH=amd64 go build -o tmp/goxenith .

.PHONY: build
# build
build: ent proto
	mkdir -p tmp/ && go build -mod=mod -ldflags "-X main" -o ./tmp/ ./...
.PHONY: generate


.PHONY: wire
# generate wire
wire:
		cd $(PROJ_ROOT)/app/cmd/ && wire ./...; \

.PHONY: all
# generate all
all:
	make proto;
	make ent;

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help