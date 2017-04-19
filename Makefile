.PHONY: vendor build
PACKAGES = $(shell go list ./...)

default: build

all: build test vet

test:
	@go test ${PACKAGES}

vet:
	@go vet ${PACKAGES}

build:
	@go get ${PACKAGES}

run:
	@go get ${PACKAGES}
	@${GOPATH}/bin/news-crawler