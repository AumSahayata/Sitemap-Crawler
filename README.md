# 🔍 Go SEO Sitemap Crawler

A fast and concurrent SEO crawler written in Go. It parses XML sitemaps (and sitemap indexes), crawls all listed URLs, and extracts essential SEO metadata from each page.

---

## 🚀 Features

- ✅ Parse nested sitemaps (sitemap index or single XML)
- ✅ Extract Title, H1, Meta Description, HTTP status code
- ✅ Concurrent and efficient scraping
- ✅ Custom User-Agent rotation
- ✅ Error-tolerant with logging

---

## 🏗️ How It Works

1. **Start with a sitemap URL**
2. **Recursively parse all linked sitemaps**
3. **Crawl each listed URL concurrently**
4. **Extract SEO metadata from HTML content**
5. **Print results**

---

## 📦 Dependencies

- [Go](https://golang.org/) 1.24+
- [goquery](https://github.com/PuerkitoBio/goquery)

Install dependencies:

```bash
go get github.com/PuerkitoBio/goquery
