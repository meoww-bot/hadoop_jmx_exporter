SHELL := /bin/bash

REV := $(shell git rev-parse HEAD)
CHANGES := $(shell test -n "$$(git status --porcelain)" && echo '-CHANGES' || true)
VERSION := $(shell cat ./VERSION)


.PHONY: \
	clean \
	test \
	vet \
	lint \
	build

all: fmt vet build

test:
	go test -v ./

vet:
	go vet -v ./

lint:
	golint ./

style:
	gofmt -d ./

build: 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -o bin/hadoop_jmx_exporter .