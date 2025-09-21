package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"media-worker/database"
	"media-worker/media"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type mediaListGetter func(ctx context.Context, q *database.Queries) ([]int32, error)

func main() {
	mode := flag.String("mode", "", "a string")
	flag.Parse()

	required := []string{"PG_HOST", "PG_USER", "PG_PASSWORD", "PG_DATABASE"}

	for _, key := range required {
		if os.Getenv(key) == "" {
			if err := godotenv.Load(".env"); err != nil {
				log.Panicf("No .env file found with err: %s\n", err.Error())
			}
		}
	}

	switch *mode {
	case "all":
		fmt.Println("starting db backfill with all media in 3 seconds")
		time.Sleep(3 * time.Second)
		uploadAllMedia()
	case "new":
		fmt.Println("starting db backfill with new media in 3 seconds")
		time.Sleep(3 * time.Second)
		uploadNewMedia()
	case "high":
		fmt.Println("updating high priority media in 3 seconds")
		time.Sleep(3 * time.Second)
		updateHighPrioMedia()
	case "low":
		fmt.Println("updating low priority media in 3 seconds")
		time.Sleep(3 * time.Second)
		updateLowPrioMedia()
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

func updateHighPrioMedia() {
	if err := runUpdate(func(ctx context.Context, q *database.Queries) ([]int32, error) {
		return q.QueryHighPrioMedia(ctx)
	}); err != nil {
		log.Fatal(err)
	}
}

func updateLowPrioMedia() {
	if err := runUpdate(func(ctx context.Context, q *database.Queries) ([]int32, error) {
		return q.QueryLowPrioMedia(ctx)
	}); err != nil {
		log.Fatal(err)
	}
}

func runUpdate(get mediaListGetter) error {
	ctx := context.Background()

	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s database=%s port=5432",
		os.Getenv("PG_HOST"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_DATABASE"),
	)

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	q := database.New(conn)

	mediaList, err := get(ctx, q)
	if err != nil {
		return err
	}
	if mediaList == nil {
		mediaList = []int32{}
	}

	fmt.Printf("Updating database with %d media\n", len(mediaList))

	media.UpdateMedia("https://graphql.anilist.co", media.UpdateFromMediaList, mediaList)
	return nil
}
