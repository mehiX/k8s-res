BIN=dist

.PHONY: build all clean test

default: all

all: build

build:
	@mkdir -p $(BIN)
	go build -o ./$(BIN)/k8s-resources ./main.go

test:
	go test -v -cover -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out

clean:
	rm -r ./$(BIN)