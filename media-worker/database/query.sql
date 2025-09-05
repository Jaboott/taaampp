-- name: PutAnime :exec
INSERT INTO anime (id,
                   titles,
                   format,
                   status,
                   season,
                   season_year,
                   episodes,
                   cover_image,
                   genres,
                   average_score,
                   studio)
VALUES ($1,
        ROW($2, $3, $4)::titles,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11,
        $12,
        $13);

-- name: PutAnimeDetails :exec
INSERT INTO anime_details (id,
                           description,
                           start_date,
                           end_date,
                           duration,
                           source,
                           trailer,
                           banner_image,
                           popularity,
                           trending,
                           favourites,
                           airing_schedule,
                           recommendations,
                           score_distribution)
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11,
        $12,
        $13,
        $14
       );

-- name: QueryHighPrioMedia :many
SELECT anime.id
FROM anime
 LEFT JOIN anime_details
ON anime.id = anime_details.id
WHERE status IN ('RELEASING', 'NOT_YET_RELEASED')
   OR (
    status NOT IN ('RELEASING', 'NOT_YET_RELEASED')
        AND end_date IS NOT NULL
        AND EXTRACT(YEAR FROM CURRENT_DATE) = (end_date).year
	AND EXTRACT(MONTH FROM CURRENT_DATE) - (end_date).month <= 1
    );
