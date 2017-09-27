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

nod-docker:
	@docker-compose exec crawler /usr/local/bin/news-crawler nod --lang german --dir /app/out/nod