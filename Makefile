GOBIN := go
PROTOCBIN := protoc
SERVER := wuserver
CLIENT := wuclient
VERSION := v0.4
LDFLAGS := -ldflags="-s -w -X \"github.com/mas9612/wrapups/pkg/version.Version=$(VERSION)\""

.PHONY: all
all: dep test build-grpc build doc

.PHONY: build
build: build-server build-client

.PHONY: build-server
build-server:
	CGO_ENABLED=0 $(GOBIN) build $(LDFLAGS) -o $(SERVER) ./cmd/wuserver

.PHONY: build-client
build-client:
	CGO_ENABLED=0 $(GOBIN) build $(LDFLAGS) -o $(CLIENT) ./cmd/wuclient

.PHONY: test
test:
	$(GOBIN) test -v ./...

.PHONY: build-grpc
build-grpc:
	$(PROTOCBIN) --go_out=plugins=grpc:. ./pkg/wrapups/wrapups.proto

.PHONY: doc
doc:
	$(PROTOCBIN) --doc_out=./doc --doc_opt=markdown,wrapups.md ./pkg/wrapups/wrapups.proto

.PHONY: clean
clean:
	$(GOBIN) clean
	rm -f $(SERVER)
	rm -f $(CLIENT)

.PHONY: dep
dep:
	dep ensure
