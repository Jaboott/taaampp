-- name: PutMedia :exec
INSERT INTO media (id,
                   titles,
                   type,
                   format,
                   status,
                   season,
                   season_year,
                   episodes,
                   chapters,
                   volumes,
                   cover_image,
                   genres,
                   average_score,
                   studios,
                   is_adult)
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
        $13,
        $14,
        $15,
        $16,
        $17);

-- name: PutMediaDetails :exec
INSERT INTO media_details (id,
                           description,
                           start_date,
                           end_date,
                           duration,
                           country,
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
        $14,
        $15
       );

name: QueryHighPrioMedia :many
SELECT media.id
FROM media
 LEFT JOIN media_details
ON media.id = media_details.id
WHERE status IN ('RELEASING', 'NOT_YET_RELEASED')
   OR (
    status NOT IN ('RELEASING', 'NOT_YET_RELEASED')
        AND end_date IS NOT NULL
        AND EXTRACT(YEAR FROM CURRENT_DATE) = (end_date).year
	AND EXTRACT(MONTH FROM CURRENT_DATE) - (end_date).month <= 1
    );
