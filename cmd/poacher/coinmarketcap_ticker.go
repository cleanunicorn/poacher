package main

import (
	"log"
	"time"

	"github.com/cleanunicorn/poacher/pkg/coinmarketcap"
	"github.com/influxdata/influxdb/client/v2"
)

func saveTickers(c client.Client, currency string) {
	// Create batch points
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "coinmarketcap",
		Precision: "s",
	})
	if err != nil {
		log.Fatalf("Could not create batchpoints, err: %v", err)
	}

	// Get all tickers and create datapoints
	var pts []*client.Point
	start := coinmarketcap.TickerStartDefault
	for {
		// Get next ticker until we reach 0 coins
		ts, err := coinmarketcap.Ticker(start, start+coinmarketcap.TickerLimitMax, currency)
		if err != nil {
			log.Fatalf("Error getting ticker from coinmarketcap, err: %v", err)
			return
		}
		if len(ts.Data) == 0 {
			break
		}

		// Create the point
		tags := map[string]string{"currency": currency}
		for _, ticker := range ts.Data {
			fields := map[string]interface{}{
				"price":              ticker.Quotes[currency].Price,
				"volume_24h":         ticker.Quotes[currency].Volume24H,
				"market_cap":         ticker.Quotes[currency].MarketCap,
				"percent_change_1h":  ticker.Quotes[currency].PercentChange1H,
				"percent_change_24h": ticker.Quotes[currency].PercentChange24H,
				"percent_change_7d":  ticker.Quotes[currency].PercentChange7D,
			}

			pt, err := client.NewPoint(ticker.Name, tags, fields, time.Unix(ts.Metadata.Timestamp, 0))
			if err != nil {
				log.Fatalf("Could not create new point, err: %v", err)
			}
			pts = append(pts, pt)
		}

		// Advance index
		start = start + coinmarketcap.TickerLimitMax
	}
	bp.AddPoints(pts)

	if err := c.Write(bp); err != nil {
		log.Fatalf("Could not write batchpoints, err: %v", err)
	}

	if err := c.Close(); err != nil {
		log.Fatalf("Could not close client, err: %v", err)
	}
}
