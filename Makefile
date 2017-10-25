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
	@docker-compose exec crawler /usr/local/bin/news-crawler nod --lang german --dir /app/out/nod --logs /app/out/events.log --timezone Europe/Berlin

nod-yesterday:
	@docker-compose run crawler /usr/local/bin/news-crawler nod --from yesterday --lang german --dir /app/out/nod --logs /app/out/events.log --timezone Europe/Berlin