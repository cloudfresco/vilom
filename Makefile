
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
EXEC = cmd/vilom
PKGS = ./...
.PHONY: all build buildp test clean fmt vet lint err sql run runp doc

all: chk buildp

chk: fmt vet lint err

build: 
	@echo "Building vilom"	
	@go build -i -o $(EXEC) $(MFILE)

buildp:
	@echo "Building vilom"	
	@go build -i -o $(EXEC) $(MFILE)

test:
	@mysql -uroot -p$(VILOM_DBPASSROOT) -e 'DROP DATABASE IF EXISTS  $(VILOM_DBNAME_TEST);'
	@mysql -uroot -p$(VILOM_DBPASSROOT) -e 'CREATE DATABASE $(VILOM_DBNAME_TEST);'
	@mysql -uroot -p$(VILOM_DBPASSROOT) -e "GRANT ALL ON *.* TO '$(VILOM_DBUSER_TEST)'@'$(VILOM_DBHOST)';"
	@mysql -uroot -p$(VILOM_DBPASSROOT) -e 'FLUSH PRIVILEGES;'
	@mysql -u$(VILOM_DBUSER_TEST) -p$(VILOM_DBPASS_TEST)  $(VILOM_DBNAME_TEST) < sql/mysql/mysqldb.sql


	@echo "Starting tests"	
	@go test -v $(PKGS)

clean:
	@rm -f $(EXEC)

fmt:
	@echo "Running gofmt"	
	@gofmt -s -l -w $(SRC)

vet:
	@echo "Running vet"	
	@go vet $(PKGS)

linter:
	@go get -u golang.org/x/lint/golint

lint: linter
	@echo "Running lint"
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d}; done 

errcheck: 
	@go get -u github.com/kisielk/errcheck

err: errcheck 
	@echo "Running errcheck"
	@errcheck $(PKGS)

safesql:
	@go get -u github.com/stripe/safesql

sql: safesql
	@echo "Running safesql"
	@safesql $(SRC)

run: build
	@echo "Starting vilom"	
	@./$(EXEC) --dev 

runp: buildp	
	@echo "Starting vilom"	
	@./$(EXEC) 

doc: 

