# `crawl`
Go package + CLI for scraping links from a website

# CLI
## Setup
```bash
# clone the repo and cd in
git clone https://github.com/WillChangeThisLater/crawl
cd crawl

# this will build the CLI package
# and link it into your path
./build.sh
ln -s "$(pwd)/crawl /usr/local/bin/crawl
```

## Usage
Most basic usage
```bash
$ crawl https://www.example.com
https://www.example.com
```

More advanced

- `-d` tells `crawl` to only scrape one level deep
- `-c 25` tells crawl to run up to 25 GET requests at a time
- `-s` silences errors
```bash
$ crawl -d 1 -c 25 -s https://www.google.com
https://www.google.com
https://www.google.com/intl/en/policies/terms/
https://www.google.com/intl/en/policies/privacy/
https://www.google.com/imghp?hl=en&tab=wi
https://www.google.com/preferences?hl=en
https://www.google.com/advanced_search?hl=en&authuser=0
https://www.google.com/intl/en/about.html
https://www.google.com/intl/en/about/products?tab=wh
http://www.google.com/history/optout?hl=en
https://www.google.com/intl/en/ads/
https://www.google.com/services/
```

You can always ask for help
```bash
$ crawl -h
Usage of crawl:
  -c int
        Maximum number of concurrent requests (default 25)
  -d int
        Maximum depth to crawl (use -1 for no limit) (default 1)
  -s    Silence stderr output
  -t int
        Request timeout (default 5)
```
