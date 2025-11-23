package main

import (
	"log"
	"net/http" //to send back http status codes

	"github.com/gin-gonic/gin" //to get context 'c'
)

type HandleService struct {
	storage *StorageService
	scraper *ScraperService
}

func HandleServiceReq(storageSvc *StorageService, scraperSvc *ScraperService) *HandleService {
	return &HandleService{
		storage: storageSvc,
		scraper: scraperSvc,
	}
}

func (h *HandleService) handleTrackProduct(c *gin.Context) {
	//create variable to hold our request data
	var req TrackRequest

	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request " + err.Error()})
		return
	}

	currentPrice, err := h.scraper.ScrapePriceService(req.URL)
	if err != nil {
		// If scraping fails (e.g., can't connect, or our
		// 'span.a-price-whole' selector wasn't found),
		// we must stop and tell the user.
		// 'http.StatusUnprocessableEntity' (422) is a good code.
		// It means "I understood your request, but I can't
		// process this specific URL."
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Failed to scrape the provided URL. Check the URL or selector.", "details": err.Error()})
		return
	}

	// Create the new tracking request in the db to track this product
	err = h.storage.CreateNewTrackRequest(req.URL, req.TaregtPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save Track Request: " + err.Error()})
		return
	}

	log.Printf("Successfully scraped initial price for %s: %.2f", req.URL, currentPrice)

	c.JSON(http.StatusCreated, gin.H{
		"message":       "Product is Being Tracked",
		"tracking_url":  req.URL,
		"target_price":  req.TaregtPrice,
		"current_price": currentPrice})
}
