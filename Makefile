GOBIN=go
PROTOCBIN=protoc
BINNAME=wrapups

all: dep test build-grpc build doc

build:
	$(GOBIN) build ./cmd/wrapups -o $(BINNAME)

test:
	$(GOBIN) test -v ./...

build-grpc:
	$(PROTOCBIN) --go_out=plugins=grpc:. ./pkg/wrapups/wrapups.proto

doc:
	$(PROTOCBIN) --doc_out=./doc --doc_opt=markdown,wrapups.md ./pkg/wrapups/wrapups.proto

clean:
	$(GOBIN) clean
	rm -f $(BINNAME)

dep:
	dep ensure
