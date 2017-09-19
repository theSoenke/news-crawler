FROM golang:1.9.0-alpine

RUN apk update && apk upgrade
RUN apk add git

RUN mkdir -p /go/src/github.com/thesoenke/news-crawler
COPY . /go/src/github.com/thesoenke/news-crawler
COPY data /app/data

WORKDIR /go/src/github.com/thesoenke/news-crawler
RUN go get
RUN go build -o /usr/local/bin/news-crawler

RUN touch crontab.tmp \
    && echo '30 * * * * /usr/local/bin/news-crawler feeds /app/data/feeds_de.txt --timezone Europe/Berlin --dir /app/out/feeds' > crontab.tmp \
    && echo '0 2 * * * /usr/local/bin/news-crawler scrape /app/out/feeds/ --timezone Europe/Berlin --dir /app/out/articles' >> crontab.tmp \
    && crontab crontab.tmp \
    && rm -rf crontab.tmp

CMD crond -l 5 -f