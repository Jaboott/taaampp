package main

import (
	"context"
	"fmt"
	"log"
	"media-worker/database"
	"media-worker/media"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

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
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	queries := database.New(conn)

	mediaList, err := queries.QueryHighPrioMedia(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if mediaList == nil {
		mediaList = []int32{}
	}

	media.UpdateMedia("https://graphql.anilist.co", media.UpdateFromMediaList, mediaList)
}
