package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
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
	getSEOData(resp *http.Response) (SEOData, error)
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
		foundSitemap := strings.Contains(page, "xml")
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
	worklist := make(chan []string)
	toCrawl := []string{}

	go func(){worklist <- []string{startURL}}()

	var n int = 1
	for n > 0 {
		list := <-worklist
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
					worklist <- sitemapFiles
				}
				for _, page := range pages {
					toCrawl = append(toCrawl, page)
				}
			}(link)
		}
	}

	return toCrawl
}

func crawlPage(url string, token chan struct{}) (*http.Response, error){
	token <- struct{}{}

	resp, err := makeRequest(url)
	<-token
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func scrapeURLs(urls []string, parser Parser, concurreny int) []SEOData{
	tokens := make(chan struct{}, concurreny)
	var n int
	worklist := make(chan []string)
	results := []SEOData{}

	go func() {worklist <- urls}()
	for n>0 {
		list := <-worklist
		for _, url := range list {
			if url != "" {
				n++
				go func(url string, token chan struct{}){
					log.Printf("Requesting URL:%s", url)
					res, err := scrapePage(url, tokens, parser)
					if err != nil {
						log.Printf("Encountered Error, URL:%s", url)
					}else {
						results = append(results, res)
					}
					worklist <- []string{}
				}(url, tokens)
			}
		}
	}

	return results
}

func extractURLs(resp *http.Response) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	results := []string{}
	sel := doc.Find("loc")

	for i := range sel.Nodes {
		loca := sel.Eq(i)
		result := loca.Text()
		results = append(results, result)
	}
	return results, nil
}

func scrapePage(url string, token chan struct{}, parser Parser) (SEOData, error){
	res, err := crawlPage(url, token)
	if err != nil {
		return SEOData{}, err
	}
	data, err := parser.getSEOData(res)
	if err != nil {
		return SEOData{}, err
	}
	return data, nil
}

func (d DefaultParser)getSEOData(resp *http.Response) (SEOData, error){
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return SEOData{}, err
	}
	result := SEOData{}
	result.URL = resp.Request.URL.String()
	result.StatusCode = resp.StatusCode
	result.Title = doc.Find("title").First().Text()
	result.H1 = doc.Find("h1").First().Text()
	result.MetaDescription, _ = doc.Find("meta[name^=description]").Attr("content")
	return result, nil
}

func scrapeSitemap(url string, parser Parser, concurrency int) []SEOData {
	results := extractSitemapURLs(url)
	res := scrapeURLs(results, parser, concurrency)
	return res
}

func main() {
	p := DefaultParser{}
	results := scrapeSitemap("https://www.quicksprout.com/sitemap_index.xml", p, 10)
	for _, res := range results {
		fmt.Println(res)
	}
}