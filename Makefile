.PHONY: setup deps update build 

NAME := NCol
#VERSION := $(shell git describe --tags --abbrev=0)

## setup
setup:
	go get github.com/Masterminds/glide
	go get golang.org/x/tools/cmd/goimports
	go get github.com/golang/lint/golint

deps:
	glide install
update:
	glide update

#test: deps
#	go test $$(glide novendor)

fmt:
	goimports -w $$(glide nv -x)

build:
	go build
