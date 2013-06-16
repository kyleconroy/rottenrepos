.PHONY: build test fmt

GOPATH = $(HOME)/gopath
export GOPATH

REVEL = $(GOPATH)/bin/revel

$(REVEL): install

install:
	go get github.com/robfig/revel
	go get github.com/robfig/revel/revel
	go get github.com/garyburd/redigo/redis
	go get github.com/soveran/redisurl

build: test
	$(REVEL) build github.com/stackmachine/rottenrepos

test: fmt $(REVEL)
	$(REVEL) test github.com/stackmachine/rottenrepos

fmt:
	gofmt -l -w .


serve: fmt $(REVEL)
	$(REVEL) run github.com/stackmachine/rottenrepos
