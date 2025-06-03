package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"math/rand"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type SEOData struct {
	URL				string
	Title			string
	H1				string
	MetaDescription	string
	StatusCode		int
}

type Parser interface {
	
}

type DefaultParser struct {

}

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.3.1 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36>",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:15.0) Gecko/20100101 Firefox/15.0.1",
	"Mozilla/5.0 (Linux; Android 5.0.2; SAMSUNG SM-T550 Build/LRX22G) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/3.3 Chrome/38.0.2125.102 Safari/537.36",
	"Mozilla/5.0 (Linux; Android 7.0; SM-T827R4 Build/NRD90M) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.116 Safari/537.36",
}

func randomUserAgent() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randNum := r.Int() % len(userAgents)
	return userAgents[randNum]
}

func isSitemap(urls []string) ([]string, []string) {
	sitemapFiles := []string{}
	pages := []string{}
	for _, page := range urls {
		if foundSitemap == true {
			fmt.Sprintf("Found sitemap", page)
			sitemapFiles = append(sitemapFiles, page)
		} else {
			pages = append(pages, page)
		}
	}
	return sitemapFiles, pages
}

func makeRequest(targetURL string) (*http.Response, error){
	client := http.Client{
		Timeout: 10*time.Second,
	}
	req, err := http.NewRequest("GET", targetURL, nil)
	req.Header.Set("User-Agent", randomUserAgent())
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func extractSitemapURLs(startURL string) []string {
	Worklist := make(chan []string)
	toCrawl := []string{}

	go func(){Worklist <- []string{startURL}}()

	var n int = 1
	for n > 0 {
		list := <-Worklist
		n--

		for _, link := range list {
			n++
			go func(link string){
				resp, err := makeRequest(link)
				if err != nil {
					log.Printf("error requesting: %s", link)
				}
				urls, err := extractURLs(resp)
				if err != nil {
					log.Printf("error extracting document from response, URL: %s", link)
				}
				sitemapFiles, pages := isSitemap(urls)
				if sitemapFiles != nil {
					Worklist <- sitemapFiles
				}
				for _, page := range pages {
					toCrawl = append(toCrawl, page)
				}
			}(link)
		}
	}

	return toCrawl
}

func crawlPage() {

}

func scrapeURLs() {

}

func scrapePage(url string, parser Parser) (SEOData, error){
	res, err := crawlPage(url)
	if err != nil {
		return SEOData{}, err
	}
	data, err := parser.getSEOData(res)
	if err != nil {
		return SEOData{}, err
	}
	return data, nil
}

func getSEOData() {

}

func scrapeSitemap(url string) []SEOData {
	results := extractSitemapURLs(url)
	res := scrapeURLs(results)
	return res
}

func main() {
	p := DefaultParser{}
	results := scrapeSitemap("")
	for _, res := range results {
		fmt.Println(res)
	}
}