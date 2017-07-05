# news-crawler

## Build
```
git clone https://github.com/thesoenke/news-crawler
cd news-crawler
make
```

## Run
Download feeds provided in an input file
```
news-crawler feeds --file data/feeds_de.txt
```

Download articles scraped by the feed downloader
```
export ELASTIC_URL="http://localhost:9200"
export ELASTIC_USER=elastic
export ELASTIC_PASSWORD=changeme
news-crawler scrape --file out/feeds/26-6-2017.json
```

## Run with Docker Compose
Run crawler and elasticsearch
```
docker-compose up
```

You might need to increase the virtual memory limit when elasticsearch does not start
```
sysctl -w vm.max_map_count=262144
```