package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type ScraperService struct {
	collect *colly.Collector
}

func NewScraperService() *ScraperService {

	// select specific domains that we want
	c := colly.NewCollector(
		colly.AllowedDomains("www.amazon.in", "amazon.in"),
	)

	// Set a User-Agent. This is CRITICAL.
	// Most websites will block a request that doesn't
	// look like a real browser. A User-Agent is the, #1 way to identify yourself.

	c.OnRequest(func(r *colly.Request) { // Use the func as soon as you make a request to the page
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
		log.Println("Visiting:", r.URL)
	}) //give this function the request object 'r' which will act like a request from browser and not scrapper
	// We're throwing away our "Go robot" name tag and putting on a "Regular Chrome User" name tag.

	// Handle any errors.
	c.OnError(func(r *colly.Response, err error) { // use this as soon as you face any error
		log.Printf("Request URL: %s failed with response: %v\nError: %v", r.Request.URL, r, err)
	})

	return &ScraperService{
		collect: c,
	}
}

// After getting the HTML, we match Pattern of Price
// c.OnHTML(): tell the collector, "ONce downloaded the HTML, please find all the elements that match "this" pattern string
func (s *ScraperService) ScrapePriceService(url string) (float64, error) {
	var foundPrice float64 = 0.0
	var foundError error

	// "this string" : span.a-price-whole (looked in amazon price inspect element)
	s.collect.OnHTML("span.a-price-whole", func(e *colly.HTMLElement) {
		// 'e' is the HTML element that was found in "span.a-price-whole" string

		priceText := strings.Replace(e.Text, ",", "", -1) // remove commans in string Price

		price, err := strconv.ParseFloat(priceText, 64)
		if err != nil {
			log.Printf("Failed to parse price text '%s': %v", priceText, err)
			foundError = err
			return
		}

		foundPrice = price

	})

	s.collect.Visit(url) //It will visit the URL, download the HTML,

	s.collect.Wait() //wait for collector to scrape the website as per our OnHTML func

	if foundError != nil {
		return 0, foundError
	}

	return foundPrice, nil
}
