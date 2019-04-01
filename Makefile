GOBIN=go
PROTOCBIN=protoc
SERVER=wuserver
CLIENT=wuclient

.PHONY: all
all: dep test build-grpc build doc

.PHONY: build
build: build-server build-client

.PHONY: build-server
build-server:
	$(GOBIN) build -o $(SERVER) ./cmd/wuserver

.PHONY: build-client
build-client:
	$(GOBIN) build -o $(CLIENT) ./cmd/wuclient

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
