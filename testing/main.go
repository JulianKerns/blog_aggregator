package main

func main() {
	const url string = "https://blog.boot.dev/index.xml"
	_, err := FetchingRSSFeed(url)
	if err != nil {
		return
	}
}
