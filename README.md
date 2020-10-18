# news-crawler

This project is crawling, extracting, indexing and processing the content of daily published news articles. The extracted content is indexed in ElasticSearch for further processing. This project also provides tooling to extract and preprocess the content for the [NoD](https://github.com/uhh-lt/NoDCore) project.

## Docker Compose
1. `git clone github.com/thesoenke/news-crawler`
2. `cd news-crawler`
3. `docker-compose up --build`

This will start the crawler, Elastisearch and Kibana. In case ElasticSearch is crashing you might need to increase the virtual memory limit for Elasticsearch

    sysctl -w vm.max_map_count=262144

When everything worked 3 containers should be running.

It is also possible to generate the NoDCore input in docker by running `make nod-docker`. To generate only the data for the day before run `make nod-yesterday`. In case you want to run `make nod-yesterday` in a cronjob you have to provide the full path to docker-compose in the Makefile e.g. `/usr/local/bin/docker-compose`. The output will be available in `out/nod/german`. This will only work when the docker-compose setup is already running and the scraper has run at least once.

## Docker Compose NoD
Starts the crawler and [NoDWeb](https://github.com/uhh-lt/NoDWeb) and automatically runs [NoDCore](https://github.com/uhh-lt/NoDCore) when new data is available.
1. `git clone github.com/thesoenke/news-crawler`
2. `cd news-crawler`
3. `docker-compose -f docker-compose.nod.yml up --build`
4. Open [localhost:9000](localhost:9000)

## Local setup
### Install
1. Make sure [go](https://golang.org) is installed
2. `go get github.com/thesoenke/news-crawler`
3. `$GOPATH/bin` should be in your `PATH` or run it directly with `$GOPATH/bin/news-crawler`

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

Generated files can be found in `out/nod/<lang>` by default. This command will only work when the scraper has run at least once to insert data into ElasticSearch.

It is also possible to only output the content from yesterday with the following command:

    news-crawler nod --from yesterday --lang german

## Logs
- Successful runs of the feedreader, scraper and nod exporter are logged with in `out/events.log`
- The feedreader writes a log of feeds that could not be fetched to `out/feeds/<lang>/failures.log`
- Articles that could not be fetched are logged in the ElasticSearch index `failures-<lang>`

## Archive
All fetched webpages are by default stored on disk in `out/content/<lang>/<date>/<md5sum of the url>.html`

## Using Kibana
When using the docker-compose setup open 4. Open [localhost:5601](localhost:5601) and add the index mapping. Index `news-*` will contain all languages. Language specific index mapping can be created by using `news-<lang>`. The index `failures-<lang>` logs all failures. \
Warning: The scraper has to be run at least once to create the indices

## Adding a new language
Before a new language can be added a list of feeds is required. After that the 3 environment variables need to be change in `docker-compose.yml`
- `CRAWLER_FEEDS_FILE` Path to a file with a list of feeds
- `CRAWLER_LANGUAGE` Language of the parser
- `CRAWLER_TIMEZONE` Timezone that should be used when storing the article publish dates
