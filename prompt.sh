#!/bin/bash

main() {
  cat <<EOF
I have the following project:

$(tree)
$(files-to-prompt . --ignore test-sites)

Right now this project is called 'crawl'. Should
I rename it to 'scrape'?
EOF
}

main
