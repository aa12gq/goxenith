GOPATH:=$(shell go env GOPATH)
.PHONY: ent
ent:
	cd app/model/ && go generate ./...