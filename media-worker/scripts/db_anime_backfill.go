package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/machinebox/graphql"
	"log"
	"media-worker/database"
	"media-worker/media"
	"net/http"
	"os"
	"strconv"
	"time"
)

type headerCapturingTransport struct {
	underlyingTransport http.RoundTripper
	headers             http.Header
}

func (h *headerCapturingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := h.underlyingTransport.RoundTrip(req)
	if err == nil {
		h.headers = resp.Header
	}
	return resp, err
}

func discoverAnime(page int) (response media.AnimeQueryResponse, timeout int, err error) {
	capturingTransport := &headerCapturingTransport{
		underlyingTransport: http.DefaultTransport,
	}

	httpClient := &http.Client{Transport: capturingTransport}
	graphqlClient := graphql.NewClient("https://graphql.anilist.co", graphql.WithHTTPClient(httpClient))

	graphqlRequest := graphql.NewRequest(media.DiscoverAnime)
	graphqlRequest.Var("page", page)
	var graphqlResponse media.AnimeQueryResponse

	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		if timeout := capturingTransport.headers.Get("Retry-After"); timeout != "" {
			timeoutSeconds, _ := strconv.Atoi(timeout)
			return media.AnimeQueryResponse{}, timeoutSeconds, nil
		}
		return media.AnimeQueryResponse{}, 0, err
	}

	return graphqlResponse, 0, nil
}

func insertAnime(
	ctx context.Context,
	pool *pgxpool.Pool,
	q *database.Queries,
	anime media.AnimeDetails,
) error {
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	qtx := q.WithTx(tx)
	defer tx.Rollback(ctx)

	toText := func(s string) pgtype.Text { return pgtype.Text{String: s, Valid: s != ""} }
	toInt4 := func(n int) pgtype.Int4 { return pgtype.Int4{Int32: int32(n), Valid: n != 0} }
	toDate := func(fuzzyDate media.FuzzyDate) sql.NullString {
		return sql.NullString{String: fmt.Sprintf("(%d, %d, %d)",
			anime.StartDate.Year, anime.StartDate.Month, anime.StartDate.Day),
			Valid: fuzzyDate != (media.FuzzyDate{})}
	}

	var studios []string
	for _, v := range anime.Studios.Nodes {
		studios = append(studios, v.Name)
	}

	if err := qtx.PutAnime(ctx, database.PutAnimeParams{
		ID:           int32(anime.ID),
		Column2:      anime.Titles.Romaji,
		Column3:      anime.Titles.English,
		Column4:      anime.Titles.Native,
		Format:       toText(anime.Format),
		Status:       anime.Status,
		Season:       toText(anime.Season),
		SeasonYear:   toInt4(anime.SeasonYear),
		Episodes:     toInt4(anime.Episodes),
		CoverImage:   toText(anime.CoverImage.Large),
		Genres:       anime.Genres,
		AverageScore: toInt4(anime.AverageScore),
		Studio:       studios,
	}); err != nil {
		return err
	}

	var recommendations []string
	var scores []string
	airingSch := sql.NullString{String: "", Valid: false}

	if len(anime.AiringSchedule.Nodes) != 0 {
		node := anime.AiringSchedule.Nodes[0]
		airingSch = sql.NullString{String: fmt.Sprintf("(%d, %d)", node.Episode, node.AiringAt), Valid: true}
	}

	for _, v := range anime.Recommendations.Nodes {
		recommendationSting := fmt.Sprintf("(%d, %d)", v.MediaRecommendation.ID, v.Rating)
		recommendations = append(recommendations, recommendationSting)
	}
	for _, v := range anime.Stats.ScoreDistribution {
		scoreString := fmt.Sprintf("(%d, %d)", v.Score, v.Amount)
		scores = append(scores, scoreString)
	}

	if err := qtx.PutAnimeDetails(ctx, database.PutAnimeDetailsParams{
		ID:          int32(anime.ID),
		Description: pgtype.Text{String: anime.Description, Valid: true},
		StartDate:   toDate(anime.StartDate),
		EndDate:     toDate(anime.EndDate),
		Duration:    toInt4(anime.Duration),
		Source:      toText(anime.Source),
		Trailer: pgtype.Text{String: fmt.Sprintf("www.%s.com/watch?v=%s",
			anime.Trailer.Site, anime.Trailer.ID), Valid: anime.Trailer.Site != ""},
		BannerImage:       toText(anime.BannerImage),
		Popularity:        int32(anime.Popularity),
		Trending:          int32(anime.Trending),
		Favourites:        int32(anime.Favourites),
		AiringSchedule:    airingSch,
		Recommendations:   recommendations,
		ScoreDistribution: scores,
	}); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func dbWorker(
	id int,
	jobs <-chan media.AnimeDetails,
	failedJobs chan media.AnimeDetails,
	pool *pgxpool.Pool,
	q *database.Queries,
) {
	for job := range jobs {
		fmt.Printf("Worker %d processing anime with id: %d\n", id, job.ID)
		if err := insertAnime(context.Background(), pool, q, job); err != nil {
			log.Printf("worker %d: %v", id, err)
			failedJobs <- job
		}
	}
}

func main() {
	// the AniList API is degraded right now so the limit is 30 instead of 90 per min
	const rateLimitPerMin = 30
	var failedList []int

	stop := false
	limiter := time.Now().Unix()
	numQueried := 0
	jobs := make(chan media.AnimeDetails)
	failedJobs := make(chan media.AnimeDetails)

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

	poolCfg, _ := pgxpool.ParseConfig(connStr)
	poolCfg.MaxConns = 10
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		log.Fatal(err)
	}

	q := database.New(pool)
	start := time.Now()

	for i := 1; i <= 10; i++ {
		go dbWorker(i, jobs, failedJobs, pool, q)
	}

	for page := 1; stop == false; {
		currentTime := time.Now().Unix()

		if numQueried%rateLimitPerMin == 0 {
			// need to wait if the time taken to query 30 is less than 1 min
			if currentTime < limiter {
				fmt.Printf("Waiting for %d seconds to stay under rate limit\n", limiter-currentTime)
				time.Sleep(time.Duration(limiter-currentTime) * time.Second)
			}
			limiter = time.Now().Unix() + 65
			numQueried = 0
		}

		fmt.Printf("Starting page: %d\n", page)
		numQueried++
		response, timeout, err := discoverAnime(page)
		if err != nil {
			fmt.Printf("Page %d failed with error: %s\n", page, err)
			failedList = append(failedList, page)
			page++
			continue
		}

		// the rate limits reset after timeout, need to start new 30 cycle
		if timeout != 0 {
			fmt.Printf("Received time out after querying page: %d\n", page)
			fmt.Printf("Waiting for: %d Seconds\n", timeout)
			time.Sleep(time.Duration(timeout) * time.Second)
			numQueried = 0
			limiter = time.Now().Unix()
			continue
		}

		// Split the individual anime to worker
		for _, anime := range response.Page.Media {
			jobs <- anime
		}

		if response.Page.PageInfo.HasNextPage == false {
			stop = true
		}
		page++
		// TODO try again for later
		time.Sleep(400 * time.Millisecond)
	}

	close(jobs)
	close(failedJobs)

	fmt.Printf("All pages queried with %d failed pages and %d failed inserts\n", len(failedList), len(failedJobs))
	fmt.Printf("Took %s\n", time.Since(start))
}
