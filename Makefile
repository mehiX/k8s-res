NAME		:= k8s-res
DEST		:= dist
OUTPUT_BIN	:= $(DEST)/$(NAME)
PACKAGE		:= $(shell grep -E "^module" go.mod | cut -d ' ' -f 2)
GIT_REV		?= $(shell git rev-parse --short HEAD)
SOURCE_DATE_EPOCH ?= $(shell date +%s)
ifeq ($(shell uname), Darwin)
DATE       ?= $(shell TZ=UTC date -j -f "%s" ${SOURCE_DATE_EPOCH} +"%Y-%m-%dT%H:%M:%SZ")
else
DATE       ?= $(shell date -u -d @${SOURCE_DATE_EPOCH} +"%Y-%m-%dT%H:%M:%SZ")
endif
VERSION    ?= v0.1

.PHONY: build all clean test cover

default: all

all: build

build:
	@echo *********** BUILD EXEC ***************
	@mkdir -p $(DEST)
	@go build \
	-ldflags "-w -s -X ${PACKAGE}/cmd.version=${VERSION} -X ${PACKAGE}/cmd.commit=${GIT_REV} -X ${PACKAGE}/cmd.date=${DATE}" \
	-a -o ./$(OUTPUT_BIN) ./main.go

test:
	go test -v ./...
	
cover:
	go test -v -cover -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out -o=coverge.html

clean:
	@echo *********** CLEANUP ***************
	@rm -r ./$(DEST)
