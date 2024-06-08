package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/JulianKerns/blog_aggregator/internal/database"
	testing "github.com/JulianKerns/blog_aggregator/internal/testing"
)

func (cfg *apiConfig) FetchFeeds() error {
	feedsToFetch, err := cfg.DB.GetNextFeedsToFetch(context.Background())
	if err != nil {
		fmt.Printf("Error: %v. Could not retrieve the feeds to fetch\n", err)
		return err
	}
	for _, feed := range feedsToFetch {
		convertedFeed := databaseFeedtoFeed(feed)
		feedPosts, err := testing.FetchingRSSFeed(convertedFeed.Url)
		if err != nil {
			fmt.Printf("%v\n", err)
			return err
		}
		markedFeed, errMarking := cfg.DB.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
			LastFetchedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
			UpdatedAt:     time.Now().UTC(),
			ID:            feed.ID,
		})
		if errMarking != nil {
			fmt.Printf("%v\n", errMarking)
			return errMarking
		}
		fmt.Printf("Getting lastest Post from RSSfeed %v called %v\n", feedPosts[0].Title, markedFeed.Name)

	}
	return nil
}
