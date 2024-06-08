package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
)

type Item struct {
	XMLName xml.Name `xml:"item"`
	Text    string   `xml:",chardata"`
	Title   string   `xml:"title"`
	Url     string   `xml:"link"`
}

func FetchingRSSFeed(Url string) ([]Item, error) {
	response, err := http.Get(Url)
	if err != nil {
		fmt.Println("Get request not succesful")
		return nil, errors.New("get request not succesful")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get valid Response")
	}

	type RSS struct {
		XMLName xml.Name `xml:"rss"`
		Text    string   `xml:",chardata"`
		Channel struct {
			XMLName xml.Name `xml:"channel"`
			Text    string   `xml:",chardata"`
			Items   []Item   `xml:"item"`
		}
	}

	decoder := xml.NewDecoder(response.Body)
	params := RSS{}
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		fmt.Printf("%v:Could not decode the XML body\n", errDecode)
		return nil, errors.New("could not decode the XML body")
	}
	//for _, post := range params.Channel.Items {
	//	fmt.Printf("%v\n", post.Title)
	//}
	return params.Channel.Items, nil
}
