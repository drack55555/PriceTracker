package main

import (
	"log"
	"sync"
	"time"
)

type SchedulerService struct {
	storage  *StorageService
	scrper   *ScraperService
	ticker   *time.Ticker
	notifier *Notifier
}

func NewSchedule(s *StorageService, sc *ScraperService, n *Notifier) *SchedulerService {
	return &SchedulerService{
		storage:  s,
		scrper:   sc,
		notifier: n,
	}
}

func (s *SchedulerService) Run() {
	// check every  20 seconds
	s.ticker = time.NewTicker(20 * time.Second) // generate a tick(signal) every 20 seconds
	log.Println("Scheduler started and Ticking at every 20 Seconds")

	// using Channel which acts like Conveyer belt for data(sends diff part of program each other messages safely)
	// using infinite loop so that it always runs in BG and keeps checking for price
	for range s.ticker.C { // the ticker pipe every 20 seconds send msg(the current time) to C channel which wakes channel up and enters the loop
		log.Println("Scheduled Ticker...Running Price Check at..")
		s.runPriceChecks()
	}
}

func (s *SchedulerService) runPriceChecks() {
	products, err := s.storage.GetAllProducts()
	if err != nil {
		log.Printf("Error getting product from storage: %v", err)
		return
	}

	log.Printf("Found %d products to check", len(products))

	var wg sync.WaitGroup

	for _, p := range products {
		// if we have 1000 of products to scan, sequentially it will take 1000 times, but we'll use
		// goroutines to do in parallel.  Goroutines will launch 1000 robots single time and scan in 5-10 secs max instead of 1000.

		wg.Add(1) // add count that 1 goroutie is running each time entering loop

		//launch gorotuine
		go func(productToCheck *Product) { // go keyword adds 1 goroutines each time it runs in this looop..
			defer wg.Done() //run at last of loop

			log.Printf("Check product with ID: %d, and at URL: %s ", productToCheck.ID, productToCheck.URL)
			currentPrice, err := s.scrper.ScrapePriceService(productToCheck.URL)
			if err != nil {
				log.Printf("Failed to scrape Product URL %s: %v", productToCheck.URL, err)
				return
			}
			log.Printf("Success! Product ID %d (%s) is now %.2f", productToCheck.ID, productToCheck.URL, currentPrice)

			newAlertState := productToCheck.AlertSent
			shouldSendEmail := false
			if currentPrice < productToCheck.TargetPrice {
				// Price is LOW
				if !productToCheck.AlertSent {
					// Scenario A: First drop! We need to alert.
					shouldSendEmail = true
					newAlertState = true
				}
				// Already alerted? newAlertState stays true, shouldSendEmail stays false.
			} else {
				// Price went back UP (or is equal).
				// We reset the flag so we can alert again next time it drops.
				newAlertState = false
			}
			// Save new price to database
			err = s.storage.UpdatePrice(productToCheck.ID, currentPrice, newAlertState)
			if err != nil {
				log.Printf("Failed to update price for %s: %v", productToCheck.URL, err)
				return
			}

			// CHECK new_price < target_price
			if shouldSendEmail {
				log.Printf("ALERT! Price for %s (%.2f) is below target (%.2f)!", productToCheck.URL, currentPrice, productToCheck.TargetPrice)

				// Notify about price change
				targetEmail := "kummarraj20@gmail.com"
				err := s.notifier.SendEmailAlert(targetEmail, productToCheck.URL, currentPrice)
				if err != nil {
					log.Printf("Failed to send email alert: %v", err)
				} else {
					log.Printf("Email alert sent successfully to %s", targetEmail)
				}
			}
		}(p) // this is like a private copy of the original item(single products), so as to solve race condition

	}

	wg.Wait()
	log.Printf("All %d product checks complete.", len(products))
}
