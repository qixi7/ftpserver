PWD      := $(realpath .)
GOPATH   := $(realpath $(PWD)/../..)
APPBIN   := ftpserver
GOOS	:=linux
GOARCH	:=amd64
CGO_ENABLED:= 0

.PHONY: build
build:
	GOPATH=$(GOPATH) CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(APPBIN)
