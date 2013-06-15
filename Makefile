.PHONY: build test fmt

GOPATH = $(shell pwd)
export GOPATH

build: test
	bin/revel build rottenrepos

test: fmt
	bin/revel test rottenrepos

fmt:
	go fmt codespy
	go fmt rottenrepos/app rottenrepos/tests


serve: fmt
	bin/revel run rottenrepos
