package testing

import (
	"fmt"
	"testing"
)

func TestGettingFeedPosts(t *testing.T) {
	cases := []struct {
		title string
		url   string
	}{
		{
			title: "The Boot.dev Beat. June 2024",
			url:   "https://blog.boot.dev/index.xml",
		},
		{
			title: "The Zen of Proverbs",
			url:   "https://www.wagslane.dev/index.xml",
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("Testing case %v", i), func(t *testing.T) {
			allItems, err := FetchingRSSFeed(c.url)
			if err != nil {
				t.Errorf("expected to fetch the RSS posts")
			}
			if allItems[0].Title != c.title {
				t.Errorf("expected the value to be equal")
			}
		})
	}
}
