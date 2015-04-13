VERSION = `cat ./version.go | grep "Version = " | cut -d" " -f4 | sed 's/[^"]*"\([^"]*\).*/\1/'`
DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

all: deps build

deps:
	go get -d -v ./...
	echo $(DEPS) | xargs -n1 go get -d

build:
	go install

test: deps
	go list ./... | xargs -n1 go test -timeout=3s

.PHONY: all deps build test
