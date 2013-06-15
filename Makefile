.PHONY: build test fmt

GOPATH = $(shell pwd)
export GOPATH

build: test
	bin/revel build rottenrepos

test: fmt
	bin/revel test rottenrepos

fmt:
	gofmt -l -w src/codespy
	gofmt -l -w src/rottenrepos


serve: fmt
	bin/revel run rottenrepos
