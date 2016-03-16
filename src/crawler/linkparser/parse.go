package linkparser

import (
	"io"
	"golang.org/x/net/html"
	"net/url"
	"errors"
)

func Parse(body io.Reader) []string{
	var links []string;

	tokenized := html.NewTokenizer(body)
	for {
		tokentype := tokenized.Next()
		if tokentype == html.ErrorToken {
			if tokenized.Err().Error() == "EOF" {
				return links
			}else{
				panic(tokenized.Err())
			}
		}

		if tokentype == html.StartTagToken || tokentype == html.SelfClosingTagToken {
			token := tokenized.Token()

			if token.Data == "a" {
				for _, attribute := range token.Attr {
					if attribute.Key == "href" {
						links = append(links, attribute.Val)
					}
				}
			}
		}
	}
}

func ProcessLinks(links []string, base string) []string{
	for i := range links {
		var err error
		links[i], err = processLink(links[i],base)
		if err != nil {
			links = append(links[:i], links[i+1:]...)
		}
	}
	return links
}

func processLink(link, base string) (string, error) {
	uri, err := url.Parse(link)
	if err != nil {
		return "", errors.New("Failed to parse link")
	}
	baseUri, err := url.Parse(base)
	if err != nil {
		return "", errors.New("Failed to parse base")
	}
	absuri := baseUri.ResolveReference(uri)
	return absuri.String(), nil
}