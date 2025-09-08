package media

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"media-worker/database"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func UpdatePage(url string, query string, pages []int) {
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

	for _, page := range pages {
		success := false

		for attempt := 1; attempt <= 5; attempt++ {
			fmt.Printf("Starting page: %d with attempt: %d\n", page, attempt)
			response, timeout, err := discoverMedia(url, query, page, nil)
			if err != nil {
				fmt.Printf("Page %d failed with error: %s\n", page, err)
				continue
			}

			if timeout != 0 {
				fmt.Printf("rate‑limit timeout (%d s); sleeping...", timeout)
				time.Sleep(time.Duration(timeout) * time.Second)
				continue
			}

			for _, media := range response.Page.Media {
				fmt.Printf("Starting media with ID: %d\n", media.ID)
				if err := insertMedia(ctx, pool, q, media); err != nil {
					fmt.Printf("Media with id %d failed with error: %s\n", media.ID, err)
				}
			}
			success = true
			break
		}

		if !success {
			fmt.Printf("Page %d failed\n", page)
		} else {
			fmt.Printf("Page %d succeeded\n", page)
		}
	}
}

func UpdateMedia(url string, query string, idList []int32) {
	const (
		rateLimitPerMin = 30
		windowSeconds   = 65.0
	)
	var failedPages []int
	var failedIds []int

	stop := false
	jobs := make(chan MediaDetails)
	failedJobs := make(chan MediaDetails)
	done := make(chan bool, 1)

	// loading the async failed ids onto failedIds
	go func(done chan bool) {
		for failed := range failedJobs {
			failedIds = append(failedIds, failed.ID)
		}
		done <- true
	}(done)

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

	var wg sync.WaitGroup
	q := database.New(pool)
	start := time.Now()

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			dbWorker(id, jobs, failedJobs, pool, q)
		}(i)
	}

	windowStart := time.Now()
	idx := 0

	for page := 1; stop == false; {
		fmt.Printf("Starting page: %d\n", page)
		response, timeout, err := discoverMedia(url, query, page, idList)
		if err != nil {
			fmt.Printf("Page %d failed with error: %s\n", page, err)
			failedPages = append(failedPages, page)
			page++
			continue
		}

		// the rate limits reset after timeout, need to start new 30 cycle
		if timeout != 0 {
			fmt.Printf("rate‑limit timeout (%d s); sleeping...", timeout)
			time.Sleep(time.Duration(timeout) * time.Second)
			idx, windowStart = 0, time.Now()
			continue
		}

		// Split the individual media to worker
		for _, media := range response.Page.Media {
			jobs <- media
		}

		fraction := float64(idx+1) / float64(rateLimitPerMin)
		targetElapsed := time.Duration(
			windowSeconds * math.Pow(fraction, 1.3) * float64(time.Second),
		)
		sleepFor := targetElapsed - time.Since(windowStart)
		time.Sleep(sleepFor)

		idx++
		if idx == rateLimitPerMin {
			idx, windowStart = 0, time.Now()
		}

		if response.Page.PageInfo.HasNextPage == false {
			stop = true
		}

		page++
	}

	close(jobs)
	wg.Wait()
	close(failedJobs)
	<-done

	for _, page := range failedPages {
		fmt.Printf("Page %d failed\n", page)
	}

	for _, id := range failedIds {
		fmt.Printf("Media with id %d failed\n", id)
	}

	fmt.Printf("All pages queried with %d failed pages and %d failed inserts\n", len(failedPages), len(failedIds))
	fmt.Printf("Took %s\n", time.Since(start))
}

func dbWorker(
	id int,
	jobs <-chan MediaDetails,
	failedJobs chan MediaDetails,
	pool *pgxpool.Pool,
	q *database.Queries,
) {
	for job := range jobs {
		fmt.Printf("Worker %d processing media with id: %d\n", id, job.ID)
		if err := insertMedia(context.Background(), pool, q, job); err != nil {
			log.Printf("worker %d: %v", id, err)
			failedJobs <- job
		}
	}
}

