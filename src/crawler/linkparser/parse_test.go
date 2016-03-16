package linkparser

import (
	"testing"
	"errors"
	"io"
	"strings"
)

func testProcessData() []string {
	data := []string{
		"#asfd",
		"https://domain.com",
		"https://google.com",
	}
	return data
}

func expectedProcessData() []string {
	data := []string{
		"https://domain.com#asfd",
		"https://domain.com",
		"https://google.com",
	}
	return data
}

func testData() io.Reader {
	data := `
        <html>
        <body>
        <a href="/link.com" />
        <img src="img.jpg" />
        </body>
        </html>
        `
	return strings.NewReader(data);
}

func expectedLinks() []string {
	return []string{
		"/link.com",
		"img.jpg",
	}
}

func testDataNoLinks() io.Reader {
	data := `
        <html>
        <body>
        </body>
        </html>
        `
	return strings.NewReader(data);
}

func testDataBadFormat() io.Reader {
	data := `
        <html>
        <head>
        <title>asdf</title>
        <title>asdfd</title>
        <body>
        <body>
        <a href="/link.com" />
        <img src="img.jpg" />
        `
	return strings.NewReader(data);
}

func TestProcessLinks(t *testing.T) {

	links := ProcessLinks(testProcessData(),"https://domain.com")
	err := compareArrays(links,expectedProcessData())
	if err != nil {
		t.Errorf("Error: ", err.Error())
	}

}

func TestParse(t *testing.T){
	links := Parse(testData())
	expected := expectedLinks()

	err := compareArrays(links,expected)
	if err != nil {
		t.Errorf("Error: ", err.Error())
	}
}

func TestParseBadData(t *testing.T){
	links := Parse(testDataBadFormat())
	expected := expectedLinks()

	err := compareArrays(links,expected)
	if err != nil {
		t.Errorf("Error: ", err.Error())
	}
}

func TestParseNoLinks(t *testing.T){
	links := Parse(testDataNoLinks())

	if len(links) != 0 {
		t.Error("found links in no links data")
	}
}

func compareArrays(a, b []string) error{
	if len(a) != len(b) {
		return errors.New("Length not matching")
	}

	for i, _ := range a {
		if a[i] != b[i] {
			return errors.New("Expected element doesn't exist")
		}
	}

	return nil
}