GOPATH:=$(shell go env GOPATH)
.PHONY: ent
ent:
	cd app/models/ && go generate ./...