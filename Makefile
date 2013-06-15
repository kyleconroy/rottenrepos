.PHONY: build test fmt

GOPATH = $(HOME)/gopath
export GOPATH

REVEL = $(GOPATH)/bin/revel

$(REVEL):
	go get github.com/robfig/revel

build: test
	$(REVEL) build github.com/stackmachine/rottenrepos

test: fmt $(REVEL)
	$(REVEL) test github.com/stackmachine/rottenrepos

fmt:
	gofmt -l -w app tests codespy


serve: fmt $(REVEL)
	$(REVEL) run github.com/stackmachine/rottenrepos
