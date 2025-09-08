package main

import (
	"flag"
	"fmt"
	"media-worker/media"
	"time"
)

func main() {
	mode := flag.String("mode", "", "a string")
	flag.Parse()

	switch *mode {
	case "all":
		fmt.Println("starting db backfill with all media in 3 seconds")
		time.Sleep(3 * time.Second)
		uploadAllMedia()
	case "new":
		fmt.Println("starting db backfill with new media in 3 seconds")
		time.Sleep(3 * time.Second)
		uploadNewMedia()
	default:
		fmt.Println("Invalid mode, please only enter either: 'all' or 'new'")
	}
}

// TODO instead of just getting the first 4 page, go until an id from query is already in db
func uploadNewMedia() {
	media.UpdatePage("https://graphql.anilist.co", media.DiscoverNewMedia, []int{1, 2, 3, 4})
}

func uploadAllMedia() {
	media.UpdateMedia("https://graphql.anilist.co", media.DiscoverMedia, nil)
}
