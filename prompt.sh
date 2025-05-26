#!/bin/bash

set -euo pipefail

reference_links=(
  ""
)

# Function to display references in a readable manner
references() {
  echo "# Reference Index"
  for reference_link in "${reference_links[@]}"; do
    # Print a header with Markdown style
    echo -e "\n## Reference: $reference_link\n"
    lynx -dump -nolist "$reference_link"
    echo -e "\n"
  done
}

about() {
    cat <<EOF

Directory structure
$(tree -I 'test-sites**' .)

Files
$(files-to-prompt . --ignore test-sites --ignore crawl --ignore prompt.sh)

EOF
}

run() {

    echo "$@" >&2

    echo "\`\`\`bash"
    echo "\$ $@"
    $@ 2>&1
    echo "\`\`\`"
}

main() {
  cat <<EOF
About:
$(about)

References:
$(references)

Some sites have references to different parts of the same page:
For instance:

$(run crawl https://eli.thegreenplace.net/)

Is there an intelligent approach for filtering these, so each page is
represented in the output only once? If there is, explain the
approach and implement it
EOF
}

main
