FROM golang:1.15.3-alpine

RUN apk update && apk upgrade
RUN apk add git

RUN mkdir -p /go/src/github.com/thesoenke/news-crawler
COPY . /go/src/github.com/thesoenke/news-crawler
COPY data /app/data

WORKDIR /go/src/github.com/thesoenke/news-crawler
RUN go build -o /usr/local/bin/news-crawler
RUN crontab crontab

CMD crond -l 5 -f
