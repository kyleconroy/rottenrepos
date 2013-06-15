.PHONY: build test fmt

GOPATH = $(shell pwd)
export GOPATH

test: fmt
	go test rotten

fmt:
	go fmt rotten
