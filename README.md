# Crawler

Can be run with:
  - Local go: go run src/crawler/main.go http://wiprodigital.com/
  - Binary: ./crawler http://wiprodigital.com/
  - Docker: make
    - This will not terminate when finished, just output `Crawl complete` 

Tests can be run with:
  - go test crawler/...
  - Run automatically when calling make (but output is near start of stream)

## Considerations

  - Only crawls `a` and `img` tags, trivial to expand linkparser to match more
  - Uses tokenizer to parse html rather than dom parser or regex
    - Regex can easilly match FQDNs but becomes hard to manage matching relative links
    - DOM parser requires fairly correct markup to be able to accurate extract elements
    - Tokenizer walks through all the tags so easy to target the required elements

## Structure
  - Currently a lot of the code is in the main file with only linkparser modularised
  - Should just use main as an entrypoint then build the crawler functionality into a module

## Focus
  - Focus was on future performance and versatility
  - Utilises a queue to recurse through pages and process
    - Easy to build on later and multi-thread if greater speed is needed
    - Easy to chain multiple queues to give a staged pipeline with more diverse workers
    - Use of the queue and asynchronous message model makes termination a hard problem
      - Workaround waiting for the queue to be empty for a few seconds to end process
  - Basic tests around linkparser but much more coverage needed around the main crawl functions

## Output
  - CSV like
  - Prefixed with:
    - Crawling - HTML page to crawl
    - Resource - Page within current domain but not html
    - External - Page external to current domain
