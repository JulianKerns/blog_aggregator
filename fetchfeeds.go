package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/lib/pq"

	"github.com/google/uuid"

	database "github.com/JulianKerns/blog_aggregator/internal/database"
)

func ScrapingFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	RSSfeedStruct, err := FetchingRSSFeed(feed.Url)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}
	_, errMarking := db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
		UpdatedAt:     time.Now().UTC(),
		ID:            feed.ID,
	})
	if errMarking != nil {
		log.Printf("%v\n", errMarking)
		return
	}

	RSSPosts := RSSfeedStruct.Channel.Items
	for _, RSSpost := range RSSPosts {
		_, err := db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       RSSpost.Title,
			Url:         RSSpost.Url,
			Description: sql.NullString{String: RSSpost.Description, Valid: true},
			PublishedAt: RSSpost.PublishedAt.Time,
			FeedID:      feed.ID,
		})

		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code == "23505" {
					continue

				} else {
					log.Printf("Database Error: %v\n", pqErr)
				}

			}
			log.Printf("Other Error: %v\n", err)

		}

	}
	log.Printf("%v RSSFeed collected... Posts found: %v... Writing to Database!\n",
		RSSfeedStruct.Channel.Title, len(RSSfeedStruct.Channel.Items))

}

// Custom decoder for the PubDate so it works with the time.Time type

type PubDate struct {
	time.Time
}

func (pb *PubDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		"Mon, 02 Jan 2006 15:04:05 MST",
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"2006-01-02T15:04:05Z07:00",
	}
	var parsedTime time.Time
	for _, layout := range layouts {
		parsedTime, err = time.Parse(layout, v)
		if err == nil {
			*pb = PubDate{parsedTime}
			return nil
		}
	}
	return errors.New("could not parse date")
}

type Item struct {
	XMLName     xml.Name `xml:"item"`
	Text        string   `xml:",chardata"`
	Title       string   `xml:"title"`
	Url         string   `xml:"link"`
	PublishedAt PubDate  `xml:"pubDate"`
	Description string   `xml:"description"`
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
		log.Println("Get request not succesful")
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
		log.Printf("%v: Could not decode the XML body\n", errDecode)
		return nil, errors.New("could not decode the XML body")
	}

	return &params, nil
}
