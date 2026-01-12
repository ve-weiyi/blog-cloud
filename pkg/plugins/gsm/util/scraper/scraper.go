package scraper

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

const DefaultTimeout = 10 * time.Second
const DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36"

type Scraper struct {
	httpClient *http.Client
}

func (s Scraper) Scrap(endpoint string) (doc *goquery.Document, err error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)

	if err != nil {
		return
	}

	req.Header.Set("User-Agent", DefaultUserAgent)

	res, err := s.httpClient.Do(req)

	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("%s %s: returned HTTP status %d", http.MethodGet, endpoint, res.StatusCode)
		err = errors.New(msg)
		return
	}

	defer res.Body.Close()

	root, err := html.Parse(res.Body)

	if err != nil {
		return
	}

	doc = goquery.NewDocumentFromNode(root)

	return
}

func NewScraper() *Scraper {
	scraper := Scraper{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
	return &scraper
}
