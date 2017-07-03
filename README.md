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
news-crawler scrape --file out/feeds/26-6-2017.json
```

## Run with Docker Compose
Docker Compose start a crawler and elasticsearch container
```
docker-compose up
```

You might need to increase the virtual memory limit when elasticsearch does not start
```
sysctl -w vm.max_map_count=262144
```