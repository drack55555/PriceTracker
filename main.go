package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	storage, err := NewStorageService("./price-alert.db")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	scraper := NewScraperService()

	emailConfig := NotifierCofig{
		SMTPServer:     "smtp.gmail.com",
		SMTPPort:       587,
		SenderEmail:    "1906117@kiit.ac.in",
		SenderPassword: "fueu mysm epap eprf",
	}
	notifier := NewNotifier(emailConfig)

	scheduler := NewSchedule(storage, scraper, notifier)

	// Creating a Goroutine that will tell program that run this in BG on it's own, and move code to next line
	// without waiting for this scheduler to finish.
	go scheduler.Run()

	handle := HandleServiceReq(storage, scraper)

	router := gin.Default()

	router.POST("/track", handle.handleTrackProduct)

	log.Printf("Starting the server at Port 8080")
	router.Run(":8080")
}
