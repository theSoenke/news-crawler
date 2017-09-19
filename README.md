# news-crawler

Uses a list of feeds to crawl and store all daily published news

## Install
1. Make sure [go](https://golang.org) is installed
2. `go get github.com/thesoenke/news-crawler`
3. `$GOPATH/bin` should be in your PATH or run directly `$GOPATH/bin/news-crawler`


## Run
Download articles from a list of feeds

    news-crawler feeds data/feeds_de.txt


Download articles found by the feed downloader and index them in Elasticsearch

    export ELASTIC_URL="http://localhost:9200"
    export ELASTIC_USER=elastic
    export ELASTIC_PASSWORD=changeme
    news-crawler scrape out/feeds/26-6-2017.json

Create NodCore input

    newscrawler nod --lang german


## Docker Compose
Start the crawler, Elastisearch and Kibana

    docker-compose up

You might need to increase the virtual memory limit for Elasticsearch

    sysctl -w vm.max_map_count=262144

When everything worked 3 containers should be running

## Using Kibana
Open localhost:5061 and add the index mapping. Index `news-*` will contain all languages. Language specific index mapping can be created by using `news-<lang>`. The index `failures-*` logs all failures. \
Warning: The scraper has to be run at least once to create the indices
