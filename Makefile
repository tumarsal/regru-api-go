MODULE   := $(shell go list -m)
BINARY   := regru
CMD      := ./cmd/regru
GOBIN    ?= $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN := $(shell go env GOPATH)/bin
endif

.PHONY: all build install run test clean help

all: build

help:
	@echo "Targets:"
	@echo "  make install  — установить $(BINARY) в $(GOBIN)"
	@echo "  make build    — собрать bin/$(BINARY)"
	@echo "  make run      — go run $(CMD)"
	@echo "  make test     — go test ./..."
	@echo "  make clean    — удалить bin/"

build:
	go build -o bin/$(BINARY) $(CMD)

install:
	go install $(CMD)

run:
	go run $(CMD) $(ARGS)

test:
	go test ./...

clean:
	rm -rf bin/
