BIN=dist

.PHONY: build all clean

default: all

all: build

build:
	@mkdir -p $(BIN)
	go build -o ./$(BIN)/k8sres ./main.go

clean:
	rm -r ./$(BIN)