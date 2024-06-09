package main

import (
	"context"
	"log"
	"sync"
	"time"

	database "github.com/JulianKerns/blog_aggregator/internal/database"
)

func StartWorker(db *database.Queries, limit int, timeBetween time.Duration) {
	log.Printf("Collecting feeds every %s on %v goroutines...", timeBetween, limit)
	ticker := time.NewTicker(timeBetween)

	for ; ; <-ticker.C {
		nextFeeds, err := db.GetNextFeedsToFetch(context.Background(), int32(limit))
		if err != nil {
			log.Println("could not retrieve the feeds from the database")
			continue

		}
		log.Printf("Found %v feeds to fetch\n", len(nextFeeds))
		wg := &sync.WaitGroup{}
		for _, feed := range nextFeeds {
			wg.Add(1)
			go ScrapingFeed(db, wg, feed)
		}
		wg.Wait()

	}
}
