#!/bin/bash

main() {
  cat <<EOF
I have the following project:

$(tree)
$(files-to-prompt . --ignore test-sites)

How can I turn this into a CLI tool? I want
the interface to simply be:

\`\`\`bash
$ crawl --input <url_list>              # if url list is provided, open file and run crawler on each URL
$ echo "http://www.example.com" | crawl # urls can also be provided via stdin. assume one url per line. strip whitespace before crawling
\`\`\`

It looks like currently the package name is 'crawl'. I suppose this makes
it easy to use the package as part of other projects like so:

\`\`\`go
import (
        "bufio"
        "flag"
        "fmt"
        "log"
        "os"
        "strings"

        "github.com/WillChangeThisLater/crawl"
)
\`\`\`

That said, even though I _do_ want to be able to impact crawl as a package as part of other
go projects, I also want to be able to use it as a standalone CLI script as discussed above
EOF
}

main
