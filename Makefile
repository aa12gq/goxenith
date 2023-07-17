GOPATH:=$(shell go env GOPATH)
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
.PHONY: ent
ent:
	cd app/models/ && go generate ./...

.PHONY: api

# generate api proto
api:
	protoc --proto_path=./proto \
			--proto_path=./proto/third_party \
 	       --go_out=paths=source_relative:./proto \
 	       --go-http_out=paths=source_relative:./proto \
 	       --go-grpc_out=paths=source_relative:./proto \
	       $(API_PROTO_FILES)
	       find ./proto -name "*.pb.go" -exec protoc-go-inject-tag -input={} \;
