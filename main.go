package main

import (
	"fmt"
	"log"
)

type SEOData struct {
	URL				string
	Title			string
	H1				string
	MetaDescription	string
	StatusCode		int
}

type parser interface {
	
}

type DefaultParser struct {

}

// userAgents

func randomUserAgent() {

}

func makeRequest() {

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

func scrapePage() {

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