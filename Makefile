DEST=dist
BIN=$(DEST)/k8s-res

.PHONY: build all clean test cover

default: all

all: build

build:
	@mkdir -p $(DEST)
	go build -o ./$(BIN) ./main.go

test:
	go test -v ./...
	
cover:
	go test -v -cover -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out

clean:
	rm -r ./$(DEST)