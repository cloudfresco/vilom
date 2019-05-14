
SHELL := /bin/bash

# The name of the executable (default is current directory name)
#TARGET := $(shell echo $${PWD\#\#*/})
#.DEFAULT_GOAL: $(TARGET)

# These will be provided to the target
#VERSION := 1.0.0
#VERSION          := $(shell git describe --tags --always --dirty="-dev")
#DATE             := $(shell date -u '+%Y-%m-%d-%H%M UTC')
#VERSION_FLAGS    := -ldflags='-X "main.Version=$(VERSION)" -X "main.BuildTime=$(DATE)"'
#BUILD := `git rev-parse HEAD`

# Use linker flags to provide version/build settings to the target
#LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
MFILE = cmd/main.go
EXEC = cmd/vilom.o
PKGS = ./...
.PHONY: all deps build buildp test clean fmt vet lint err sql run runp doc

all: 

chk: fmt vet lint err

deps: 
	@dep ensure

build: 
	@go build -i -o $(EXEC) $(MFILE)

buildp:
	@go build -i -o $(EXEC) $(MFILE)

test: deps buildp
	@go test -v $(PKGS)

clean:
	@rm -f $(EXEC)

fmt:
	@gofmt -s -l -w $(SRC)

vet:
	@go vet $(PKGS)

linter:
	@go get -u github.com/golang/lint/golint

lint: linter
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d}; done 

errcheck: 
	@go get github.com/kisielk/errcheck

err: errcheck 
	@errcheck $(PKGS)

safesql:
	@go get github.com/stripe/safesql

sql: safesql
	@safesql $(SRC)

run: build	
	@./$(EXEC) --dev 

runp: buildp	
	@./$(EXEC) 

doc: 

