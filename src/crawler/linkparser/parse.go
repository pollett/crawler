package linkparser

import (
	"io"
	"golang.org/x/net/html"
)

func Parse(body io.Reader) []string{
	var urls []string;

	tokenized := html.NewTokenizer(body)
	for {
		tokentype := tokenized.Next()
		if tokentype == html.ErrorToken {
			if tokenized.Err().Error() == "EOF" {
				return urls
			}else{
				panic(tokenized.Err())
			}
		}

		if tokentype == html.StartTagToken || tokentype == html.SelfClosingTagToken {
			token := tokenized.Token()

			if token.Data == "a" {
				for _, attribute := range token.Attr {
					if attribute.Key == "href" {
						urls = append(urls, attribute.Val)
					}
				}
			}
		}
	}
}