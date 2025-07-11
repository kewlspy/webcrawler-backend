package services

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/kewlspy/web-backend/models"
)

func CrawlURL(url *models.URL) {
	url.Status = "running"
	models.DB.Save(url)

	c := colly.NewCollector()

	internal, external, broken := 0, 0, 0
	loginForm := false

	c.OnHTML("html", func(e *colly.HTMLElement) {
		doc := e.DOM

		url.Title = doc.Find("title").Text()

		for i := 1; i <= 6; i++ {
			url.InternalLinks += doc.Find("h" + string(rune(i+'0'))).Length()
		}

		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			link, _ := s.Attr("href")
			if strings.HasPrefix(link, "/") || strings.Contains(link, url.Link) {
				internal++
			} else {
				external++
			}

			resp, err := http.Get(link)
			if err != nil || resp.StatusCode >= 400 {
				broken++
				models.DB.Create(&models.BrokenLink{
					URLID:  url.ID,
					Link:   link,
					Status: resp.StatusCode,
				})
			}
		})

		doc.Find("form").Each(func(i int, s *goquery.Selection) {
			if s.Find("input[type='password']").Length() > 0 {
				loginForm = true
			}
		})
	})

	err := c.Visit(url.Link)
	if err != nil {
		url.Status = "error"
		models.DB.Save(url)
		return
	}

	url.Status = "done"
	url.InternalLinks = internal
	url.ExternalLinks = external
	url.BrokenLinks = broken
	url.HasLoginForm = loginForm

	models.DB.Save(url)
}
