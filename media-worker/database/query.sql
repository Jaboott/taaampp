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
                   is_adult,
                   last_updated)
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
        $17,
        NOW()
        )
ON CONFLICT (id) DO UPDATE
SET titles        = ROW($2, $3, $4)::titles,
type          = $5,
format        = $6,
status        = $7,
season        = $8,
season_year   = $9,
episodes      = $10,
chapters      = $11,
volumes       = $12,
cover_image   = $13,
genres        = $14,
average_score = $15,
studios       = $16,
is_adult      = $17,
last_updated = NOW();


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
VALUES ($1,
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
        $15)
ON CONFLICT (id) DO UPDATE
SET description        = $2,
start_date         = $3,
end_date           = $4,
duration           = $5,
country            = $6,
source             = $7,
trailer            = $8,
banner_image       = $9,
popularity         = $10,
trending           = $11,
favourites         = $12,
airing_schedule    = $13,
recommendations    = $14,
score_distribution = $15;


-- name: QueryHighPrioMedia :many
SELECT media.id
FROM media
LEFT JOIN media_details
ON media.id = media_details.id
WHERE (
    type = 'ANIME'
        AND (
        status IN ('RELEASING', 'NOT_YET_RELEASED')
            OR (
            status NOT IN ('RELEASING', 'NOT_YET_RELEASED')
                AND end_date IS NOT NULL
                AND EXTRACT(YEAR FROM CURRENT_DATE) = (end_date).year
		    AND EXTRACT(MONTH FROM CURRENT_DATE) - (end_date).month <= 1
            )
        )
    )
OR (
type = 'MANGA'
    AND popularity >= 500
    AND (
    status IN ('RELEASING', 'NOT_YET_RELEASED')
        OR (
        status NOT IN ('RELEASING', 'NOT_YET_RELEASED')
            AND end_date IS NOT NULL
            AND EXTRACT(YEAR FROM CURRENT_DATE) = (end_date).year
        AND EXTRACT(MONTH FROM CURRENT_DATE) - (end_date).month <= 1
        )
    )
);