func discoverMedia(url string, query string, page int, idList []int32) (response MediaQueryResponse, timeout int, err error) {
	handler := NewGraphQLHandler(url)
	variables := map[string]interface{}{
		"page": page,
	}

	if idList != nil {
		variables["ids"] = idList
	}

	var graphqlResponse MediaQueryResponse

	if headers, err := handler.Query(query, variables, &graphqlResponse); err != nil {
		if timeout := headers.Get("Retry-After"); timeout != "" {
			timeoutSeconds, _ := strconv.Atoi(timeout)
			return MediaQueryResponse{}, timeoutSeconds, nil
		}
		return MediaQueryResponse{}, 0, err
	}

	return graphqlResponse, 0, nil
}

func insertMedia(
	ctx context.Context,
	pool *pgxpool.Pool,
	q *database.Queries,
	media MediaDetails,
) error {
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	qtx := q.WithTx(tx)
	defer tx.Rollback(ctx)

	toText := func(s string) pgtype.Text { return pgtype.Text{String: s, Valid: s != ""} }
	toInt4 := func(n int) pgtype.Int4 { return pgtype.Int4{Int32: int32(n), Valid: n != 0} }
	toBool := func(b bool) pgtype.Bool { return pgtype.Bool{Bool: b, Valid: true} }
	toDate := func(fuzzyDate FuzzyDate) sql.NullString {
		return sql.NullString{String: fmt.Sprintf("(%d, %d, %d)",
			media.StartDate.Year, media.StartDate.Month, media.StartDate.Day),
			Valid: fuzzyDate != (FuzzyDate{})}
	}
	toNullMediaType := func(s string) database.NullMediaType {
		return database.NullMediaType{MediaType: database.MediaType(s), Valid: s != ""}
	}

	var studios []string
	for _, v := range media.Studios.Nodes {
		studios = append(studios, v.Name)
	}

	if err := qtx.PutMedia(ctx, database.PutMediaParams{
		ID:           int32(media.ID),
		Column2:      media.Titles.Romaji,
		Column3:      media.Titles.English,
		Column4:      media.Titles.Native,
		Type:         toNullMediaType(media.Type),
		Format:       toText(media.Format),
		Status:       media.Status,
		Season:       toText(media.Season),
		SeasonYear:   toInt4(media.SeasonYear),
		Episodes:     toInt4(media.Episodes),
		Chapters:     toInt4(media.Chapters),
		Volumes:      toInt4(media.Volumes),
		CoverImage:   toText(media.CoverImage.Large),
		Genres:       media.Genres,
		AverageScore: toInt4(media.AverageScore),
		Studios:      studios,
		IsAdult:      toBool(media.IsAdult),
	}); err != nil {
		return err
	}

	var recommendations []string
	var scores []string
	airingSch := sql.NullString{String: "", Valid: false}

	if len(media.AiringSchedule.Nodes) != 0 {
		node := media.AiringSchedule.Nodes[0]
		airingSch = sql.NullString{String: fmt.Sprintf("(%d, %d)", node.Episode, node.AiringAt), Valid: true}
	}

	for _, v := range media.Recommendations.Nodes {
		recommendationSting := fmt.Sprintf("(%d, %d)", v.MediaRecommendation.ID, v.Rating)
		recommendations = append(recommendations, recommendationSting)
	}
	for _, v := range media.Stats.ScoreDistribution {
		scoreString := fmt.Sprintf("(%d, %d)", v.Score, v.Amount)
		scores = append(scores, scoreString)
	}

	if err := qtx.PutMediaDetails(ctx, database.PutMediaDetailsParams{
		ID:          int32(media.ID),
		Description: pgtype.Text{String: media.Description, Valid: true},
		StartDate:   toDate(media.StartDate),
		EndDate:     toDate(media.EndDate),
		Duration:    toInt4(media.Duration),
		Country:     toText(media.Country),
		Source:      toText(media.Source),
		Trailer: pgtype.Text{String: fmt.Sprintf("www.%s.com/watch?v=%s",
			media.Trailer.Site, media.Trailer.ID), Valid: media.Trailer != (Trailer{})},
		BannerImage:       toText(media.BannerImage),
		Popularity:        int32(media.Popularity),
		Trending:          int32(media.Trending),
		Favourites:        int32(media.Favourites),
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
