package main

import "time"

type TrackRequest struct {
	URL         string  `json:"url" binding:"required"`
	TaregtPrice float64 `json:"target_price" binding:"required"`
}

type Product struct {
	ID          int       `json:"id"`
	URL         string    `json:"url"`
	TargetPrice float64   `json:"targetPrice"`
	LastPrice   float64   `json:"lastPrice"`
	LastChecked time.Time `json:"lastChecked"`
	AlertSent   bool      `json:"alert_sent"`
}
