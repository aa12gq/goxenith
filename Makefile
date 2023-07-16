GOPATH:=$(shell go env GOPATH)
.PHONY: ent
ent:
	cd app/ && go generate ./...