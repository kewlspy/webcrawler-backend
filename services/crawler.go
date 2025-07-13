package services

import (
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/kewlspy/web-backend/models"
)

func CrawlURL(dbURL *models.URL) {
	dbURL.Status = "running"
	models.DB.Save(dbURL)

	c := colly.NewCollector(
		colly.MaxDepth(1),
	)

	internal := 0
	external := 0
	broken := 0
	loginForm := false
	headingCount := map[string]int{}

	var brokenLinks []models.BrokenLink

	c.OnHTML("html", func(e *colly.HTMLElement) {
		doc := e.DOM
		dbURL.Title = doc.Find("title").Text()

		// Count headings h1-h6
		for i := 1; i <= 6; i++ {
			tag := "h" + string(rune('0'+i))
			headingCount[tag] = doc.Find(tag).Length()
		}

		// Count and validate links
		doc.Find("a").Each(func(_ int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if !exists || href == "" {
				return
			}

			fullLink := resolveURL(dbURL.Link, href)

			if isInternalLink(dbURL.Link, fullLink) {
				internal++
			} else {
				external++
			}

			resp, err := http.Head(fullLink)
			status := 0
			if err != nil || resp == nil || resp.StatusCode >= 400 {
				if resp != nil {
					status = resp.StatusCode
				}
				broken++
				brokenLinks = append(brokenLinks, models.BrokenLink{
					URLID:  dbURL.ID,
					Link:   fullLink,
					Status: status,
				})
			} else {
				resp.Body.Close()
			}
		})

		// Detect login form
		doc.Find("form").Each(func(i int, s *goquery.Selection) {
			if s.Find("input[type='password']").Length() > 0 {
				loginForm = true
			}
		})
	})

	err := c.Visit(dbURL.Link)
	if err != nil {
		log.Printf("Failed to crawl URL: %v", err)
		dbURL.Status = "error"
		models.DB.Save(dbURL)
		return
	}

	// Update results
	dbURL.Status = "done"
	dbURL.InternalLinks = internal
	dbURL.ExternalLinks = external
	dbURL.BrokenLinks = broken
	dbURL.HasLoginForm = loginForm
	dbURL.HTMLVersion = "HTML5" 

	models.DB.Save(dbURL)

	if len(brokenLinks) > 0 {
		models.DB.Create(&brokenLinks)
	}

	log.Printf("Crawl complete for: %s", dbURL.Link)
}

func resolveURL(base, ref string) string {
	u, err := url.Parse(ref)
	if err != nil || u.Scheme == "" {
		baseURL, _ := url.Parse(base)
		return baseURL.ResolveReference(u).String()
	}
	return ref
}

func isInternalLink(base string, link string) bool {
	baseHost := getHost(base)
	linkHost := getHost(link)
	return baseHost == linkHost
}

func getHost(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return parsed.Hostname()
}
