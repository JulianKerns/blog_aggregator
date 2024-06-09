package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	database "github.com/JulianKerns/blog_aggregator/internal/database"
)

func ScrapingFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	feedPosts, err := FetchingRSSFeed(feed.Url)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	_, errMarking := db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
		UpdatedAt:     time.Now().UTC(),
		ID:            feed.ID,
	})
	if errMarking != nil {
		fmt.Printf("%v\n", errMarking)
		return
	}
	fmt.Printf("%v RSSFeed collected ...Posts found: %v\n", feedPosts.Channel.Title, len(feedPosts.Channel.Items))

}

type Item struct {
	XMLName xml.Name `xml:"item"`
	Text    string   `xml:",chardata"`
	Title   string   `xml:"title"`
	Url     string   `xml:"link"`
}
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Channel struct {
		XMLName xml.Name `xml:"channel"`
		Text    string   `xml:",chardata"`
		Title   string   `xml:"title"`
		Items   []Item   `xml:"item"`
	}
}

func FetchingRSSFeed(Url string) (*RSS, error) {
	response, err := http.Get(Url)
	if err != nil {
		fmt.Println("Get request not succesful")
		return nil, errors.New("get request not succesful")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get valid Response")
	}

	decoder := xml.NewDecoder(response.Body)
	params := RSS{}
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		fmt.Printf("%v:Could not decode the XML body\n", errDecode)
		return nil, errors.New("could not decode the XML body")
	}

	return &params, nil
}
