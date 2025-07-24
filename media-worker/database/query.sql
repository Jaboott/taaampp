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
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11
       );
