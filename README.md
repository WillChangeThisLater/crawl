**Simple Link Crawler**
This is a simple link crawler, written in go.

The crawler only follows links with the same domain name
as the original link you provide it. The crawler trawls
links concurrently, but can be set to go through links
one at a time.

**Installation**
TODO: figure this out :/

**Features**
I'm working on this link crawler as part of a larger
personal project. It suits my needs for now, so updates
may be infrequent. That said, I have a few feature ideas

- Add support for crawling multiple sites 
- Add more ways of calling the crawler
- Identify duplicate links based on hash value

**How it works**
Most webpages you see on the internet link to one or more webpages.
You can see how they do that by looking at the HTML of the page.
For instance, if you navigate to `eli.thegreenplace.net` (not my site,
but a great programming blog with lots of articles about Go) and hit
`View Page Source`, you'll see HTML which looks something like this

```html
Eli Bendersky's website            </a>
        </div>
        <div class="collapse navbar-collapse navbar-ex1-collapse">
            <ul class="nav navbar-nav navbar-right">
                <li>
                    <a href="https://eli.thegreenplace.net/pages/about">
                        <i class="fa fa-question"></i>
                        <span class="icon-label">About</span>
                    </a>
                </li>
                <li>
                    <a href="https://eli.thegreenplace.net/pages/projects">
                        <i class="fa fa-github"></i>
                        <span class="icon-label">Projects</span>
                    </a>
                </li>
                <li>
                    <a href="https://eli.thegreenplace.net/archives/all">
                        <i class="fa fa-th-list"></i>
                        <span class="icon-label">Archives</span>
                    </a>
                </li>
            </ul>
        </div>
        <!-- /.navbar-collapse -->
    </div>
</div> <!-- /.navbar -->
```

Those href attributes are the links to other pages. Often these are just other pages
on the website, but they can also link to things you don't normally see directly which a
web browser uses behind the scenes (CSS, JavaScript, etc.)

What a crawler does is relatively simple. It starts with a link to a page. It downloads
the HTML for the page and parses the HTML looking for those links. It then navigates to
each link it finds and repeats the same process.

One important note: it's critical that the crawler maintains some kind of information about
the links it has already seen. That's because two pages can link to one another. If the crawler
doesn't keep track of the links it has seen it may get caught in an infinite loop, going back
and forth between two pages forever.

My crawler is very bare bones. It's intended to run on a single local machine, and scrape
a single site. It has some support for concurrency (e.g. it can scrape multiple pages at
the same time), but it doesn't have much else.

If you're curious and want to see what a more advanced web crawler looks like
check out [colly](https://github.com/gocolly/colly). The colly README also has a cool
list of projects that use scraping/crawling if you're interested.
