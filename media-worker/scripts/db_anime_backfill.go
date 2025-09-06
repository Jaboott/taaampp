package main

import "media-worker/media"

func main() {
	media.UpdateMedia("https://graphql.anilist.co", media.DiscoverMedia, nil)
}
