CREATE TYPE titles AS
(
    romaji  TEXT,
    english TEXT,
    native  TEXT
);

CREATE TYPE fuzzy_date AS
(
    year  INTEGER,
    month INTEGER,
    day   INTEGER
);

CREATE TYPE airing_schedule AS
(
    episode   INTEGER,
    airing_at BIGINT
);

CREATE TYPE recommendation AS
(
    id    INTEGER,
    likes INTEGER
);

CREATE TYPE score_distribution AS
(
    score  INTEGER,
    amount INTEGER
);

CREATE TYPE watchlist_status AS ENUM (
    'watching',
    'completed',
    'dropped',
    'plan_to_watch'
    );

CREATE TYPE media_type AS ENUM ('ANIME', 'MANGA');

CREATE TABLE users
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         TEXT UNIQUE NOT NULL,
    password_hash TEXT        NOT NULL,
    username      TEXT UNIQUE NOT NULL,
    avatar_url    TEXT,
    created_at    TIMESTAMPTZ      DEFAULT NOW()
);

CREATE TABLE watchlist
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID REFERENCES users (id) ON DELETE CASCADE,
    media_id      INTEGER REFERENCES media (id),
    private       BOOLEAN          DEFAULT FALSE,
    status        watchlist_status,
    score         INTEGER,
    progress      INTEGER,
    rewatch_count INTEGER,
    start_date    TIMESTAMP,
    end_date      TIMESTAMP
);

CREATE TABLE media
(
    id            INTEGER PRIMARY KEY,
    titles        titles                    NOT NULL,
    type          media_type,
    format        TEXT,
    status        TEXT                      NOT NULL,
    season        TEXT,
    season_year   INTEGER,
    episodes      INTEGER,
    chapters      INTEGER,
    volumes       INTEGER,
    cover_image   TEXT,
    genres        TEXT[],
    average_score INTEGER,
    studios       TEXT[],
    is_adult      BOOLEAN,
    last_updated  TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE media_details
(
    id                 INTEGER PRIMARY KEY REFERENCES media ON DELETE CASCADE,
    description        TEXT,
    start_date         fuzzy_date,
    end_date           fuzzy_date,
    duration           INTEGER,
    country            TEXT,
    source             TEXT,
    trailer            TEXT,
    banner_image       TEXT,
    popularity         INTEGER NOT NULL DEFAULT 0,
    trending           INTEGER NOT NULL DEFAULT 0,
    favourites         INTEGER NOT NULL DEFAULT 0,
    airing_schedule    airing_schedule,
    recommendations    recommendation[],
    score_distribution score_distribution[]
);