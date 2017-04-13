FROM golang:1.8.1-alpine

RUN apk update && apk upgrade
RUN apk add git

RUN mkdir -p /go/src/github.com/thesoenke/news-crawler

COPY . /go/src/github.com/thesoenke/news-crawler

WORKDIR /go/src/github.com/thesoenke/news-crawler

ARG REVISION=HEAD

RUN go get
RUN go build -o /usr/local/bin/news-crawler
ENTRYPOINT ["/usr/local/bin/news-crawler"]