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

## Run in Docker
```
docker build -t news-crawler .
docker run -t -v $PWD/out:/app/out news-crawler
```