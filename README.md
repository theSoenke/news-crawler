# news-crawler

Uses a list of feeds to crawl and store all daily published news.

## Docker Compose Setup
1. `git clone github.com/thesoenke/news-crawler`
2. `cd news-crawler`
3. `docker-compose up --build`

This will start the crawler, Elastisearch and Kibana. In case ElasticSearch is crashing you might need to increase the virtual memory limit for Elasticsearch

    sysctl -w vm.max_map_count=262144

When everything worked 3 containers should be running

## Local setup
### Install
1. Make sure [go](https://golang.org) is installed
2. `git clone github.com/thesoenke/news-crawler`
3. `cd news-crawler && make`
3. `$GOPATH/bin` should be in your PATH or run it directly with `$GOPATH/bin/news-crawler`

### Run
#### Feedreader
Download articles from a list of feeds

    news-crawler feeds data/feeds_de.txt --lang german

#### Scraper
The scraper downloads articles found by the feedreader and indexes them in Elasticsearch.

    news-crawler scrape out/feeds/german/26-6-2017.json --lang german

An ElasticSearch instance needs to be running. Credentials can be set with the environment variables `ELASTIC_URL`, `ELASTIC_USER` and `ELASTIC_PASSWORD`. Defaults are url: `http://localhost:9200`, user: `elastic`, password: `changeme`

#### Create NoDCore input

    newscrawler nod --lang german

Generated files can be found in `out/nod/<lang>` by default.

## Using Kibana
When using the docker-compose setup open `http://localhost:5061` and add the index mapping. Index `news-*` will contain all languages. Language specific index mapping can be created by using `news-<lang>`. The index `failures-*` logs all failures. \
Warning: The scraper has to be run at least once to create the indices
