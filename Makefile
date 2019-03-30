GOBIN=go
PROTOCBIN=protoc
SERVER=wuserver
CLIENT=wuclient

all: dep test build-grpc build doc

build:
	$(GOBIN) build -o $(SERVER) ./cmd/wuserver
	$(GOBIN) build -o $(CLIENT) ./cmd/wuclient

test:
	$(GOBIN) test -v ./...

build-grpc:
	$(PROTOCBIN) --go_out=plugins=grpc:. ./pkg/wrapups/wrapups.proto

doc:
	$(PROTOCBIN) --doc_out=./doc --doc_opt=markdown,wrapups.md ./pkg/wrapups/wrapups.proto

clean:
	$(GOBIN) clean
	rm -f $(SERVER)

dep:
	dep ensure
