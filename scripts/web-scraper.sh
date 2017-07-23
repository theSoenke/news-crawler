#!/bin/sh

/usr/local/bin/news-crawler feeds /app/data/feeds_de.txt --timezone Europe/Berlin --out /app/out/feeds/
/usr/local/bin/news-crawler scrape /app/out/feeds/ --timezone Europe/Berlin