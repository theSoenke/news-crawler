# news-crawler

## Build

```
git clone https://github.com/thesoenke/news-crawler
cd news-crawler
make
```

## Run in Docker
```
docker build -t news-crawler .
docker run -t -v $PWD/out:/app/out news-crawler
```