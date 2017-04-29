FROM golang:1.8.1-alpine

RUN apk update && apk upgrade
RUN apk add git gcc libc-dev

RUN mkdir -p /go/src/github.com/thesoenke/news-crawler
COPY . /go/src/github.com/thesoenke/news-crawler
WORKDIR /go/src/github.com/thesoenke/news-crawler

RUN go get
RUN go build -o /usr/local/bin/news-crawler

COPY feeds /feeds
COPY feed-crawler /etc/periodic/15min/feed-crawler

CMD ["/usr/local/bin/news-crawler --feeds /feeds/news_de.json"]
# CMD crond -l 2 -f